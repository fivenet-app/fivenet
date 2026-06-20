package qualifications

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) GetExamInfo(
	ctx context.Context,
	req *pbqualifications.GetExamInfoRequest,
) (*pbqualifications.GetExamInfoResponse, error) {
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_TAKE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualificationShort(
		ctx,
		req.GetQualificationId(),
		userInfo,
		false,
	)
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

	questionCount, err := s.store.CountExamQuestions(ctx, req.GetQualificationId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.GetExamInfoResponse{
		Qualification: quali,
		QuestionCount: questionCount,
		ExamUser:      examUser,
	}, nil
}

func (s *Server) checkIfUserCanTakeExam(
	ctx context.Context,
	quali *qualifications.QualificationShort,
	userInfo *userinfo.UserInfo,
) (bool, error) {
	if quali.GetExamMode() <= qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED {
		return false, errorsqualifications.ErrExamDisabled
	} else if quali.GetExamMode() == qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_REQUEST_NEEDED {
		request, err := s.getQualificationRequest(
			ctx,
			quali.GetId(),
			userInfo.GetUserId(),
			userInfo,
		)
		if err != nil {
			return false, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if request == nil || request.Status == nil ||
			(request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED && request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_EXAM_STARTED) {
			return false, nil
		}
	}

	return true, nil
}

func (s *Server) TakeExam(
	ctx context.Context,
	req *pbqualifications.TakeExamRequest,
) (*pbqualifications.TakeExamResponse, error) {
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_TAKE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualificationShort(
		ctx,
		req.GetQualificationId(),
		userInfo,
		false,
	)
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

	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	timesUp := examUser != nil && examUser.GetEndsAt() != nil &&
		time.Since(examUser.GetEndsAt().AsTime()) > 10*time.Second

	var exam *qualificationsexam.ExamQuestions
	if examUser == nil || !timesUp {
		exam, err = s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	var responses *qualificationsexam.ExamResponses
	if examUser != nil && !timesUp {
		responses, _, err = s.store.GetExamResponses(
			ctx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := s.store.UpdateRequestStatus(
		ctx,
		s.db,
		req.GetQualificationId(),
		userInfo.GetUserId(),
		qualifications.RequestStatus_REQUEST_STATUS_EXAM_STARTED,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// No end time for the exam? Need to create an entry
	if examUser == nil || examUser.GetEndsAt() == nil {
		examTime := quali.GetExamSettings().GetTime().AsDuration()
		if err := s.store.CreateExamUser(
			ctx,
			s.db,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			time.Now().Add(examTime),
		); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}
	}

	examUser, err = s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbqualifications.TakeExamResponse{
		Exam:      exam,
		ExamUser:  examUser,
		Responses: responses,

		TimesUp: timesUp,
	}, nil
}

func (s *Server) SubmitExam(
	ctx context.Context,
	req *pbqualifications.SubmitExamRequest,
) (*pbqualifications.SubmitExamResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_TAKE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualification(
		ctx,
		req.GetQualificationId(),
		userInfo,
		false,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	var duration time.Duration
	endedAt := time.Now()
	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if examUser != nil && examUser.GetStartedAt() != nil {
		duration = endedAt.Sub(examUser.GetStartedAt().AsTime())
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.UpsertExamResponses(
		ctx,
		tx,
		req.GetQualificationId(),
		userInfo.GetUserId(),
		req.GetResponses(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Only update the exam user if this is not a partial update, otherwise we might "end" the exam prematurely when the user is still working on it
	if !req.GetPartial() {
		if err := s.store.UpsertExamUserEndedAt(
			ctx,
			tx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			endedAt,
		); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}

			return nil, errorsqualifications.ErrFailedQuery
		}

		if err := s.gradeExam(
			ctx,
			tx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			quali,
			req.GetResponses(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbqualifications.SubmitExamResponse{
		Duration: durationpb.New(duration),
	}, nil
}

func (s *Server) gradeExam(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	quali *qualifications.Qualification,
	responses *qualificationsexam.ExamResponses,
) error {
	if quali.GetExamSettings() != nil && quali.GetExamSettings().GetAutoGrade() {
		exam, err := s.store.GetExamQuestions(ctx, tx, qualificationId, true)
		if err != nil {
			return err
		}
		if exam != nil && len(exam.GetQuestions()) > 0 {
			// Auto grading is enabled, we can grade the exam now
			score, grading := exam.Grade(
				quali.GetExamSettings().GetAutoGradeMode(),
				responses,
			)
			var status qualifications.ResultStatus
			if score >= float32(quali.GetExamSettings().GetMinimumPoints()) {
				status = qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL
			} else {
				status = qualifications.ResultStatus_RESULT_STATUS_FAILED
			}

			if _, err := s.createOrUpdateQualificationResult(
				ctx,
				tx,
				qualificationId,
				0,
				&userinfo.UserInfo{
					Superuser: true,
					Job:       quali.GetCreatorJob(),
					UserId:    0,
				},
				userId,
				status,
				&score,
				"",
				grading,
			); err != nil {
				return err
			}
		}

		if err := s.store.UpdateRequestStatus(
			ctx,
			tx,
			qualificationId,
			userId,
			qualifications.RequestStatus_REQUEST_STATUS_COMPLETED,
		); err != nil {
			return err
		}
	} else {
		if err := s.store.UpdateRequestStatus(
			ctx,
			tx,
			qualificationId,
			userId,
			qualifications.RequestStatus_REQUEST_STATUS_EXAM_GRADING,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) GetUserExam(
	ctx context.Context,
	req *pbqualifications.GetUserExamRequest,
) (*pbqualifications.GetUserExamResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		qualificationIDLogFieldKey, req.GetQualificationId(),
		userIDLogFieldKey, req.GetUserId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resp := &pbqualifications.GetUserExamResponse{}

	exam, err := s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.Exam = exam

	resp.Responses, resp.Grading, err = s.store.GetExamResponses(
		ctx,
		req.GetQualificationId(),
		req.GetUserId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), req.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.ExamUser = examUser

	return resp, nil
}
