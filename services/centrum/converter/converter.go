package converter

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	// GKSPhone Converter.
	tGksPhoneJMsg     = table.GksphoneJobMessage
	tGksPhoneSettings = table.GksphoneSettings

	// LBPhone Converter.
	tPhonePhones           = table.PhonePhones
	tPhoneServicesChannels = table.PhoneServicesChannels
	tPhoneServicesMessages = table.PhoneServicesMessages
)

type Converter struct {
	logger *zap.Logger
	db     *sql.DB

	dispatches *dispatches.DispatchDB

	converterType string
	convertJobs   []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	Config *config.Config

	Dispatches *dispatches.DispatchDB
}

func New(p Params) *Converter {
	ctxCancel, cancel := context.WithCancel(context.Background())

	c := &Converter{
		logger: p.Logger.Named("centrum.converter"),
		db:     p.DB,

		dispatches: p.Dispatches,

		converterType: p.Config.DispatchCenter.Type,
		convertJobs:   p.Config.DispatchCenter.ConvertJobs,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		c.convertPhoneJobMsgToDispatch(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return c
}

func (s *Converter) convertPhoneJobMsgToDispatch(ctx context.Context) {
	if len(s.convertJobs) == 0 {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}

		switch s.converterType {
		case "lbphone":
			if err := s.convertLBPhoneJobMsgToDispatch(ctx); err != nil {
				s.logger.Error(
					"failed to convert lbphone job messages to dispatches",
					zap.Error(err),
				)
			}

		case "gksphone":
			if err := s.convertGKSPhoneJobMsgToDispatch(ctx); err != nil {
				s.logger.Error(
					"failed to convert gksphone job messages to dispatches",
					zap.Error(err),
				)
			}

		default:
			s.logger.Error(
				"unknown phone dispatch converter type",
				zap.String("converter_type", s.converterType),
			)
		}
	}
}

func (s *Converter) convertGKSPhoneJobMsgToDispatch(ctx context.Context) error {
	tUsers := tables.User()

	stmt := tGksPhoneJMsg.
		SELECT(
			tGksPhoneJMsg.ID,
			tGksPhoneJMsg.Jobm,
			tGksPhoneJMsg.Anon,
			tGksPhoneJMsg.Gps,
			tGksPhoneJMsg.Message,
			tUsers.ID.AS("userid"),
		).
		FROM(
			tGksPhoneJMsg.
				INNER_JOIN(tGksPhoneSettings,
					tGksPhoneSettings.PhoneNumber.EQ(tGksPhoneJMsg.Number),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.EQ(tGksPhoneSettings.Identifier),
				),
		).
		WHERE(jet.AND(
			tGksPhoneJMsg.Jobm.REGEXP_LIKE(
				jet.String("\\[\"("+strings.Join(s.convertJobs, "|")+")\"\\]"),
			),
			tGksPhoneJMsg.Owner.EQ(jet.Int32(0)),
		)).
		LIMIT(15)

	var dest []struct {
		*model.GksphoneJobMessage

		UserId int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	s.logger.Debug("converting gksphone dispatch to fivenet", zap.Int("dispatch_count", len(dest)))
	for _, msg := range dest {
		job := strings.TrimSuffix(strings.TrimPrefix(*msg.Jobm, "[\""), "\"]")
		gps, _ := strings.CutPrefix(*msg.Gps, "GPS: ")
		gpsSplit := strings.Split(gps, ", ")
		x, err := strconv.ParseFloat(gpsSplit[0], 32)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(gpsSplit[1], 32)
		if err != nil {
			continue
		}

		anon := false
		if msg.Anon != nil && *msg.Anon == "1" {
			anon = true
		}

		message := "N/A"
		if msg.Message != nil {
			if len(*msg.Message) > 250 {
				message = utils.StringFirstN(*msg.Message, 250) + "..."
			} else {
				message = *msg.Message
			}
		}

		dsp := &centrum.Dispatch{
			CreatedAt:  timestamp.Now(),
			Attributes: &centrum.DispatchAttributes{},
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						Name: job,
					},
				},
			},
			Message:   message,
			X:         x,
			Y:         y,
			Anon:      anon,
			CreatorId: &msg.UserId,
		}

		s.logger.Debug(
			"converted gksphone dispatch to fivenet",
			zap.String("job", job),
			zap.Int32("creator_id", msg.UserId),
			zap.Int32("phone_dsp_id", msg.ID),
		)
		if _, err := s.dispatches.Create(ctx, dsp); err != nil {
			return err
		}

		if err := s.closeGKSPhoneJobMsg(ctx, msg.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Converter) closeGKSPhoneJobMsg(ctx context.Context, id int32) error {
	stmt := tGksPhoneJMsg.
		UPDATE(
			tGksPhoneJMsg.Owner,
		).
		SET(
			tGksPhoneJMsg.Owner.SET(jet.Int32(1)),
		).
		WHERE(
			tGksPhoneJMsg.ID.EQ(jet.Int32(id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Converter) convertLBPhoneJobMsgToDispatch(ctx context.Context) error {
	tUsers := tables.User()

	stmt := tPhoneServicesChannels.
		SELECT(
			tPhoneServicesChannels.ID,
			tPhoneServicesChannels.Company,
			tPhoneServicesChannels.PhoneNumber,
			tPhoneServicesMessages.Message,
			tPhoneServicesMessages.XPos,
			tPhoneServicesMessages.YPos,
			tUsers.ID.AS("userid"),
		).
		FROM(
			tPhoneServicesChannels.
				INNER_JOIN(tPhoneServicesMessages,
					tPhoneServicesMessages.ChannelID.EQ(tPhoneServicesChannels.ID),
				).
				INNER_JOIN(tPhonePhones,
					tPhonePhones.PhoneNumber.EQ(tPhoneServicesChannels.PhoneNumber),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.EQ(tPhonePhones.OwnerID),
				),
		).
		WHERE(jet.AND(
			tPhoneServicesChannels.Company.REGEXP_LIKE(
				jet.String("\\[\"(" + strings.Join(s.convertJobs, "|") + ")\"\\]"),
			),
		)).
		LIMIT(15)

	var dest []struct {
		ID          int32
		Job         string `alias:"company"`
		PhoneNumber string
		Message     string
		XPos        int32
		YPos        int32
		UserId      int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	s.logger.Debug("converting lbphone dispatch to fivenet", zap.Int("dispatch_count", len(dest)))
	for _, msg := range dest {
		message := "N/A"
		if msg.Message != "" {
			if len(msg.Message) > 250 {
				message = utils.StringFirstN(msg.Message, 250) + "..."
			} else {
				message = msg.Message
			}
		}

		dsp := &centrum.Dispatch{
			CreatedAt:  timestamp.Now(),
			Attributes: &centrum.DispatchAttributes{},
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						Name: msg.Job,
					},
				},
			},
			Message:   message,
			X:         float64(msg.XPos),
			Y:         float64(msg.YPos),
			Anon:      false,
			CreatorId: &msg.UserId,
		}

		s.logger.Debug(
			"converted lbphone dispatch to fivenet",
			zap.String("job", msg.Job),
			zap.Int32("creator_id", msg.UserId),
			zap.Int32("phone_dsp_id", msg.ID),
		)
		if _, err := s.dispatches.Create(ctx, dsp); err != nil {
			return err
		}

		if err := s.closeLBPhoneJobMsg(ctx, msg.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Converter) closeLBPhoneJobMsg(ctx context.Context, id int32) error {
	stmt := tPhoneServicesChannels.
		DELETE().
		WHERE(
			tPhoneServicesChannels.ID.EQ(jet.Int32(id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
