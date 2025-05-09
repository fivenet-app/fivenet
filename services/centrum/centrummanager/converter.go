package centrummanager

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	// Converter
	tGksPhoneJMsg     = table.GksphoneJobMessage
	tGksPhoneSettings = table.GksphoneSettings
)

func (s *Housekeeper) ConvertPhoneJobMsgToDispatch() error {
	if len(s.convertJobs) == 0 {
		return nil
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil

		case <-time.After(2 * time.Second):
		}

		if err := s.convertPhoneJobMsgToDispatch(); err != nil {
			s.logger.Error("failed to convert gksphone job messages to dispatches", zap.Error(err))
		}
	}
}

func (s *Housekeeper) convertPhoneJobMsgToDispatch() error {
	tUsers := tables.Users()

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
			tGksPhoneJMsg.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(s.convertJobs, "|")+")\"\\]")),
			tGksPhoneJMsg.Owner.EQ(jet.Int32(0)),
		)).
		LIMIT(15)

	var dest []struct {
		*model.GksphoneJobMessage
		UserId int32
	}
	if err := stmt.QueryContext(s.ctx, s.db, &dest); err != nil {
		return err
	}

	s.logger.Debug("converting phone dispatch to fivenet", zap.Int("dispatch_count", len(dest)))
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
			Job:        job,
			Message:    message,
			X:          x,
			Y:          y,
			Anon:       anon,
			CreatorId:  &msg.UserId,
		}

		s.logger.Debug("converted phone dispatch to fivenet", zap.String("job", job), zap.Int32("creator_id", msg.UserId), zap.Int32("phone_dsp_id", msg.ID))
		if _, err := s.CreateDispatch(s.ctx, dsp); err != nil {
			return err
		}

		if err := s.closePhoneJobMsg(s.ctx, msg.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Housekeeper) closePhoneJobMsg(ctx context.Context, id int32) error {
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
