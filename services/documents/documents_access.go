package documents

import (
	context "context"
	"errors"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var documentSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_COMMENT),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS),
		int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

const documentAccessEntryLimit = 20

func canUserAccessDocument(
	ctx context.Context,
	v2 *access.SubjectObjectAccess,
	targetID int64,
	userInfo *userinfo.UserInfo,
	level documentsaccess.AccessLevel,
) (bool, error) {
	if userInfo.GetSuperuser() {
		return true, nil
	}

	allowed, err := v2.CanUserAccessTarget(ctx, targetID, userInfo, int32(level))
	if err != nil {
		return false, err
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
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

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

	docAccess, err := s.getDocumentAccess(ctx, req.GetDocumentId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	requiredAccess, err := s.getDocumentRequiredAccess(ctx, req.GetDocumentId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	docAccess, err = access.NormalizeAccess(
		docAccess,
		requiredAccess,
		nil,
		documentAccessEntryLimit,
	)
	if err != nil {
		if isAccessEntryLimitError(err) {
			return nil, errorsdocuments.ErrDocRequiredAccessTemplate
		}
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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

func (s *Server) SetDocumentAccess(
	ctx context.Context,
	req *pbdocuments.SetDocumentAccessRequest,
) (*pbdocuments.SetDocumentAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

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

func (s *Server) getDocumentAccess(
	ctx context.Context,
	documentId int64,
) (*documentsaccess.DocumentAccess, error) {
	return s.subjectAccess.ListTargetAccess(ctx, s.db, documentId, documentSubjectAccessOptions)
}

func (s *Server) getDocumentRequiredAccess(
	ctx context.Context,
	documentId int64,
	userInfo *userinfo.UserInfo,
) (*resourcesaccess.Access, error) {
	doc, err := s.getDocument(ctx, tDocument.ID.EQ(mysql.Int64(documentId)), userInfo, false)
	if err != nil {
		return nil, err
	}
	if doc.GetTemplateId() <= 0 {
		return nil, nil
	}

	tmpl, err := s.getTemplate(ctx, doc.GetTemplateId())
	if err != nil {
		return nil, err
	}
	if tmpl.GetContentAccess() == nil || tmpl.GetContentAccess().IsEmpty() {
		return nil, nil
	}

	return tmpl.GetContentAccess(), nil
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

	// Validate job access entries
	valid, err := access.ValidateJobAccessEntries(s.jobs, &docAccess.Jobs, true)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !valid {
		return errorsdocuments.ErrDocAccessInvalid
	}

	requiredAccess, err := s.getDocumentRequiredAccess(ctx, documentId, userInfo)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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

	docAccess, err = access.NormalizeAccess(
		docAccess,
		requiredAccess,
		fallbackAccess,
		documentAccessEntryLimit,
	)
	if err != nil {
		if isAccessEntryLimitError(err) {
			return errorsdocuments.ErrDocRequiredAccessTemplate
		}
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if documentsaccess.DocumentAccessHasDuplicates(docAccess) {
		return errorsdocuments.ErrDocAccessDuplicate
	}

	changes, err := s.subjectAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		documentId,
		docAccess,
		documentSubjectAccessOptions,
	)
	if err != nil {
		if dbutils.IsDuplicateError(err) {
			return errswrap.NewError(err, errorsdocuments.ErrDocAccessDuplicate)
		}
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if addActivity && !changes.IsEmpty() {
		if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
			DocumentId:   documentId,
			ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Data: &documentsactivity.DocActivityData{
				Data: &documentsactivity.DocActivityData_AccessUpdated{
					AccessUpdated: &documentsactivity.DocAccessUpdated{
						Jobs: &documentsactivity.DocAccessJobsDiff{
							ToCreate: changes.Jobs.ToCreate,
							ToUpdate: changes.Jobs.ToUpdate,
							ToDelete: changes.Jobs.ToDelete,
						},
						Users: &documentsactivity.DocAccessUsersDiff{
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

func isAccessEntryLimitError(err error) bool {
	var limitErr *access.AccessEntryLimitError
	return errors.As(err, &limitErr)
}
