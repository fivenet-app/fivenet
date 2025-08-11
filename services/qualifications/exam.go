package qualifications

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	tExamQuestions = table.FivenetQualificationsExamQuestions.AS("exam_question")
	tExamResponses = table.FivenetQualificationsExamResponses.AS("exam_response")
	tExamUser      = table.FivenetQualificationsExamUsers.AS("exam_user")
)

func (s *Server) GetExamInfo(ctx context.Context, req *pbqualifications.GetExamInfoRequest) (*pbqualifications.GetExamInfoResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.QualificationId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
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

	return &pbqualifications.GetExamInfoResponse{
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

func (s *Server) TakeExam(ctx context.Context, req *pbqualifications.TakeExamRequest) (*pbqualifications.TakeExamResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.QualificationId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "TakeExam",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
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

	if err := s.updateRequestStatus(ctx, s.db, req.QualificationId, userInfo.UserId, qualifications.RequestStatus_REQUEST_STATUS_EXAM_STARTED); err != nil {
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

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbqualifications.TakeExamResponse{
		Exam:     exam,
		ExamUser: examUser,
	}, nil
}

func (s *Server) SubmitExam(ctx context.Context, req *pbqualifications.SubmitExamRequest) (*pbqualifications.SubmitExamResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.QualificationId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "SubmitExam",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_TAKE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualification(ctx, req.QualificationId, nil, userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
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

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
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
			tExamResponses.Grading,
		).
		VALUES(
			req.QualificationId,
			userInfo.UserId,
			req.Responses,
			jet.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tExamResponses.Responses.SET(jet.RawString("VALUES(`responses`)")),
		)

	if _, err := respStmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if quali.ExamSettings != nil && quali.ExamSettings.AutoGrade {
		exam, err := s.getExamQuestions(ctx, tx, req.QualificationId, true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if exam != nil && len(exam.Questions) > 0 {
			// Auto grading is enabled, we can grade the exam now
			score, grading := exam.Grade(quali.ExamSettings.AutoGradeMode, req.Responses)
			var status qualifications.ResultStatus
			if score >= float32(quali.ExamSettings.MinimumPoints) {
				status = qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL
			} else {
				status = qualifications.ResultStatus_RESULT_STATUS_FAILED
			}

			if _, err := s.createOrUpdateQualificationResult(ctx, tx, req.QualificationId, 0, &userinfo.UserInfo{
				Superuser: true,
				Job:       quali.CreatorJob,
				UserId:    0,
			}, userInfo.UserId, status, &score, "", grading); err != nil {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}

		if err := s.updateRequestStatus(ctx, tx, req.QualificationId, userInfo.UserId, qualifications.RequestStatus_REQUEST_STATUS_COMPLETED); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	} else {
		if err := s.updateRequestStatus(ctx, tx, req.QualificationId, userInfo.UserId, qualifications.RequestStatus_REQUEST_STATUS_EXAM_GRADING); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbqualifications.SubmitExamResponse{
		Duration: durationpb.New(duration),
	}, nil
}

func (s *Server) GetUserExam(ctx context.Context, req *pbqualifications.GetUserExamRequest) (*pbqualifications.GetUserExamResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.qualifications.id", req.QualificationId,
		"fivenet.user_id", req.UserId,
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "GetUserExam",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resp := &pbqualifications.GetUserExamResponse{}

	exam, err := s.getExamQuestions(ctx, s.db, req.QualificationId, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.Exam = exam

	resp.Responses, resp.Grading, err = s.getExamResponses(ctx, req.QualificationId, req.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	examUser, err := s.getExamUser(ctx, req.QualificationId, req.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.ExamUser = examUser

	return resp, nil
}

func (s *Server) deleteExamUser(ctx context.Context, tx qrm.DB, qualificationId uint64, userId int32) error {
	tExamUser := table.FivenetQualificationsExamUsers

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
