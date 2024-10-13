package qualifications

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	permsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/perms"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

const QualificationsPageSize = 10

var tQuali = table.FivenetQualifications.AS("qualification")

type Server struct {
	QualificationsServiceServer

	logger   *zap.Logger
	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	notif    notifi.INotifi
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher *mstlystcdata.UserAwareEnricher
	Audit             audit.IAuditer
	Config            *config.Config
	Notif             notifi.INotifi
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("jobs"),

		db:       p.DB,
		ps:       p.Perms,
		enricher: p.UserAwareEnricher,
		aud:      p.Audit,
		notif:    p.Notif,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterQualificationsServiceServer(srv, s)
}

func (s *Server) ListQualifications(ctx context.Context, req *ListQualificationsRequest) (*ListQualificationsResponse, error) {
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
		condition, jet.ProjectionList{jet.COUNT(jet.DISTINCT(tQuali.ID)).AS("datacount.totalcount")}, userInfo)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*qualifications.Qualification{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.listQualificationsQuery(condition, nil, userInfo).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tQuali.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Qualifications); i++ {
		if resp.Qualifications[i].Creator != nil {
			jobInfoFn(resp.Qualifications[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Qualifications))

	return resp, nil
}

func (s *Server) GetQualification(ctx context.Context, req *GetQualificationRequest) (*GetQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
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

	canGrade, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !canContent {
		canContent = canGrade
	}

	resp := &GetQualificationResponse{}
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

	qualiAccess, err := s.GetQualificationAccess(ctx, &GetQualificationAccessRequest{
		QualificationId: req.QualificationId,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if qualiAccess != nil {
		resp.Qualification.Access = qualiAccess.Access
	}

	if req.WithExam != nil && *req.WithExam {
		exam, err := s.getExamQuestions(ctx, req.QualificationId, canGrade)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		resp.Qualification.Exam = exam
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateQualification(ctx context.Context, req *CreateQualificationRequest) (*CreateQualificationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
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
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.Content,
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
			req.Qualification.Abbreviation,
			req.Qualification.Title,
			req.Qualification.Description,
			req.Qualification.Content,
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

	if err := s.handleQualificationAccessChanges(ctx, tx, qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, uint64(lastId), req.Qualification.Access); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.handleQualificationRequirementsChanges(ctx, tx, uint64(lastId), req.Qualification.Requirements); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Exam != nil && req.Qualification.Exam.Questions != nil {
		if err := s.handleExamQuestionsChanges(ctx, tx, uint64(lastId), req.Qualification.Exam); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateQualificationResponse{
		QualificationId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateQualification(ctx context.Context, req *UpdateQualificationRequest) (*UpdateQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.Qualification.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "UpdateQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.Qualification.Id, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	var onlyUpdateAccess bool
	if !check && !userInfo.SuperUser {
		onlyUpdateAccess, err = s.checkIfUserHasAccessToQuali(ctx, req.Qualification.Id, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errorsqualifications.ErrFailedQuery
		}
		if !onlyUpdateAccess {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.Qualification.Id,
		tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceUpdateQualificationPerm, permsqualifications.QualificationsServiceUpdateQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if !onlyUpdateAccess {
		if req.Qualification.Description != nil {
			*req.Qualification.Description = strings.TrimSuffix(*req.Qualification.Description, "<br>")
		}

		tQuali := table.FivenetQualifications
		stmt := tQuali.
			UPDATE(
				tQuali.Weight,
				tQuali.Closed,
				tQuali.Abbreviation,
				tQuali.Title,
				tQuali.Description,
				tQuali.Content,
				tQuali.DiscordSettings,
				tQuali.ExamMode,
				tQuali.ExamSettings,
			).
			SET(
				req.Qualification.Weight,
				req.Qualification.Closed,
				req.Qualification.Abbreviation,
				req.Qualification.Title,
				req.Qualification.Description,
				req.Qualification.Content,
				req.Qualification.DiscordSettings,
				req.Qualification.ExamMode,
				req.Qualification.ExamSettings,
			).
			WHERE(
				tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.handleQualificationAccessChanges(ctx, tx, qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, req.Qualification.Id, req.Qualification.Access); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.handleQualificationRequirementsChanges(ctx, tx, req.Qualification.Id, req.Qualification.Requirements); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if req.Qualification.Exam != nil && req.Qualification.Exam.Questions != nil {
		if err := s.handleExamQuestionsChanges(ctx, tx, req.Qualification.Id, req.Qualification.Exam); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateQualificationResponse{
		QualificationId: req.Qualification.Id,
	}, nil
}

func (s *Server) DeleteQualification(ctx context.Context, req *DeleteQualificationRequest) (*DeleteQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.QualificationId,
		tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsqualifications.QualificationsServicePerm, permsqualifications.QualificationsServiceDeleteQualificationPerm, permsqualifications.QualificationsServiceDeleteQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsqualifications.ErrFailedQuery
	}

	stmt := tQuali.
		UPDATE(
			tQuali.DeletedAt,
		).
		SET(
			tQuali.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tQuali.ID.EQ(jet.Uint64(req.QualificationId)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationResponse{}, nil
}
