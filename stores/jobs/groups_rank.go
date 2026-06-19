package jobsstore

import (
	"context"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	storesrank "github.com/fivenet-app/fivenet/v2026/stores/internal/rank"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type jobGroupRankGroup struct {
	job string
}

func (g jobGroupRankGroup) ListRanks(
	ctx context.Context,
	q qrm.DB,
	excludeID int64,
) ([]storesrank.Row, error) {
	condition := mysql.AND(
		tJobGroups.Job.EQ(mysql.String(g.job)),
		tJobGroups.State.NOT_EQ(mysql.Int32(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED))),
	)
	if excludeID > 0 {
		condition = condition.AND(tJobGroups.ID.NOT_EQ(mysql.Int64(excludeID)))
	}

	stmt := tJobGroups.
		SELECT(
			tJobGroups.ID.AS("id"),
			tJobGroups.SortRank.AS("sort_rank"),
		).
		FROM(tJobGroups).
		WHERE(condition).
		ORDER_BY(tJobGroups.SortRank.ASC(), tJobGroups.ID.ASC()).
		FOR(mysql.UPDATE())

	ranks := []storesrank.Row{}
	if err := stmt.QueryContext(ctx, q, &ranks); err != nil {
		return nil, err
	}

	return ranks, nil
}

func (g jobGroupRankGroup) UpdateRank(
	ctx context.Context,
	q qrm.DB,
	id int64,
	sortRank string,
) error {
	updateStmt := tJobGroups.
		UPDATE(tJobGroups.SortRank).
		SET(tJobGroups.SortRank.SET(mysql.String(sortRank))).
		WHERE(mysql.AND(
			tJobGroups.ID.EQ(mysql.Int64(id)),
			tJobGroups.Job.EQ(mysql.String(g.job)),
		)).
		LIMIT(1)

	_, err := updateStmt.ExecContext(ctx, q)
	return err
}

func (s *Store) nextGroupRank(
	ctx context.Context,
	db qrm.DB,
	job string,
	excludeID int64,
) (string, error) {
	return storesrank.Next(ctx, db, jobGroupRankGroup{job: job}, excludeID)
}
