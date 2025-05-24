package documents

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) GetDocumentAccess(ctx context.Context, req *pbdocuments.GetDocumentAccessRequest) (*pbdocuments.GetDocumentAccessResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocAccessViewDenied
	}

	access, err := s.getDocumentAccess(ctx, req.DocumentId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	for i := range access.Jobs {
		s.enricher.EnrichJobInfo(access.Jobs[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range access.Users {
		if access.Users[i].User != nil {
			jobInfoFn(access.Users[i].User)
		}
	}

	resp := &pbdocuments.GetDocumentAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *pbdocuments.SetDocumentAccessRequest) (*pbdocuments.SetDocumentAccessResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "SetDocumentAccess",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_ACCESS)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocAccessEditDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.handleDocumentAccessChange(ctx, tx, req.DocumentId, userInfo, req.Access, true); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.SetDocumentAccessResponse{}, nil
}

func (s *Server) getDocumentAccess(ctx context.Context, documentId uint64) (*documents.DocumentAccess, error) {
	jobAccess, err := s.access.Jobs.List(ctx, s.db, documentId)
	if err != nil {
		return nil, err
	}

	userAccess, err := s.access.Users.List(ctx, s.db, documentId)
	if err != nil {
		return nil, err
	}

	return &documents.DocumentAccess{
		Jobs:  jobAccess,
		Users: userAccess,
	}, nil
}

func (s *Server) handleDocumentAccessChange(ctx context.Context, tx qrm.DB, documentId uint64, userInfo *userinfo.UserInfo, access *documents.DocumentAccess, addActivity bool) error {
	changes, err := s.access.HandleAccessChanges(ctx, tx, documentId, access.Jobs, access.Users, nil)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return errswrap.NewError(err, errorsdocuments.ErrDocAccessDuplicate)
		}
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if addActivity && !changes.IsEmpty() {
		if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   documentId,
			ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
			Data: &documents.DocActivityData{
				Data: &documents.DocActivityData_AccessUpdated{
					AccessUpdated: &documents.DocAccessUpdated{
						Jobs: &documents.DocAccessJobsDiff{
							ToCreate: changes.Jobs.ToCreate,
							ToUpdate: changes.Jobs.ToUpdate,
							ToDelete: changes.Jobs.ToDelete,
						},
						Users: &documents.DocAccessUsersDiff{
							ToCreate: changes.Users.ToCreate,
							ToUpdate: changes.Users.ToUpdate,
							ToDelete: changes.Users.ToDelete,
						},
					},
				},
			},
		}); err != nil {
			return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return nil
}
