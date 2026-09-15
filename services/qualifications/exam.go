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
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

const examSubmissionGracePeriod = qualificationsstore.ExamSubmissionGracePeriod

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
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualificationShort(
		ctx,
		req.GetQualificationId(),
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// A completed or pending-grading attempt remains viewable even though its
	// request no longer passes the eligibility check for starting a new exam.
	if examUser == nil {
		check, err = s.checkIfUserCanTakeExam(ctx, quali, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !check {
			return nil, errorsqualifications.ErrExamDisabled
		}
	}

	questionCount, err := s.store.CountExamQuestions(ctx, req.GetQualificationId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.GetExamInfoResponse{
		Qualification: quali,
		QuestionCount: questionCount,
		ExamUser:      publicExamUser(examUser),
	}, nil
}

func (s *Server) checkIfUserCanTakeExam(
	ctx context.Context,
	quali *qualifications.QualificationShort,
	userInfo *userinfo.UserInfo,
) (bool, error) {
	if quali.GetClosed() || quali.GetDraft() {
		return false, nil
	}
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

	requirementsMet, err := s.store.CheckRequirementsMetForQualification(
		ctx,
		quali.GetId(),
		userInfo.GetUserId(),
	)
	if err != nil {
		return false, err
	}
	if !requirementsMet {
		return false, nil
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
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualificationShort(
		ctx,
		req.GetQualificationId(),
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if req.GetCancel() {
		if examUser == nil || examUser.GetEndedAt() != nil ||
			examUser.GetEndsAt() == nil || !time.Now().Before(examUser.GetEndsAt().AsTime()) {
			return nil, errorsqualifications.ErrExamDisabled
		}
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		defer tx.Rollback()
		active, err := s.store.ClaimActiveExamUser(
			ctx,
			tx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			examUser.GetAttemptId(),
			false,
			0,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !active {
			return nil, errorsqualifications.ErrExamDisabled
		}
		if err := s.store.DeleteExamResponses(
			ctx,
			tx,
			examUser.GetAttemptId(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if err := s.store.DeleteExamUser(
			ctx,
			tx,
			examUser.GetAttemptId(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if err := s.store.UpdateRequestStatus(
			ctx,
			tx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if err := tx.Commit(); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		return &pbqualifications.TakeExamResponse{}, nil
	}

	check, err = s.checkIfUserCanTakeExam(ctx, quali, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrExamDisabled
	}

	timesUp := examUser != nil && examUser.GetEndsAt() != nil &&
		!time.Now().Before(examUser.GetEndsAt().AsTime())
	if examUser != nil && (examUser.GetEndedAt() != nil || timesUp) {
		return &pbqualifications.TakeExamResponse{
			ExamUser: publicExamUser(examUser),
			TimesUp:  true,
		}, nil
	}

	var exam *qualificationsexam.ExamQuestions
	if examUser != nil && examUser.GetSnapshot() != nil {
		exam = examForCandidate(examUser.GetSnapshot().GetExam())
	} else if examUser == nil || !timesUp {
		exam, err = s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	var responses *qualificationsexam.ExamResponses
	if examUser != nil && !timesUp {
		responses, _, err = s.store.GetExamResponses(
			ctx,
			s.db,
			examUser.GetAttemptId(),
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
		examWithAnswers, err := s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if _, err := s.store.CreateExamUser(
			ctx,
			s.db,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			time.Now().Add(examTime),
			&qualificationsexam.ExamSnapshot{
				Exam:     examWithAnswers,
				Settings: quali.GetExamSettings(),
			},
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
		ExamUser:  publicExamUser(examUser),
		Responses: publicExamResponses(exam, responses),

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
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualification(
		ctx,
		req.GetQualificationId(),
		userInfo,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if quali.GetExamMode() <= qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED {
		return nil, errorsqualifications.ErrExamDisabled
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
	if examUser == nil || examUser.GetEndedAt() != nil || examUser.GetEndsAt() == nil ||
		!endedAt.Before(examUser.GetEndsAt().AsTime().Add(examSubmissionGracePeriod)) {
		return nil, status.Error(codes.FailedPrecondition, "exam attempt is not active")
	}

	exam := examUser.GetSnapshot().GetExam()
	if exam == nil {
		exam, err = s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}
	responses, err := normalizeExamResponses(exam, req.GetResponses(), req.GetPartial())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid exam submission: %v", err)
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()
	publishNotifications := make([]func(context.Context) error, 0, 1)
	active, err := s.store.ClaimActiveExamUser(
		ctx,
		tx,
		req.GetQualificationId(),
		userInfo.GetUserId(),
		examUser.GetAttemptId(),
		!req.GetPartial(),
		examSubmissionGracePeriod,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !active {
		return nil, status.Error(codes.FailedPrecondition, "exam attempt is not active")
	}
	if req.GetPartial() {
		existing, _, err := s.store.GetExamResponses(
			ctx,
			tx,
			examUser.GetAttemptId(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		responses = mergeExamResponses(exam, existing, responses)
	}

	if err := s.store.UpsertExamResponses(
		ctx,
		tx,
		req.GetQualificationId(),
		userInfo.GetUserId(),
		examUser.GetAttemptId(),
		responses,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Only update the exam user if this is not a partial update, otherwise we might "end" the exam prematurely when the user is still working on it
	if !req.GetPartial() {
		if err := s.gradeExam(
			ctx,
			tx,
			req.GetQualificationId(),
			userInfo.GetUserId(),
			quali,
			examUser.GetSnapshot(),
			responses,
			&publishNotifications,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	notifi.PublishAfterCommit(
		ctx,
		s.logger,
		"qualification_result_updated",
		publishNotifications...)

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
	snapshot *qualificationsexam.ExamSnapshot,
	responses *qualificationsexam.ExamResponses,
	publishNotifications *[]func(context.Context) error,
) error {
	settings := quali.GetExamSettings()
	if snapshot.GetSettings() != nil {
		settings = snapshot.GetSettings()
	}
	exam := snapshot.GetExam()
	if exam == nil {
		var err error
		exam, err = s.store.GetExamQuestions(ctx, tx, qualificationId, true)
		if err != nil {
			return err
		}
	}
	if settings != nil && settings.GetAutoGrade() &&
		validateExamAutoGrading(exam, settings) == nil {
		if exam != nil && len(exam.GetQuestions()) > 0 {
			// Auto grading is enabled, we can grade the exam now
			score, grading := exam.Grade(
				settings.GetAutoGradeMode(),
				responses,
			)
			var status qualifications.ResultStatus
			if score >= float32(settings.GetMinimumPoints()) {
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
				false,
				publishNotifications,
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
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resp := &pbqualifications.GetUserExamResponse{}
	examUser, err := s.store.GetExamUser(ctx, req.GetQualificationId(), req.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if examUser == nil {
		return &pbqualifications.GetUserExamResponse{}, nil
	}

	resp.Responses, resp.Grading, err = s.store.GetExamResponses(
		ctx,
		s.db,
		examUser.GetAttemptId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	resp.ExamUser = publicExamUser(examUser)
	if examUser.GetSnapshot().GetExam() != nil {
		resp.Exam = examUser.GetSnapshot().GetExam()
	} else {
		resp.Exam, err = s.store.GetExamQuestions(ctx, s.db, req.GetQualificationId(), true)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}
	resp.Responses = publicExamResponses(resp.Exam, resp.Responses)

	return resp, nil
}
