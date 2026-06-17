package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListDocumentActivity(
	ctx context.Context,
	q ListDocumentActivityQuery,
) (resourcesdatabase.DataCount, []*documentsactivity.DocActivity, error) {
	tDocActivity := table.FivenetDocumentsActivity.AS("doc_activity")
	condition := tDocActivity.DocumentID.EQ(mysql.Int64(q.DocumentID))
	if len(q.ActivityTypes) > 0 {
		ids := make([]mysql.Expression, len(q.ActivityTypes))
		for i := range q.ActivityTypes {
			ids[i] = mysql.Int32(int32(q.ActivityTypes[i]))
		}
		condition = condition.AND(tDocActivity.ActivityType.IN(ids...))
	}

	countStmt := tDocActivity.
		SELECT(mysql.COUNT(tDocActivity.ID).AS("data_count.total")).
		FROM(tDocActivity).
		WHERE(condition)

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	_, limit := q.Pagination.GetResponseWithPageSize(count.Total, 10)
	if count.Total <= 0 {
		return count, []*documentsactivity.DocActivity{}, nil
	}

	tCreator := table.FivenetUser.AS("creator")
	stmt := tDocActivity.
		SELECT(
			tDocActivity.ID,
			tDocActivity.CreatedAt,
			tDocActivity.DocumentID,
			tDocActivity.ActivityType,
			tDocActivity.CreatorID,
			tDocActivity.CreatorJob,
			tDocActivity.Reason,
			tDocActivity.Data,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tDocActivity.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tDocActivity.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(q.Pagination.GetOffset()).
		ORDER_BY(tDocActivity.ID.DESC()).
		LIMIT(limit)

	var activity []*documentsactivity.DocActivity
	if err := stmt.QueryContext(ctx, s.db, &activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, activity, nil
}
