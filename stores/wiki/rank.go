package wikistore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	storesrank "github.com/fivenet-app/fivenet/v2026/stores/internal/rank"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type PageOrderInfo struct {
	ID        int64
	Job       string
	ParentID  *int64
	Startpage bool
	SortRank  string
}

type pageRankGroup struct {
	job       string
	parentID  *int64
	startpage bool
}

func (s *Store) GetPageOrderInfo(
	ctx context.Context,
	q qrm.DB,
	pageID int64,
) (*PageOrderInfo, error) {
	tPage := table.FivenetWikiPages.AS("page_order_info")

	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.Job,
			tPage.ParentID,
			tPage.Startpage,
			tPage.SortRank,
		).
		FROM(tPage).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(pageID)),
			tPage.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	dest := &PageOrderInfo{}
	if err := stmt.QueryContext(ctx, q, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (g pageRankGroup) ListRanks(
	ctx context.Context,
	q qrm.DB,
	excludeID int64,
) ([]storesrank.Row, error) {
	tPage := table.FivenetWikiPages.AS("row")

	condition := mysql.AND(
		tPage.Job.EQ(mysql.String(g.job)),
		tPage.DeletedAt.IS_NULL(),
	)
	if g.parentID == nil {
		condition = condition.AND(tPage.ParentID.IS_NULL())
		condition = condition.AND(tPage.Startpage.EQ(mysql.Bool(g.startpage)))
	} else {
		condition = condition.AND(tPage.ParentID.EQ(mysql.Int64(*g.parentID)))
	}
	if excludeID > 0 {
		condition = condition.AND(tPage.ID.NOT_EQ(mysql.Int64(excludeID)))
	}

	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.SortRank,
		).
		FROM(tPage).
		WHERE(condition).
		ORDER_BY(
			tPage.SortRank.ASC(),
			tPage.ID.ASC(),
		).
		FOR(mysql.UPDATE())

	rows := []storesrank.Row{}
	if err := stmt.QueryContext(ctx, q, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return rows, nil
}

func (g pageRankGroup) UpdateRank(
	ctx context.Context,
	q qrm.DB,
	pageID int64,
	sortRank string,
) error {
	tPage := table.FivenetWikiPages

	_, err := tPage.
		UPDATE().
		SET(tPage.SortRank.SET(mysql.String(sortRank))).
		WHERE(mysql.AND(
			tPage.ID.EQ(mysql.Int64(pageID)),
			tPage.Job.EQ(mysql.String(g.job)),
			tPage.DeletedAt.IS_NULL(),
		)).
		LIMIT(1).
		ExecContext(ctx, q)

	return err
}

func (s *Store) NextPageGroupRank(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
) (string, error) {
	return storesrank.Next(
		ctx,
		q,
		pageRankGroup{job: job, parentID: parentID, startpage: startpage},
		excludeID,
	)
}

func (s *Store) InsertPageGroupRank(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
	beforeID, afterID *int64,
) (string, error) {
	return storesrank.Insert(
		ctx,
		q,
		pageRankGroup{job: job, parentID: parentID, startpage: startpage},
		excludeID,
		beforeID,
		afterID,
		errorswiki.ErrPageNotFound,
		errorswiki.ErrFailedQuery,
	)
}
