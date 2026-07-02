package documents

import (
	context "context"
	"errors"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const documentAccessEntryLimit = 15

func canUserAccessDocument(
	ctx context.Context,
	access *access.SubjectObjectAccess,
	targetID int64,
	userInfo *userinfo.UserInfo,
	level documentsaccess.AccessLevel,
) (bool, error) {
	if userInfo.GetJobAdmin() {
		return true, nil
	}

	allowed, err := access.CanUserAccessTarget(ctx, targetID, userInfo, int32(level))
	if err != nil {
		return false, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if allowed {
		return true, nil
	}

	return false, nil
}

func (s *Server) canUserAccessDocument(
	ctx context.Context,
	targetID int64,
	userInfo *userinfo.UserInfo,
	level documentsaccess.AccessLevel,
) (bool, error) {
	return canUserAccessDocument(ctx, s.subjectAccess, targetID, userInfo, level)
}

func (s *Server) canUserAccessDocumentIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	level documentsaccess.AccessLevel,
	targetIDs ...int64,
) ([]int64, error) {
	out := make([]int64, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		allowed, err := s.canUserAccessDocument(ctx, targetID, userInfo, level)
		if err != nil {
			return nil, err
		}
		if allowed {
			out = append(out, targetID)
		}
	}

	return out, nil
}

func (s *Server) canUserAccessDocuments(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	level documentsaccess.AccessLevel,
	targetIDs ...int64,
) (bool, error) {
	out, err := s.canUserAccessDocumentIDs(ctx, userInfo, level, targetIDs...)
	return len(out) == len(targetIDs), err
}

func (s *Server) GetDocumentAccess(
	ctx context.Context,
	req *pbdocuments.GetDocumentAccessRequest,
) (*pbdocuments.GetDocumentAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocAccessViewDenied
	}

	docAccess, err := s.store.GetDocumentAccess(ctx, req.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	docAccess, err = s.normalizeDocumentAccess(ctx, req.GetDocumentId(), userInfo, docAccess, nil)
	if err != nil {
		return nil, err
	}

	for i := range docAccess.GetJobs() {
		s.enricher.EnrichJobInfo(docAccess.GetJobs()[i])
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range docAccess.GetUsers() {
		if docAccess.GetUsers()[i].GetUser() != nil {
			jobInfoFn(docAccess.GetUsers()[i].GetUser())
		}
	}

	resp := &pbdocuments.GetDocumentAccessResponse{
		Access: docAccess,
	}

	return resp, nil
}

func (s *Server) normalizeDocumentAccess(
	ctx context.Context,
	documentID int64,
	userInfo *userinfo.UserInfo,
	docAccess *documentsaccess.DocumentAccess,
	fallback *resourcesaccess.Access,
) (*documentsaccess.DocumentAccess, error) {
	sanitizedDocAccess, err := access.SanitizeAccessJobs(s.jobs, docAccess)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	requiredAccess, err := s.store.GetDocumentRequiredAccess(ctx, documentID, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	sanitizedRequiredAccess, err := access.SanitizeAccessJobs(s.jobs, requiredAccess)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	sanitizedRequiredAccess = access.NormalizeRequiredAccessFloors(sanitizedRequiredAccess)

	normalizedAccess, err := access.NormalizeAccess(
		sanitizedDocAccess,
		sanitizedRequiredAccess,
		fallback,
		documentAccessEntryLimit,
	)
	if err != nil {
		if _, ok := errors.AsType[access.AccessEntryLimitError](err); ok {
			return nil, errorsdocuments.ErrDocRequiredAccessTemplate
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return normalizedAccess, nil
}

func (s *Server) SetDocumentAccess(
	ctx context.Context,
	req *pbdocuments.SetDocumentAccessRequest,
) (*pbdocuments.SetDocumentAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS,
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

	if err := s.handleDocumentAccessChange(
		ctx,
		tx,
		req.GetDocumentId(),
		userInfo,
		req.GetAccess(),
		true,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.SetDocumentAccessResponse{}, nil
}

func (s *Server) handleDocumentAccessChange(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	userInfo *userinfo.UserInfo,
	docAccess *documentsaccess.DocumentAccess,
	addActivity bool,
) error {
	if docAccess == nil {
		docAccess = &documentsaccess.DocumentAccess{}
	}

	fallbackAccess := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			},
		},
	}

	normalizedAccess, err := s.normalizeDocumentAccess(
		ctx,
		documentId,
		userInfo,
		docAccess,
		fallbackAccess,
	)
	if err != nil {
		return err
	}

	if err := s.store.UpdateDocumentAccess(
		ctx,
		tx,
		documentId,
		userInfo,
		normalizedAccess,
		addActivity,
	); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.subjectAccess.RefreshTargetVisibility(ctx, tx, documentId); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return nil
}
