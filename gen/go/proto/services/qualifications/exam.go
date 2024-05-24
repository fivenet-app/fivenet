package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tExam          = table.FivenetQualificationsExam.AS("exam")
	tExamResponses = table.FivenetQualificationsExamResponses.AS("examuserresponse")
)

func (s *Server) GetExam(ctx context.Context, req *GetExamRequest) (*GetExamResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetExam",
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

	questions, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if userInfo.SuperUser {
		check = true
	}

	exam, err := s.getExam(ctx, req.QualificationId, userInfo, questions, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &GetExamResponse{
		Exam: exam,
	}, nil
}

func (s *Server) CreateOrUpdateExam(ctx context.Context, req *CreateOrUpdateExamRequest) (*CreateOrUpdateExamResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateExam",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.Exam.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if req.Exam.Id > 0 {
		stmt := tExam.
			UPDATE(
				tExam.Settings,
				tExam.Questions,
			).
			SET(
				req.Exam.Settings,
				req.Exam.Questions,
			).
			WHERE(tExam.QualificationID.EQ(jet.Uint64(req.Exam.QualificationId)))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {
		stmt := tExam.
			INSERT(
				tExam.QualificationID,
				tExam.Settings,
				tExam.Questions,
			).
			VALUES(
				req.Exam.QualificationId,
				req.Exam.Settings,
				req.Exam.Questions,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	}

	exam, err := s.getExam(ctx, req.Exam.QualificationId, userInfo, true, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &CreateOrUpdateExamResponse{
		Exam: exam,
	}, nil
}

func (s *Server) TakeExam(ctx context.Context, req *TakeExamRequest) (*TakeExamResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "TakeExam",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	// TODO

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return nil, nil
}

func (s *Server) SubmitExam(ctx context.Context, req *SubmitExamRequest) (*SubmitExamResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "SubmitExam",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	// TODO

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return nil, nil
}
