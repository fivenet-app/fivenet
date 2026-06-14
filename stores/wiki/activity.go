package wikistore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountPageActivity(ctx context.Context, pageID int64) (int64, error) {
	tPActivity := table.FivenetWikiPagesActivity.AS("page_activity")
	condition := tPActivity.PageID.EQ(mysql.Int64(pageID))

	countStmt := tPActivity.
		SELECT(
			mysql.COUNT(tPActivity.ID).AS("data_count.total"),
		).
		FROM(
			tPActivity,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListPageActivity(
	ctx context.Context,
	pageID int64,
	offset int64,
	limit int64,
) ([]*wikiactivity.PageActivity, error) {
	tPActivity := table.FivenetWikiPagesActivity.AS("page_activity")
	condition := tPActivity.PageID.EQ(mysql.Int64(pageID))

	tCreator := table.FivenetUser.AS("creator")

	stmt := tPActivity.
		SELECT(
			tPActivity.ID,
			tPActivity.CreatedAt,
			tPActivity.PageID,
			tPActivity.ActivityType,
			tPActivity.CreatorID,
			tPActivity.CreatorJob,
			tPActivity.Reason,
			tPActivity.Data,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tPActivity.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tPActivity.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(offset).
		ORDER_BY(
			tPActivity.ID.DESC(),
		).
		LIMIT(limit)

	activity := []*wikiactivity.PageActivity{}
	if err := stmt.QueryContext(ctx, s.db, &activity); err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *Store) AddPageActivity(
	ctx context.Context,
	tx qrm.DB,
	activity *wikiactivity.PageActivity,
) (int64, error) {
	stmt := table.FivenetWikiPagesActivity.
		INSERT(
			table.FivenetWikiPagesActivity.PageID,
			table.FivenetWikiPagesActivity.ActivityType,
			table.FivenetWikiPagesActivity.CreatorID,
			table.FivenetWikiPagesActivity.CreatorJob,
			table.FivenetWikiPagesActivity.Reason,
			table.FivenetWikiPagesActivity.Data,
		).
		VALUES(
			activity.GetPageId(),
			activity.GetActivityType(),
			activity.CreatorId,
			activity.CreatorJob,
			activity.Reason,
			activity.GetData(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Store) CountPageChildren(ctx context.Context, pageID int64) (int64, error) {
	tPage := table.FivenetWikiPages
	countStmt := tPage.
		SELECT(
			mysql.COUNT(tPage.ID).AS("data_count.total"),
		).
		FROM(tPage).
		WHERE(mysql.AND(
			tPage.ParentID.EQ(mysql.Int64(pageID)),
			tPage.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}
