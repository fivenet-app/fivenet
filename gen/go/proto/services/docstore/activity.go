package docstore

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tDocActivity = table.FivenetDocumentsActivity
)

func (s *Server) AddDocumentActivity(ctx context.Context, tx qrm.DB, activitiy *documents.DocActivity) error {
	stmt := tDocActivity.
		INSERT(
			tDocActivity.DocumentID,
			tDocActivity.ActivityType,
			tDocActivity.CreatorID,
			tDocActivity.CreatorJob,
			tDocActivity.Reason,
			tDocActivity.Data,
		).
		VALUES(
			activitiy.DocumentId,
			activitiy.ActivityType,
			activitiy.CreatorId,
			activitiy.CreatorJob,
			activitiy.Reason,
			activitiy.Data,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Server) ListDocumentActivity(ctx context.Context, req *ListDocumentActivityRequest) (*ListDocumentActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, ErrFailedQuery
	}
	if !ok {
		return nil, ErrDocViewDenied
	}

	// TODO add audit log entry

	tDocActivity := table.FivenetDocumentsActivity.AS("doc_activity")

	condition := tDocActivity.DocumentID.EQ(jet.Uint64(req.DocumentId))

	countStmt := tDocActivity.
		SELECT(
			jet.COUNT(tDocActivity.ID).AS("datacount.totalcount"),
		).
		FROM(
			tDocActivity,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.Pagination.GetResponse()
	resp := &ListDocumentActivityResponse{
		Pagination: pag,
		Activity:   []*documents.DocActivity{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

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
			tCreator.Identifier,
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
		OFFSET(
			req.Pagination.Offset,
		).
		ORDER_BY(
			tDocActivity.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, err
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].Creator != nil {
			jobInfoFn(resp.Activity[i].Creator)
		}
	}

	return resp, nil
}
