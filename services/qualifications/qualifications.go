package qualifications

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	permsqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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

	condition := jet.Bool(true)

	if req.Search != nil && req.GetSearch() != "" {
		*req.Search = strings.TrimSpace(req.GetSearch())
		*req.Search = strings.ReplaceAll(req.GetSearch(), "%", "")
		*req.Search = strings.ReplaceAll(req.GetSearch(), " ", "%")
		*req.Search = "%" + req.GetSearch() + "%"
		condition = condition.AND(jet.OR(
			tQuali.Abbreviation.LIKE(jet.String(req.GetSearch())),
			tQuali.Title.LIKE(jet.String(req.GetSearch())),
		))
	}

	countStmt := s.listQualificationsQuery(
		condition,
		jet.ProjectionList{jet.COUNT(jet.DISTINCT(tQuali.ID)).AS("data_count.total")},
		userInfo,
	)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*qualifications.Qualification{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{tQuali.Draft.ASC()}
	if req.GetSort() != nil {
		var column jet.Column
		switch req.GetSort().GetColumn() {
		case "abbreviation":
			column = tQuali.Abbreviation
		case "id":
			fallthrough
		default:
			column = tQualiResults.ID
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tQualiResults.ID.DESC())
	}

	stmt := s.listQualificationsQuery(condition, nil, userInfo).
		OFFSET(req.GetPagination().GetOffset()).
		GROUP_BY(tQuali.ID).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetQualifications() {
		if resp.GetQualifications()[i].GetCreator() != nil {
			jobInfoFn(resp.GetQualifications()[i].GetCreator())
		}
	}

	resp.GetPagination().Update(len(resp.GetQualifications()))

	return resp, nil
}

