package jobs

import (
	"context"
	"errors"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type qualificationDisplayRow struct {
	ID           int64  `db:"id"`
	Abbreviation string `db:"abbreviation"`
	Title        string `db:"title"`
}

func int64Expressions(ids []int64) []mysql.Expression {
	exprs := make([]mysql.Expression, len(ids))
	for i := range ids {
		exprs[i] = mysql.Int64(ids[i])
	}

	return exprs
}

func (s *Server) enrichGroupRuleQualifications(
	ctx context.Context,
	rules ...*jobsgroups.GroupRule,
) error {
	ids := make([]int64, 0)
	seen := map[int64]struct{}{}

	for _, rule := range rules {
		qualification := rule.GetQualification()
		if qualification == nil {
			continue
		}

		for _, id := range qualification.GetQualificationIds() {
			if id <= 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	tQuali := table.FivenetQualifications.AS("qualification_display_row")
	stmt := tQuali.
		SELECT(
			tQuali.ID,
			tQuali.Abbreviation,
			tQuali.Title,
		).
		FROM(tQuali).
		WHERE(tQuali.ID.IN(int64Expressions(ids)...)).
		ORDER_BY(tQuali.ID.ASC()).
		LIMIT(int64(len(ids)))

	rows := []*qualificationDisplayRow{}
	if err := stmt.QueryContext(ctx, s.db, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	byID := make(map[int64]*resqualifications.QualificationShort, len(rows))
	for _, row := range rows {
		byID[row.ID] = &resqualifications.QualificationShort{
			Id:           row.ID,
			Abbreviation: row.Abbreviation,
			Title:        row.Title,
		}
	}

	for _, rule := range rules {
		qualification := rule.GetQualification()
		if qualification == nil {
			continue
		}

		hydrated := make(
			[]*resqualifications.QualificationShort,
			0,
			len(qualification.GetQualificationIds()),
		)
		for _, id := range qualification.GetQualificationIds() {
			if quali := byID[id]; quali != nil {
				hydrated = append(hydrated, quali)
				continue
			}

			hydrated = append(hydrated, &resqualifications.QualificationShort{Id: id})
		}
		qualification.Qualifications = hydrated
	}

	return nil
}
