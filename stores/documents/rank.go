package documentsstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const documentsTemplateRankKeyStep = utils.RankStep

type TemplateOrderInfo struct {
	ID         int64
	CreatorJob string
	SortRank   string
}

type TemplateRankRow struct {
	ID       int64
	SortRank string
}

func getTemplateRankBounds(
	rows []TemplateRankRow,
	beforeID, afterID *int64,
) (string, string, error) {
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
			return "", "", errorsdocuments.ErrTemplateNoPerms
		}
		lower := ""
		if idx > 0 {
			lower = rows[idx-1].SortRank
		}
		return lower, rows[idx].SortRank, nil

	case afterID != nil:
		idx := findIndex(*afterID)
		if idx < 0 {
			return "", "", errorsdocuments.ErrTemplateNoPerms
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

func (s *Store) GetTemplateOrderInfo(
	ctx context.Context,
	q qrm.DB,
	templateID int64,
) (*TemplateOrderInfo, error) {
	tTemplate := table.FivenetDocumentsTemplates.AS("template_order_info")

	stmt := tTemplate.
		SELECT(
			tTemplate.ID,
			tTemplate.CreatorJob,
			tTemplate.SortRank,
		).
		FROM(tTemplate).
		WHERE(mysql.AND(
			tTemplate.ID.EQ(mysql.Int64(templateID)),
			tTemplate.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	dest := &TemplateOrderInfo{}
	if err := stmt.QueryContext(ctx, q, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Store) listTemplateGroupRanks(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
) ([]TemplateRankRow, error) {
	tTemplate := table.FivenetDocumentsTemplates.AS("template_rank_row")

	condition := mysql.AND(
		tTemplate.CreatorJob.EQ(mysql.String(creatorJob)),
		tTemplate.DeletedAt.IS_NULL(),
	)
	if excludeID > 0 {
		condition = condition.AND(tTemplate.ID.NOT_EQ(mysql.Int64(excludeID)))
	}

	stmt := tTemplate.
		SELECT(
			tTemplate.ID,
			tTemplate.SortRank,
		).
		FROM(tTemplate).
		WHERE(condition).
		ORDER_BY(
			tTemplate.SortRank.ASC(),
			tTemplate.ID.ASC(),
		).
		FOR(mysql.UPDATE())

	rows := []TemplateRankRow{}
	if err := stmt.QueryContext(ctx, q, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return rows, nil
}

func (s *Store) rebalanceTemplateGroupRanks(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
) error {
	rows, err := s.listTemplateGroupRanks(ctx, q, creatorJob, excludeID)
	if err != nil {
		return err
	}

	tTemplate := table.FivenetDocumentsTemplates
	for idx, row := range rows {
		rank := utils.FormatRank(int64(idx+1) * documentsTemplateRankKeyStep)
		if _, err := tTemplate.
			UPDATE().
			SET(tTemplate.SortRank.SET(mysql.String(rank))).
			WHERE(mysql.AND(
				tTemplate.ID.EQ(mysql.Int64(row.ID)),
				tTemplate.CreatorJob.EQ(mysql.String(creatorJob)),
				tTemplate.DeletedAt.IS_NULL(),
			)).
			LIMIT(1).
			ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) NextTemplateGroupRank(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
) (string, error) {
	rows, err := s.listTemplateGroupRanks(ctx, q, creatorJob, excludeID)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return utils.FormatRank(documentsTemplateRankKeyStep), nil
	}

	return utils.NextRank(rows[len(rows)-1].SortRank)
}

func (s *Store) InsertTemplateGroupRank(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
	beforeID, afterID *int64,
) (string, error) {
	rows, err := s.listTemplateGroupRanks(ctx, q, creatorJob, excludeID)
	if err != nil {
		return "", err
	}

	if beforeID != nil && afterID != nil {
		return "", errors.New("before_id and after_id are mutually exclusive")
	}

	if len(rows) == 0 {
		if beforeID != nil || afterID != nil {
			return "", errorsdocuments.ErrTemplateNoPerms
		}
		return utils.FormatRank(documentsTemplateRankKeyStep), nil
	}

	lower, upper, err := getTemplateRankBounds(rows, beforeID, afterID)
	if err != nil {
		return "", err
	}

	rank, ok := utils.RankBetween(lower, upper)
	if ok {
		return rank, nil
	}

	if err := s.rebalanceTemplateGroupRanks(ctx, q, creatorJob, excludeID); err != nil {
		return "", err
	}

	rows, err = s.listTemplateGroupRanks(ctx, q, creatorJob, excludeID)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return utils.FormatRank(documentsTemplateRankKeyStep), nil
	}

	lower, upper, err = getTemplateRankBounds(rows, beforeID, afterID)
	if err != nil {
		return "", err
	}

	rank, ok = utils.RankBetween(lower, upper)
	if !ok {
		return "", errorsdocuments.ErrFailedQuery
	}

	return rank, nil
}

func (s *Store) UpdateTemplateSortRank(
	ctx context.Context,
	q qrm.DB,
	templateID int64,
	creatorJob string,
	sortRank string,
) error {
	tTemplate := table.FivenetDocumentsTemplates

	_, err := tTemplate.
		UPDATE(
			tTemplate.SortRank,
		).
		SET(
			tTemplate.SortRank.SET(mysql.String(sortRank)),
		).
		WHERE(mysql.AND(
			tTemplate.ID.EQ(mysql.Int64(templateID)),
			tTemplate.CreatorJob.EQ(mysql.String(creatorJob)),
			tTemplate.DeletedAt.IS_NULL(),
		)).
		LIMIT(1).
		ExecContext(ctx, q)

	return err
}
