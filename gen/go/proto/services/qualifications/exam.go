package qualifications

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	tExamQuestions = table.FivenetQualificationsExamQuestions.AS("examquestion")
	tExamResponses = table.FivenetQualificationsExamResponses.AS("examresponse")
	tExamUser      = table.FivenetQualificationsExamUsers.AS("examuser")
)

func (s *Server) GetExamInfo(ctx context.Context, req *GetExamInfoRequest) (*GetExamInfoResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualificationShort(ctx, req.QualificationId, tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	check, err = s.checkIfUserCanTakeExam(ctx, quali, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrExamDisabled
	}

	questionCount, err := s.countExamQuestions(ctx, req.QualificationId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	examUser, err := s.getExamUser(ctx, req.QualificationId, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &GetExamInfoResponse{
		Qualification: quali,
		QuestionCount: int32(questionCount),
		ExamUser:      examUser,
	}, nil
}

func (s *Server) checkIfUserCanTakeExam(ctx context.Context, quali *qualifications.QualificationShort, userInfo *userinfo.UserInfo) (bool, error) {
	if quali.ExamMode <= qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED {
		return false, errorsqualifications.ErrExamDisabled
	} else if quali.ExamMode == qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_REQUEST_NEEDED {
		request, err := s.getQualificationRequest(ctx, quali.Id, userInfo.UserId, userInfo)
		if err != nil {
			return false, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if request == nil || request.Status == nil || (*request.Status != qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED && *request.Status != qualifications.RequestStatus_REQUEST_STATUS_EXAM_STARTED) {
			return false, nil
		}
	}

	return true, nil
}

func (s *Server) getExamUser(ctx context.Context, qualificationId uint64, userId int32) (*qualifications.ExamUser, error) {
	stmt := tExamUser.
		SELECT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.CreatedAt,
			tExamUser.StartedAt,
			tExamUser.EndsAt,
			tExamUser.EndedAt,
		).
		FROM(tExamUser).
		WHERE(jet.AND(
			tExamUser.QualificationID.EQ(jet.Uint64(qualificationId)),
			tExamUser.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	var dest qualifications.ExamUser
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.QualificationId == 0 || dest.UserId == 0 {
		return nil, nil
	}

	return &dest, nil
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

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualificationShort(ctx, req.QualificationId, tQuali.ID.EQ(jet.Uint64(req.QualificationId)), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	check, err = s.checkIfUserCanTakeExam(ctx, quali, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrExamDisabled
	}

	examUser, err := s.getExamUser(ctx, req.QualificationId, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	var exam *qualifications.ExamQuestions
	if examUser == nil || (examUser.EndsAt != nil && time.Since(examUser.EndsAt.AsTime()) < quali.ExamSettings.Time.AsDuration()) {
		exam, err = s.getExamQuestions(ctx, s.db, req.QualificationId, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.updateRequestStatus(ctx, req.QualificationId, userInfo.UserId, qualifications.RequestStatus_REQUEST_STATUS_EXAM_STARTED); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// No end time for the exam? Need to create an entry
	if examUser == nil || examUser.EndsAt == nil {
		examTime := quali.ExamSettings.Time.AsDuration()

		tExamUser := table.FivenetQualificationsExamUsers
		stmt := tExamUser.
			INSERT(
				tExamUser.QualificationID,
				tExamUser.UserID,
				tExamUser.StartedAt,
				tExamUser.EndsAt,
				tExamUser.EndedAt,
			).
			VALUES(
				req.QualificationId,
				userInfo.UserId,
				jet.CURRENT_TIMESTAMP(),
				jet.CURRENT_TIMESTAMP().ADD(jet.INTERVALd(examTime)),
				jet.NULL,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}
	}

	examUser, err = s.getExamUser(ctx, req.QualificationId, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &TakeExamResponse{
		Exam:     exam,
		ExamUser: examUser,
	}, nil
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

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	duration := 0 * time.Second
	endedAt := time.Now()
	examUser, err := s.getExamUser(ctx, req.QualificationId, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if examUser != nil && examUser.StartedAt != nil {
		duration = endedAt.Sub(examUser.StartedAt.AsTime())
	}

	tExamUser := table.FivenetQualificationsExamUsers
	stmt := tExamUser.
		INSERT(
			tExamUser.QualificationID,
			tExamUser.UserID,
			tExamUser.EndedAt,
		).
		VALUES(
			req.QualificationId,
			userInfo.UserId,
			jet.TimestampT(endedAt),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamUser.EndedAt.SET(jet.TimestampT(endedAt)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		return nil, errorsqualifications.ErrFailedQuery
	}

	tExamResponses := table.FivenetQualificationsExamResponses
	respStmt := tExamResponses.
		INSERT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.Responses,
		).
		VALUES(
			req.QualificationId,
			userInfo.UserId,
			req.Responses,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamResponses.Responses.SET(jet.RawString("VALUES(`responses`)")),
		)

	if _, err := respStmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.updateRequestStatus(ctx, req.QualificationId, userInfo.UserId, qualifications.RequestStatus_REQUEST_STATUS_EXAM_GRADING); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &SubmitExamResponse{
		Duration: durationpb.New(duration),
	}, nil
}

func (s *Server) GetUserExam(ctx context.Context, req *GetUserExamRequest) (*GetUserExamResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int("fivenet.user_id", int(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetUserExam",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resp := &GetUserExamResponse{}

	exam, err := s.getExamQuestions(ctx, s.db, req.QualificationId, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.Exam = exam

	responses, err := s.getExamResponses(ctx, req.QualificationId, req.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.Responses = responses

	examUser, err := s.getExamUser(ctx, req.QualificationId, req.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.ExamUser = examUser

	return resp, nil
}

func (s *Server) deleteExamUser(ctx context.Context, tx qrm.DB, qualificationId uint64, userId int32) error {
	stmt := tExamUser.
		DELETE().
		WHERE(jet.AND(
			tExamUser.QualificationID.EQ(jet.Uint64(qualificationId)),
			tExamUser.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
