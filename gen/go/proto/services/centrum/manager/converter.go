package manager

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
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
			tGksPhoneJMsg.Time.LT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(2, jet.SECOND)),
			),
		))

	var dest []struct {
		*model.GksphoneJobMessage
		UserId int32
	}
	if err := stmt.QueryContext(s.ctx, s.db, &dest); err != nil {
		return err
	}

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
			Attributes: &centrum.Attributes{},
			Job:        job,
			Message:    message,
			X:          x,
			Y:          y,
			Anon:       anon,
			CreatorId:  &msg.UserId,
		}

		if postal := s.postals.Closest(x, y); postal != nil {
			dsp.Postal = postal.Code
		}

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
