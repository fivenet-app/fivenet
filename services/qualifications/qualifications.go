package qualifications

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	QualificationsPageSize = 10

	QualificationsLabelDefaultFormat = "%abbr%: %name%"
)

func (s *Server) ListQualifications(ctx context.Context, req *pbqualifications.ListQualificationsRequest) (*pbqualifications.ListQualificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)

	if req.Search != nil && *req.Search != "" {
		*req.Search = strings.TrimSpace(*req.Search)
		*req.Search = strings.ReplaceAll(*req.Search, "%", "")
		*req.Search = strings.ReplaceAll(*req.Search, " ", "%")
		*req.Search = "%" + *req.Search + "%"
		condition = condition.AND(jet.OR(
			tQuali.Abbreviation.LIKE(jet.String(*req.Search)),
			tQuali.Title.LIKE(jet.String(*req.Search)),
		))
	}

	countStmt := s.listQualificationsQuery(
		condition, jet.ProjectionList{jet.COUNT(jet.DISTINCT(tQuali.ID)).AS("data_count.total")}, userInfo)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*qualifications.Qualification{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "abbreviation":
			column = tQuali.Abbreviation
		case "id":
			fallthrough
		default:
			column = tQualiResults.ID
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tQualiResults.ID.DESC())
	}

	stmt := s.listQualificationsQuery(condition, nil, userInfo).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tQuali.ID).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Qualifications {
		if resp.Qualifications[i].Creator != nil {
			jobInfoFn(resp.Qualifications[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Qualifications))

	return resp, nil
}

func (s *Server) GetQualification(ctx context.Context, req *pbqualifications.GetQualificationRequest) (*pbqualifications.GetQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetQualification",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	request, err := s.getQualificationRequest(ctx, req.QualificationId, userInfo.UserId, userInfo)
	if err != nil {
		return nil, errorsqualifications.ErrFailedQuery
	}

	canContent := false

	// If user's request is accepted or user has GRADE or higher perm to qualification, show content
	if request != nil {
		canContent = request.Status != nil && *request.Status >= qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED
	}

	canGrade, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !canContent {
		canContent = canGrade
	}

	// Allow content if the qualification has the exam mode enabled and the user has the access to take the qualification
	canTake, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if canTake {
		quali, err := s.getQualificationShort(ctx, req.QualificationId, nil, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if quali.ExamMode == qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_ENABLED {
			canContent = true
		}
	}

	resp := &pbqualifications.GetQualificationResponse{}
	resp.Qualification, err = s.getQualification(ctx, req.QualificationId,
		tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo, canContent)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if resp.Qualification == nil || resp.Qualification.Id <= 0 {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if resp.Qualification.Exam == nil {
		resp.Qualification.Exam = &qualifications.ExamQuestions{
			Questions: []*qualifications.ExamQuestion{},
		}
	}
	if resp.Qualification.ExamSettings == nil {
		resp.Qualification.ExamSettings = &qualifications.QualificationExamSettings{
			Time: durationpb.New(10 * time.Minute),
		}
	}

	if resp.Qualification.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.Qualification.Creator)
	}

	qualiAccess, err := s.GetQualificationAccess(ctx, &pbqualifications.GetQualificationAccessRequest{
		QualificationId: req.QualificationId,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if qualiAccess != nil {
		resp.Qualification.Access = qualiAccess.Access
	}

	if canGrade && req.WithExam != nil && *req.WithExam {
		exam, err := s.getExamQuestions(ctx, s.db, req.QualificationId, canGrade)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		resp.Qualification.Exam = exam
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateQualification(ctx context.Context, req *pbqualifications.CreateQualificationRequest) (*pbqualifications.CreateQualificationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateQualification",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Field Permission Check
	fields, err := s.perms.AttrStringList(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceCreateQualificationPerm, permsqualifications.QualificationsServiceCreateQualificationFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !fields.Contains("Public") {
		req.Qualification.Public = false
	}

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
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.Content,
			tQuali.DiscordSyncEnabled,
			tQuali.DiscordSettings,
			tQuali.ExamMode,
			tQuali.ExamSettings,
			tQuali.CreatorID,
			tQuali.CreatorJob,
		).
		VALUES(
			userInfo.Job,
			req.Qualification.Weight,
			req.Qualification.Closed,
			req.Qualification.Public,
			req.Qualification.Abbreviation,
			req.Qualification.Title,
			req.Qualification.Description,
			req.Qualification.Content,
			req.Qualification.DiscordSyncEnabled,
			req.Qualification.DiscordSettings,
			req.Qualification.ExamMode,
			req.Qualification.ExamSettings,
			userInfo.UserId,
			userInfo.Job,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Access != nil {
		if _, err := s.access.HandleAccessChanges(ctx, tx, uint64(lastId), req.Qualification.Access.Jobs, nil, nil); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.handleQualificationRequirementsChanges(ctx, tx, uint64(lastId), req.Qualification.Requirements); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Exam != nil && req.Qualification.Exam.Questions != nil {
		if err := s.handleExamQuestionsChanges(ctx, tx, uint64(lastId), req.Qualification.Exam); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if _, _, err := s.qualiFileHandler.HandleFileChangesForParent(ctx, tx, req.Qualification.Id, req.Qualification.Files); err != nil {
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

func (s *Server) UpdateQualification(ctx context.Context, req *pbqualifications.UpdateQualificationRequest) (*pbqualifications.UpdateQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.Qualification.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "UpdateQualification",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Qualification.Id, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualification(ctx, req.Qualification.Id,
		tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	ownAccess, err := s.perms.AttrStringList(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceUpdateQualificationPerm, permsqualifications.QualificationsServiceUpdateQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasAccess(ownAccess, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	fields, err := s.perms.AttrStringList(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceCreateQualificationPerm, permsqualifications.QualificationsServiceCreateQualificationFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !fields.Contains("Public") {
		req.Qualification.Public = quali.Public
	}

	// Make sure that the qualification doesn't require itself
	if len(req.Qualification.Requirements) > 0 {
		for _, req := range req.Qualification.Requirements {
			if req.TargetQualificationId == quali.Id {
				return nil, errorsqualifications.ErrRequirementSelfRef
			}
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Description != nil {
		*req.Qualification.Description = strings.TrimSuffix(*req.Qualification.Description, "<br>")
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
			tQuali.Content,
			tQuali.DiscordSyncEnabled,
			tQuali.DiscordSettings,
			tQuali.ExamMode,
			tQuali.ExamSettings,
			tQuali.LabelSyncEnabled,
			tQuali.LabelSyncFormat,
		).
		SET(
			req.Qualification.Weight,
			req.Qualification.Closed,
			req.Qualification.Public,
			req.Qualification.Abbreviation,
			req.Qualification.Title,
			req.Qualification.Description,
			req.Qualification.Content,
			req.Qualification.DiscordSyncEnabled,
			req.Qualification.DiscordSettings,
			req.Qualification.ExamMode,
			req.Qualification.ExamSettings,
			req.Qualification.LabelSyncEnabled,
			req.Qualification.LabelSyncFormat,
		).
		WHERE(
			tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Access != nil {
		if _, err := s.access.HandleAccessChanges(ctx, tx, req.Qualification.Id, req.Qualification.Access.Jobs, nil, nil); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.handleQualificationRequirementsChanges(ctx, tx, req.Qualification.Id, req.Qualification.Requirements); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Exam != nil {
		if err := s.handleExamQuestionsChanges(ctx, tx, req.Qualification.Id, req.Qualification.Exam); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbqualifications.UpdateQualificationResponse{
		QualificationId: req.Qualification.Id,
	}, nil
}

func (s *Server) DeleteQualification(ctx context.Context, req *pbqualifications.DeleteQualificationRequest) (*pbqualifications.DeleteQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualification",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		if !userInfo.Superuser {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.QualificationId,
		tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	fields, err := s.perms.AttrStringList(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceDeleteQualificationPerm, permsqualifications.QualificationsServiceDeleteQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !access.CheckIfHasAccess(fields, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if quali.DeletedAt != nil && userInfo.Superuser {
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
			tQuali.ID.EQ(jet.Uint64(req.QualificationId)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbqualifications.DeleteQualificationResponse{}, nil
}
