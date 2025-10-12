package documents

import (
	context "context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) GetDocumentAccess(
	ctx context.Context,
	req *pbdocuments.GetDocumentAccessRequest,
) (*pbdocuments.GetDocumentAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocAccessViewDenied
	}

	access, err := s.getDocumentAccess(ctx, req.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range access.GetUsers() {
		if access.GetUsers()[i].GetUser() != nil {
			jobInfoFn(access.GetUsers()[i].GetUser())
		}
	}

	resp := &pbdocuments.GetDocumentAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(
	ctx context.Context,
	req *pbdocuments.SetDocumentAccessRequest,
) (*pbdocuments.SetDocumentAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_ACCESS,
	)
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

	if err := s.handleDocumentAccessChange(ctx, tx, req.GetDocumentId(), userInfo, req.GetAccess(), true); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.SetDocumentAccessResponse{}, nil
}

func (s *Server) getDocumentAccess(
	ctx context.Context,
	documentId int64,
) (*documents.DocumentAccess, error) {
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

func (s *Server) handleDocumentAccessChange(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	userInfo *userinfo.UserInfo,
	access *documents.DocumentAccess,
	addActivity bool,
) error {
	// Validate job access entries
	valid, err := s.access.Jobs.Validate(s.jobs, &access.Jobs, true)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !valid {
		return errorsdocuments.ErrDocAccessInvalid
	}

	changes, err := s.access.HandleAccessChanges(
		ctx,
		tx,
		documentId,
		access.GetJobs(),
		access.GetUsers(),
		nil,
	)
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
			CreatorJob:   userInfo.GetJob(),
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
