package documents

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tDPins = table.FivenetDocumentsPins

func (s *Server) ListDocumentPins(ctx context.Context, req *pbdocuments.ListDocumentPinsRequest) (*pbdocuments.ListDocumentPinsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tDocumentShort.ID.IN(
		tDPins.
			SELECT(
				tDPins.DocumentID,
			).
			FROM(tDPins).
			WHERE(tDPins.Job.EQ(jet.String(userInfo.Job))),
	)

	countStmt := s.listDocumentsQuery(
		condition, jet.ProjectionList{jet.COUNT(jet.DISTINCT(tDocumentShort.ID)).AS("data_count.total")}, userInfo)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, DocsDefaultPageSize)
	resp := &pbdocuments.ListDocumentPinsResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := s.listDocumentsQuery(condition, nil, userInfo).
		ORDER_BY(
			tDocumentShort.CreatedAt.DESC(),
			tDocumentShort.UpdatedAt.DESC(),
		).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tDocumentShort.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Documents {
		if resp.Documents[i].Creator != nil {
			jobInfoFn(resp.Documents[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Documents))

	return resp, nil
}

func (s *Server) ToggleDocumentPin(ctx context.Context, req *pbdocuments.ToggleDocumentPinRequest) (*pbdocuments.ToggleDocumentPinResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "ToggleDocumentPin",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrDocViewDenied
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
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return &pbdocuments.ToggleDocumentPinResponse{
		State: req.State,
	}, nil
}
