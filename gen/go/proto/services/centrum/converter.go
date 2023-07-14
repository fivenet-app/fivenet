package centrum

import (
	"context"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tGksPhoneJMsg = table.GksphoneJobMessage
)

func (s *Server) ConvertPhoneJobMsgToDispatch() error {
	stmt := tGksPhoneJMsg.
		SELECT(
			tGksPhoneJMsg.ID,
		).
		FROM().
		WHERE(
			tGksPhoneJMsg.ID.GT_EQ(jet.Int32(1)),
		)

	var dest []*model.GksphoneJobMessage
	if err := stmt.QueryContext(s.ctx, s.db, &dest); err != nil {
		return err
	}

	for _, m := range dest {
		job := strings.TrimSuffix(strings.TrimPrefix(*m.Jobm, "[\""), "\"]")
		gps, _ := strings.CutPrefix(*m.Gps, "GPS: ")
		gpsSplit := strings.Split(gps, ", ")
		x, _ := strconv.ParseFloat(gpsSplit[0], 32)
		y, _ := strconv.ParseFloat(gpsSplit[1], 32)
		var anon bool
		if m.Anon != nil && *m.Anon == "1" {
			anon = true
		}
		userId := int32(26061)

		dsp := &model.FivenetCentrumDispatches{
			Job:     &job,
			Message: *m.Message,
			X:       &x,
			Y:       &y,
			Anon:    &anon,
			UserID:  userId,
		}
		_ = dsp
		// TODO create dispatch in fivenet dispatch table

		if err := s.closePhoneJobMsg(s.ctx, m.ID); err != nil {
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
