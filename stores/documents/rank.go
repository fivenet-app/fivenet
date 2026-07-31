package documentsstore

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	storesrank "github.com/fivenet-app/fivenet/v2026/stores/internal/rank"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type TemplateOrderInfo struct {
	ID         int64
	CreatorJob string
	SortRank   string
}

type templateRankGroup struct {
	creatorJob string
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

func (g templateRankGroup) ListRanks(
	ctx context.Context,
	q qrm.DB,
	excludeID int64,
) ([]storesrank.Row, error) {
	tTemplate := table.FivenetDocumentsTemplates.AS("row")

	condition := mysql.AND(
		tTemplate.CreatorJob.EQ(mysql.String(g.creatorJob)),
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

	rows := []storesrank.Row{}
	if err := stmt.QueryContext(ctx, q, &rows); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return rows, nil
}

func (g templateRankGroup) UpdateRank(
	ctx context.Context,
	q qrm.DB,
	templateID int64,
	sortRank string,
) error {
	tTemplate := table.FivenetDocumentsTemplates

	_, err := tTemplate.
		UPDATE().
		SET(tTemplate.SortRank.SET(mysql.String(sortRank))).
		WHERE(mysql.AND(
			tTemplate.ID.EQ(mysql.Int64(templateID)),
			tTemplate.CreatorJob.EQ(mysql.String(g.creatorJob)),
			tTemplate.DeletedAt.IS_NULL(),
		)).
		LIMIT(1).
		ExecContext(ctx, q)

	return err
}

func (s *Store) NextTemplateGroupRank(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
) (string, error) {
	return storesrank.Next(ctx, q, templateRankGroup{creatorJob: creatorJob}, excludeID)
}

func (s *Store) InsertTemplateGroupRank(
	ctx context.Context,
	q qrm.DB,
	creatorJob string,
	excludeID int64,
	beforeID, afterID *int64,
) (string, error) {
	return storesrank.Insert(
		ctx,
		q,
		templateRankGroup{creatorJob: creatorJob},
		excludeID,
		beforeID,
		afterID,
		errorsdocuments.ErrTemplateNoPerms,
		errorsdocuments.ErrFailedQuery,
	)
}

func (s *Store) UpdateTemplateSortRank(
	ctx context.Context,
	q qrm.DB,
	templateID int64,
	creatorJob string,
	sortRank string,
) error {
	return templateRankGroup{creatorJob: creatorJob}.UpdateRank(ctx, q, templateID, sortRank)
}
