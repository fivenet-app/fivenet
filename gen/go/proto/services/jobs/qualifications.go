package jobs

import (
	"context"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const QualificationsPageSize = 5

var (
	tQuali = table.FivenetJobsQualifications.AS("qualification")
)

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
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*jobs.Qualification{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.listQualificationsQuery(condition, nil, userInfo).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tQuali.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
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
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "GetQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsjobs.ErrFailedQuery
	}

	resp := &GetQualificationResponse{}
	resp.Qualification, err = s.getQualification(ctx, req.QualificationId,
		tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if resp.Qualification == nil || resp.Qualification.Id <= 0 {
		return nil, errorsjobs.ErrFailedQuery
	}

	if resp.Qualification.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.Qualification.Creator)
	}

	qualiAccess, err := s.GetQualificationAccess(ctx, &GetQualificationAccessRequest{
		QualificationId: req.QualificationId,
	})
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if qualiAccess.Access != nil {
		resp.Qualification.Access = qualiAccess.Access
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateQualification(ctx context.Context, req *CreateQualificationRequest) (*CreateQualificationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tQuali := table.FivenetJobsQualifications
	stmt := tQuali.
		INSERT(
			tQuali.Job,
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.Content,
			tQuali.CreatorID,
			tQuali.CreatorJob,
			tQuali.DiscordSettings,
		).
		VALUES(
			userInfo.Job,
			req.Qualification.Weight,
			req.Qualification.Closed,
			req.Qualification.Abbreviation,
			req.Qualification.Title,
			req.Qualification.Description,
			req.Qualification.Content,
			userInfo.UserId,
			userInfo.Job,
			req.Qualification.DiscordSettings,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if err := s.handleQualificationAccessChanges(ctx, tx, jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, uint64(lastId), req.Qualification.Access); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateQualificationResponse{
		QualificationId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateQualification(ctx context.Context, req *UpdateQualificationRequest) (*UpdateQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.qualifications.id", int64(req.Qualification.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "UpdateQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.Qualification.Id, userInfo, jobs.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var onlyUpdateAccess bool
	if !check && !userInfo.SuperUser {
		onlyUpdateAccess, err = s.checkIfUserHasAccessToQuali(ctx, req.Qualification.Id, userInfo, jobs.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errorsjobs.ErrFailedQuery
		}
		if !onlyUpdateAccess {
			return nil, errorsjobs.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.Qualification.Id,
		tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
		userInfo)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsQualificationsServicePerm, permsjobs.JobsQualificationsServiceUpdateQualificationPerm, permsjobs.JobsQualificationsServiceUpdateQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if !onlyUpdateAccess {
		if req.Qualification.Description != nil {
			*req.Qualification.Description = strings.TrimSuffix(*req.Qualification.Description, "<br>")
		}

		tQuali := table.FivenetJobsQualifications
		stmt := tQuali.
			UPDATE(
				tQuali.Weight,
				tQuali.Closed,
				tQuali.Abbreviation,
				tQuali.Title,
				tQuali.Description,
				tQuali.Content,
				tQuali.DiscordSettings,
			).
			SET(
				req.Qualification.Weight,
				req.Qualification.Closed,
				req.Qualification.Abbreviation,
				req.Qualification.Title,
				req.Qualification.Description,
				req.Qualification.Content,
				req.Qualification.DiscordSettings,
			).
			WHERE(
				tQuali.ID.EQ(jet.Uint64(req.Qualification.Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	if err := s.handleQualificationAccessChanges(ctx, tx, jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, req.Qualification.Id, req.Qualification.Access); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateQualificationResponse{
		QualificationId: req.Qualification.Id,
	}, nil
}

func (s *Server) DeleteQualification(ctx context.Context, req *DeleteQualificationRequest) (*DeleteQualificationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualification",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsjobs.ErrFailedQuery
		}
	}

	quali, err := s.getQualification(ctx, req.QualificationId,
		tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsQualificationsServicePerm, permsjobs.JobsQualificationsServiceDeleteQualificationPerm, permsjobs.JobsQualificationsServiceDeleteQualificationAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, quali.CreatorJob, quali.Creator) {
		return nil, errorsjobs.ErrFailedQuery
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
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationResponse{}, nil
}
