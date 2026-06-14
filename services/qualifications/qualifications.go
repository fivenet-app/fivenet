package qualifications

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	permsqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	QualificationsPageSize = 10

	QualificationsLabelDefaultFormat = "%abbr%: %name%"
)

func (s *Server) ListQualifications(
	ctx context.Context,
	req *pbqualifications.ListQualificationsRequest,
) (*pbqualifications.ListQualificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	condition := mysql.Bool(true)
	if req.GetSearch() != "" {
		search := dbutils.PrepareForLikeSearch(req.GetSearch())
		condition = condition.AND(mysql.OR(
			tQuali.Abbreviation.LIKE(mysql.String(search)),
			tQuali.Title.LIKE(mysql.String(search)),
		))
	}

	if !userInfo.GetSuperuser() {
		accessExists := s.access.ACLAccessExistsCondition(
			tQuali.ID,
			userInfo,
			int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		)

		condition = condition.AND(mysql.AND(
			tQuali.DeletedAt.IS_NULL(),
			mysql.OR(
				tQuali.Public.IS_TRUE(),
				mysql.AND(
					tQuali.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
					tQuali.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
				),
				accessExists,
			),
		))
	}

	includePhoneNumber := false
	if fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(
		s.perms,
		userInfo,
	); err == nil {
		includePhoneNumber = fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		)
	}

	resp, err := s.store.ListQualifications(ctx, req, userInfo, condition, includePhoneNumber)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetQualifications() {
		if resp.GetQualifications()[i].GetCreator() != nil {
			jobInfoFn(resp.GetQualifications()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) GetQualification(
	ctx context.Context,
	req *pbqualifications.GetQualificationRequest,
) (*pbqualifications.GetQualificationResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	quali, err := s.store.GetQualificationShort(ctx, req.GetQualificationId(), nil, userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if !check && !userInfo.GetSuperuser() && !quali.GetPublic() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	request, err := s.getQualificationRequest(
		ctx,
		req.GetQualificationId(),
		userInfo.GetUserId(),
		userInfo,
	)
	if err != nil {
		return nil, errorsqualifications.ErrFailedQuery
	}

	canContent := false

	// If user's request is accepted or user has GRADE or higher perm to qualification, show content
	if request != nil {
		canContent = request.Status != nil &&
			request.GetStatus() >= qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED
	}

	canGrade, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !canContent {
		canContent = canGrade
	}

	// Allow content if the qualification has the exam mode enabled and the user has the access to take the qualification
	canTake, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_TAKE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if canTake &&
		quali.GetExamMode() == qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED {
		canContent = true
	}

	resp := &pbqualifications.GetQualificationResponse{}
	resp.Qualification, err = s.store.GetQualification(ctx, req.GetQualificationId(),
		tQuali.ID.EQ(mysql.Int64(req.GetQualificationId())), userInfo, canContent, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if resp.GetQualification() == nil || resp.GetQualification().GetId() <= 0 {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if resp.GetQualification().GetExam() == nil {
		resp.Qualification.Exam = &qualificationsexam.ExamQuestions{
			Questions: []*qualificationsexam.ExamQuestion{},
		}
	}
	if resp.GetQualification().GetExamSettings() == nil {
		resp.Qualification.ExamSettings = &qualificationsexam.QualificationExamSettings{
			Time:          durationpb.New(10 * time.Minute),
			AutoGradeMode: qualificationsexam.AutoGradeMode_AUTO_GRADE_MODE_STRICT,
		}
	}

	if resp.GetQualification().GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.GetQualification().GetCreator())
	}

	qualiAccess, err := s.GetQualificationAccess(
		ctx,
		&pbqualifications.GetQualificationAccessRequest{
			QualificationId: req.GetQualificationId(),
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if qualiAccess != nil {
		resp.Qualification.Access = qualiAccess.GetAccess()
	}

	files, err := s.fHandler.ListFilesForParentID(ctx, req.GetQualificationId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.Qualification.Files = files

	if canGrade && req.WithExam != nil && req.GetWithExam() {
		exam, err := s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), canGrade)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		resp.Qualification.Exam = exam
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
}

func (s *Server) CreateQualification(
	ctx context.Context,
	req *pbqualifications.CreateQualificationRequest,
) (*pbqualifications.CreateQualificationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	lastId, err := s.store.CreateQualification(ctx, tx, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobAccess := []*qualificationsaccess.QualificationJobAccess{}

	job := s.enricher.GetJobByName(userInfo.GetJob())
	if job != nil {
		highestGrade := int32(-1)
		if len(job.GetGrades()) > 0 {
			highestGrade = job.GetGrades()[len(job.GetGrades())-1].GetGrade()
		}

		jobAccess = append(jobAccess, &qualificationsaccess.QualificationJobAccess{
			TargetId:     lastId,
			Job:          job.GetName(),
			MinimumGrade: highestGrade,
			Access:       int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
		})
	}

	if _, err := s.access.ReplaceTargetAccess(
		ctx,
		tx,
		s.accessResolver,
		lastId,
		qualificationJobAccess(jobAccess),
		qualificationSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbqualifications.CreateQualificationResponse{
		QualificationId: lastId,
	}, nil
}

func (s *Server) UpdateQualification(
	ctx context.Context,
	req *pbqualifications.UpdateQualificationRequest,
) (*pbqualifications.UpdateQualificationResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.qualifications.id", req.GetQualification().GetId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualification().GetId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	oldQuali, err := s.store.GetQualification(ctx, req.GetQualification().GetId(),
		tQuali.ID.EQ(mysql.Int64(req.GetQualification().GetId())),
		userInfo, true, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	ownAccess, err := permsqualifications.QualificationsService.UpdateQualification.AccessTyped.Get(
		s.perms,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		ownAccess.StringList(),
		userInfo,
		oldQuali.GetCreatorJob(),
		oldQuali.GetCreator(),
	) {
		return nil, errorsqualifications.ErrQualiUpdateDenied
	}

	fields, err := permsqualifications.QualificationsService.UpdateQualification.FieldsTyped.Get(
		s.perms,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !fields.Contains(
		permsqualifications.QualificationsServiceUpdateQualificationFieldsPermValuePublic,
	) {
		req.Qualification.Public = oldQuali.GetPublic()
	}

	// Make sure that the qualification doesn't require itself
	if len(req.GetQualification().GetRequirements()) > 0 {
		for _, req := range req.GetQualification().GetRequirements() {
			if req.GetTargetQualificationId() == oldQuali.GetId() {
				return nil, errorsqualifications.ErrRequirementSelfRef
			}
		}
	}

	// A qualification can only be switched to published once
	if !oldQuali.GetDraft() && oldQuali.GetDraft() != req.GetQualification().GetDraft() {
		// Allow a super user to change the draft state
		if !userInfo.GetSuperuser() {
			req.Qualification.Draft = oldQuali.GetDraft()
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpdateQualification(ctx, tx, req.GetQualification()); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.GetQualification().GetAccess() != nil {
		if _, err := s.access.ReplaceTargetAccess(
			ctx,
			tx,
			s.accessResolver,
			req.GetQualification().GetId(),
			qualificationJobAccess(req.GetQualification().GetAccess().GetJobs()),
			qualificationSubjectAccessOptions,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.store.HandleQualificationRequirementsChanges(
		ctx,
		tx,
		req.GetQualification().GetId(),
		req.GetQualification().GetRequirements(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	files := []*file.File{}
	if req.Qualification.Files != nil {
		files = append(files, req.GetQualification().GetFiles()...)
	}

	if req.GetQualification().GetExam() != nil {
		questFiles, err := s.store.HandleExamQuestionsChanges(
			ctx,
			tx,
			req.GetQualification().GetId(),
			req.GetQualification().GetExam(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		files = append(files, questFiles...)
	}

	if _, _, err := s.fHandler.HandleFileChangesForParent(
		ctx,
		tx,
		req.GetQualification().GetId(),
		files,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbqualifications.UpdateQualificationResponse{
		QualificationId: req.GetQualification().GetId(),
	}, nil
}

func (s *Server) DeleteQualification(
	ctx context.Context,
	req *pbqualifications.DeleteQualificationRequest,
) (*pbqualifications.DeleteQualificationResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		if !userInfo.GetSuperuser() {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	quali, err := s.store.GetQualification(ctx, req.GetQualificationId(),
		tQuali.ID.EQ(mysql.Int64(req.GetQualificationId())), userInfo, true, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	fields, err := permsqualifications.QualificationsService.DeleteQualification.AccessTyped.Get(
		s.perms,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		quali.GetCreatorJob(),
		quali.GetCreator(),
	) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	var deletedAtTime *timestamp.Timestamp
	if quali.GetDeletedAt() == nil || !userInfo.GetSuperuser() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.DeleteQualification(
		ctx,
		tx,
		req.GetQualificationId(),
		deletedAtTime,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.DeleteQualificationResponse{}, nil
}