func (s *Server) GetQualification(
	ctx context.Context,
	req *pbqualifications.GetQualificationRequest,
) (*pbqualifications.GetQualificationResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetQualification",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	quali, err := s.getQualificationShort(ctx, req.GetQualificationId(), nil, userInfo)
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
		qualifications.AccessLevel_ACCESS_LEVEL_GRADE,
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
		qualifications.AccessLevel_ACCESS_LEVEL_TAKE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if canTake &&
		quali.GetExamMode() == qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED {
		canContent = true
	}

	resp := &pbqualifications.GetQualificationResponse{}
	resp.Qualification, err = s.getQualification(ctx, req.GetQualificationId(),
		tQuali.ID.EQ(jet.Uint64(req.GetQualificationId())), userInfo, canContent)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if resp.GetQualification() == nil || resp.GetQualification().GetId() <= 0 {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if resp.GetQualification().GetExam() == nil {
		resp.Qualification.Exam = &qualifications.ExamQuestions{
			Questions: []*qualifications.ExamQuestion{},
		}
	}
	if resp.GetQualification().GetExamSettings() == nil {
		resp.Qualification.ExamSettings = &qualifications.QualificationExamSettings{
			Time:          durationpb.New(10 * time.Minute),
			AutoGradeMode: qualifications.AutoGradeMode_AUTO_GRADE_MODE_STRICT,
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
		exam, err := s.getExamQuestions(ctx, s.db, req.GetQualificationId(), canGrade)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		resp.Qualification.Exam = exam
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateQualification(
	ctx context.Context,
	req *pbqualifications.CreateQualificationRequest,
) (*pbqualifications.CreateQualificationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateQualification",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tQuali := table.FivenetQualifications
	stmt := tQuali.
		INSERT(
			tQuali.Job,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.ContentType,
			tQuali.Content,
			tQuali.CreatorID,
			tQuali.CreatorJob,
		).
		VALUES(
			userInfo.GetJob(),
			false,
			true,
			false,
			"",
			"",
			"",
			content.ContentType_CONTENT_TYPE_HTML,
			"",
			userInfo.GetUserId(),
			userInfo.GetJob(),
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobAccess := []*qualifications.QualificationJobAccess{}

	job := s.enricher.GetJobByName(userInfo.GetJob())
	if job != nil {
		highestGrade := int32(-1)
		if len(job.GetGrades()) > 0 {
			highestGrade = job.GetGrades()[len(job.GetGrades())-1].GetGrade()
		}

		jobAccess = append(jobAccess, &qualifications.QualificationJobAccess{
			TargetId:     uint64(lastId),
			Job:          job.GetName(),
			MinimumGrade: highestGrade,
			Access:       qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
		})
	}

	if _, err := s.access.HandleAccessChanges(ctx, tx, uint64(lastId), jobAccess, nil, nil); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbqualifications.CreateQualificationResponse{
		QualificationId: uint64(lastId),
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

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "UpdateQualification",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualification().GetId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	oldQuali, err := s.getQualification(ctx, req.GetQualification().GetId(),
		tQuali.ID.EQ(jet.Uint64(req.GetQualification().GetId())),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	ownAccess, err := s.perms.AttrStringList(
		userInfo,
		permsqualifications.QualificationsServicePerm,
		permsqualifications.QualificationsServiceUpdateQualificationPerm,
		permsqualifications.QualificationsServiceUpdateQualificationAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		ownAccess,
		userInfo,
		oldQuali.GetCreatorJob(),
		oldQuali.GetCreator(),
	) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	fields, err := s.perms.AttrStringList(
		userInfo,
		permsqualifications.QualificationsServicePerm,
		permsqualifications.QualificationsServiceUpdateQualificationPerm,
		permsqualifications.QualificationsServiceUpdateQualificationFieldsPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !fields.Contains("Public") {
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

	if req.Qualification.Description != nil {
		*req.Qualification.Description = strings.TrimSuffix(
			req.GetQualification().GetDescription(),
			"<br>",
		)
	}

	tQuali := table.FivenetQualifications
	stmt := tQuali.
		UPDATE(
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.ContentType,
			tQuali.Content,
			tQuali.DiscordSyncEnabled,
			tQuali.DiscordSettings,
			tQuali.ExamMode,
			tQuali.ExamSettings,
			tQuali.LabelSyncEnabled,
			tQuali.LabelSyncFormat,
		).
		SET(
			req.GetQualification().GetWeight(),
			req.GetQualification().GetClosed(),
			req.GetQualification().GetDraft(),
			req.GetQualification().GetPublic(),
			req.GetQualification().GetAbbreviation(),
			req.GetQualification().GetTitle(),
			req.GetQualification().GetDescription(),
			content.ContentType_CONTENT_TYPE_HTML,
			req.GetQualification().GetContent(),
			req.GetQualification().GetDiscordSyncEnabled(),
			req.GetQualification().GetDiscordSettings(),
			req.GetQualification().GetExamMode(),
			req.GetQualification().GetExamSettings(),
			req.GetQualification().GetLabelSyncEnabled(),
			req.GetQualification().GetLabelSyncFormat(),
		).
		WHERE(
			tQuali.ID.EQ(jet.Uint64(req.GetQualification().GetId())),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.GetQualification().GetAccess() != nil {
		if _, err := s.access.HandleAccessChanges(ctx, tx, req.GetQualification().GetId(), req.GetQualification().GetAccess().GetJobs(), nil, nil); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	files := []*file.File{}
	if req.Qualification.Files != nil {
		files = append(files, req.GetQualification().GetFiles()...)
	}

	if err := s.handleQualificationRequirementsChanges(ctx, tx, req.GetQualification().GetId(), req.GetQualification().GetRequirements()); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.GetQualification().GetExam() != nil {
		questFiles, err := s.handleExamQuestionsChanges(
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

	if _, _, err := s.fHandler.HandleFileChangesForParent(ctx, tx, req.GetQualification().GetId(), files); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

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

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualification",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		if !userInfo.GetSuperuser() {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.GetQualificationId(),
		tQuali.ID.EQ(jet.Uint64(req.GetQualificationId())), userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	fields, err := s.perms.AttrStringList(
		userInfo,
		permsqualifications.QualificationsServicePerm,
		permsqualifications.QualificationsServiceDeleteQualificationPerm,
		permsqualifications.QualificationsServiceDeleteQualificationAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, quali.GetCreatorJob(), quali.GetCreator()) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if quali.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	tQuali := table.FivenetQualifications
	stmt := tQuali.
		UPDATE(
			tQuali.DeletedAt,
		).
		SET(
			tQuali.DeletedAt.SET(deletedAtTime),
		).
		WHERE(
			tQuali.ID.EQ(jet.Uint64(req.GetQualificationId())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbqualifications.DeleteQualificationResponse{}, nil
}
