package wiki

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const wikiPageRankKeyStep = utils.RankStep

type pageOrderInfo struct {
	ID        int64
	Job       string
	ParentID  *int64
	Startpage bool
	SortRank  string
}

type pageRankRow struct {
	ID       int64
	SortRank string
}

func getPageRankBounds(rows []pageRankRow, beforeID, afterID *int64) (string, string, error) {
	if beforeID != nil && afterID != nil {
		return "", "", errors.New("before_id and after_id are mutually exclusive")
	}
	if len(rows) == 0 {
		return "", "", nil
	}

	findIndex := func(id int64) int {
		for idx, row := range rows {
			if row.ID == id {
				return idx
			}
		}
		return -1
	}

	switch {
	case beforeID != nil:
		idx := findIndex(*beforeID)
		if idx < 0 {
			return "", "", errorswiki.ErrPageNotFound
		}
		lower := ""
		if idx > 0 {
			lower = rows[idx-1].SortRank
		}
		return lower, rows[idx].SortRank, nil
	case afterID != nil:
		idx := findIndex(*afterID)
		if idx < 0 {
			return "", "", errorswiki.ErrPageNotFound
		}
		upper := ""
		if idx < len(rows)-1 {
			upper = rows[idx+1].SortRank
		}
		return rows[idx].SortRank, upper, nil
	default:
		return rows[len(rows)-1].SortRank, "", nil
	}
}

func (s *Server) getPageOrderInfo(
	ctx context.Context,
	q qrm.DB,
	pageID int64,
) (*pageOrderInfo, error) {
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

	dest := &pageOrderInfo{}
	if err := stmt.QueryContext(ctx, q, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) listPageGroupRanks(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
) ([]pageRankRow, error) {
	tPage := table.FivenetWikiPages.AS("page_rank_row")

	condition := mysql.AND(
		tPage.Job.EQ(mysql.String(job)),
		tPage.DeletedAt.IS_NULL(),
	)
	if parentID == nil {
		condition = condition.AND(tPage.ParentID.IS_NULL())
		condition = condition.AND(tPage.Startpage.EQ(mysql.Bool(startpage)))
	} else {
		condition = condition.AND(tPage.ParentID.EQ(mysql.Int64(*parentID)))
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

	rows := []pageRankRow{}
	if err := stmt.QueryContext(ctx, q, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	return rows, nil
}

func (s *Server) rebalancePageGroupRanks(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
) error {
	rows, err := s.listPageGroupRanks(ctx, q, job, parentID, startpage, excludeID)
	if err != nil {
		return err
	}

	tPage := table.FivenetWikiPages
	for idx, row := range rows {
		rank := utils.FormatRank(int64(idx+1) * wikiPageRankKeyStep)
		if _, err := tPage.
			UPDATE().
			SET(tPage.SortRank.SET(mysql.String(rank))).
			WHERE(mysql.AND(
				tPage.ID.EQ(mysql.Int64(row.ID)),
				tPage.Job.EQ(mysql.String(job)),
				tPage.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, q); err != nil {
			return errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	return nil
}

func (s *Server) nextPageGroupRank(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
) (string, error) {
	rows, err := s.listPageGroupRanks(ctx, q, job, parentID, startpage, excludeID)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return utils.FormatRank(wikiPageRankKeyStep), nil
	}

	return utils.NextRank(rows[len(rows)-1].SortRank)
}

func (s *Server) insertPageGroupRank(
	ctx context.Context,
	q qrm.DB,
	job string,
	parentID *int64,
	startpage bool,
	excludeID int64,
	beforeID, afterID *int64,
) (string, error) {
	rows, err := s.listPageGroupRanks(ctx, q, job, parentID, startpage, excludeID)
	if err != nil {
		return "", err
	}

	if beforeID != nil && afterID != nil {
		return "", errors.New("before_id and after_id are mutually exclusive")
	}

	if len(rows) == 0 {
		if beforeID != nil || afterID != nil {
			return "", errorswiki.ErrPageNotFound
		}
		return utils.FormatRank(wikiPageRankKeyStep), nil
	}

	lower, upper, err := getPageRankBounds(rows, beforeID, afterID)
	if err != nil {
		return "", err
	}

	rank, ok := utils.RankBetween(lower, upper)
	if ok {
		return rank, nil
	}

	if err := s.rebalancePageGroupRanks(ctx, q, job, parentID, startpage, excludeID); err != nil {
		return "", err
	}

	rows, err = s.listPageGroupRanks(ctx, q, job, parentID, startpage, excludeID)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return utils.FormatRank(wikiPageRankKeyStep), nil
	}

	lower, upper, err = getPageRankBounds(rows, beforeID, afterID)
	if err != nil {
		return "", err
	}

	rank, ok = utils.RankBetween(lower, upper)
	if !ok {
		return "", errorswiki.ErrFailedQuery
	}

	return rank, nil
}
