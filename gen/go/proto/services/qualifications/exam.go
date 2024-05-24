package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	// TODO

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return nil, nil
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

	if req.Exam.Id > 0 {

		// TODO

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {

		// TODO

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	}
	// TODO

	return nil, nil
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

	// TODO

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return nil, nil
}
