package docstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tDPins = table.FivenetDocumentsPins
)

func (s *Server) ListDocumentPins(ctx context.Context, req *ListDocumentPinsRequest) (*ListDocumentPinsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tDocumentShort.ID.IN(
		tDPins.
			SELECT(
				tDPins.DocumentID,
			).
			FROM(tDPins).
			WHERE(tDPins.Job.EQ(jet.String(userInfo.Job))),
	)

	countStmt :=
		s.listDocumentsQuery(
			condition, jet.ProjectionList{jet.COUNT(jet.DISTINCT(tDocumentShort.ID)).AS("datacount.totalcount")}, userInfo)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, DocsDefaultPageSize)
	resp := &ListDocumentPinsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.listDocumentsQuery(condition, nil, userInfo).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tDocumentShort.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Documents); i++ {
		if resp.Documents[i].Creator != nil {
			jobInfoFn(resp.Documents[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Documents))

	return resp, nil
}

func (s *Server) ToggleDocumentPin(ctx context.Context, req *ToggleDocumentPinRequest) (*ToggleDocumentPinResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "ToggleDocumentPin",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrDocViewDenied
	}

	if req.State {
		stmt := tDPins.
			INSERT(
				tDPins.DocumentID,
				tDPins.Job,
				tDPins.CreatorID,
			).
			VALUES(
				req.DocumentId,
				userInfo.Job,
				userInfo.UserId,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
			}
		}
	} else {
		stmt := tDPins.
			DELETE().
			WHERE(jet.AND(
				tDPins.DocumentID.EQ(jet.Uint64(req.DocumentId)),
				tDPins.Job.EQ(jet.String(userInfo.Job)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	return &ToggleDocumentPinResponse{
		State: req.State,
	}, nil
}
