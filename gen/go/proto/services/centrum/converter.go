package centrum

import (
	"context"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tGksPhoneJMsg     = table.GksphoneJobMessage
	tGksPhoneSettings = table.GksphoneSettings
)

func (s *Server) ConvertPhoneJobMsgToDispatch() error {
	stmt := tGksPhoneJMsg.
		SELECT(
			tGksPhoneJMsg.ID,
			tGksPhoneJMsg.Jobm,
			tGksPhoneJMsg.Anon,
			tGksPhoneJMsg.Gps,
			tGksPhoneJMsg.Message,
			tGksPhoneJMsg.Time,
			tUser.ID.AS("userId"),
		).
		FROM(
			tGksPhoneJMsg.
				INNER_JOIN(tGksPhoneSettings,
					tGksPhoneSettings.PhoneNumber.EQ(tGksPhoneJMsg.Number),
				).
				INNER_JOIN(tUser,
					tUser.Identifier.EQ(tGksPhoneSettings.Identifier),
				),
		).
		WHERE(
			tGksPhoneJMsg.Jobm.REGEXP_LIKE(jet.String("\\[\"(" + strings.Join(s.visibleJobs, "|") + ")\"\\]")),
		)

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
		x, _ := strconv.ParseFloat(gpsSplit[0], 32)
		y, _ := strconv.ParseFloat(gpsSplit[1], 32)
		var anon bool
		if msg.Anon != nil && *msg.Anon == "1" {
			anon = true
		}

		dsp := &dispatch.Dispatch{
			Job:     job,
			Message: *msg.Message,
			X:       float32(x),
			Y:       float32(y),
			Anon:    &anon,
			UserId:  &msg.UserId,
		}

		_, err := s.createDispatch(s.ctx, dsp)
		if err != nil {
			return err
		}

		if err := s.closePhoneJobMsg(s.ctx, msg.ID); err != nil {
			return err
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
