package documentsstore

import (
	"context"
	"errors"

	resourcesdatabase "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	documentsrequests "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/requests"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListDocumentReqs(
	ctx context.Context,
	q ListDocumentReqsQuery,
) (resourcesdatabase.DataCount, []*documentsrequests.DocRequest, error) {
	tDocRequest := table.FivenetDocumentsRequests.AS("doc_request")
	condition := tDocRequest.DocumentID.EQ(mysql.Int64(q.DocumentID))

	countStmt := tDocRequest.
		SELECT(mysql.COUNT(tDocRequest.ID).AS("data_count.total")).
		FROM(tDocRequest).
		WHERE(condition)

	var count resourcesdatabase.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	pag, limit := q.Pagination.GetResponseWithPageSize(count.Total, 10)
	_ = pag
	if count.Total <= 0 {
		return count, []*documentsrequests.DocRequest{}, nil
	}

	tCreator := table.FivenetUser.AS("creator")
	stmt := tDocRequest.
		SELECT(
			tDocRequest.ID,
			tDocRequest.CreatedAt,
			tDocRequest.UpdatedAt,
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(
			tDocRequest.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tDocRequest.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(q.Pagination.GetOffset()).
		ORDER_BY(tDocRequest.ID.DESC()).
		LIMIT(limit)

	var reqs []*documentsrequests.DocRequest
	if err := stmt.QueryContext(ctx, s.db, &reqs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return resourcesdatabase.DataCount{}, nil, err
		}
	}

	return count, reqs, nil
}
