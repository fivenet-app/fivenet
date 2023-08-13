package centrum

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	tGksPhoneJMsg     = table.GksphoneJobMessage
	tGksPhoneSettings = table.GksphoneSettings
)

func (s *Server) ConvertPhoneJobMsgToDispatch() error {
	for {
		if err := s.convertPhoneJobMsgToDispatch(); err != nil {
			s.logger.Error("failed to convert gksphone job messages to dispatches", zap.Error(err))
		}

		<-time.After(2 * time.Second)
	}
}

func (s *Server) convertPhoneJobMsgToDispatch() error {
	stmt := tGksPhoneJMsg.
		SELECT(
			tGksPhoneJMsg.ID,
			tGksPhoneJMsg.Jobm,
			tGksPhoneJMsg.Anon,
			tGksPhoneJMsg.Gps,
			tGksPhoneJMsg.Message,
			tGksPhoneJMsg.Time,
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
			tGksPhoneJMsg.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(s.visibleJobs, "|")+")\"\\]")),
			tGksPhoneJMsg.Owner.EQ(jet.Int32(0)),
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
		var anon bool
		if msg.Anon != nil && *msg.Anon == "1" {
			anon = true
		}

		dsp := &dispatch.Dispatch{
			Job:     job,
			Message: *msg.Message,
			X:       x,
			Y:       y,
			Anon:    &anon,
			UserId:  &msg.UserId,
		}

		if _, err := s.createDispatch(s.ctx, dsp); err != nil {
			return err
		}

		if true {
			if err := s.closePhoneJobMsg(s.ctx, msg.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) closePhoneJobMsg(ctx context.Context, id int32) error {
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
