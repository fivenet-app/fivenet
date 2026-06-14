package qualificationsstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	ListQualifications(
		ctx context.Context,
		req *pbqualifications.ListQualificationsRequest,
		userInfo *userinfo.UserInfo,
		where mysql.BoolExpression,
		includePhoneNumber bool,
	) (*pbqualifications.ListQualificationsResponse, error)
	GetQualification(
		ctx context.Context,
		qualificationId int64,
		where mysql.BoolExpression,
		userInfo *userinfo.UserInfo,
		selectContent bool,
		includePhoneNumber bool,
	) (*resqualifications.Qualification, error)
	GetQualificationShort(
		ctx context.Context,
		qualificationId int64,
		where mysql.BoolExpression,
		userInfo *userinfo.UserInfo,
		includePhoneNumber bool,
	) (*resqualifications.QualificationShort, error)
	GetQualificationRequirements(
		ctx context.Context,
		qualificationId int64,
	) ([]*resqualifications.QualificationRequirement, error)
	CheckRequirementsMetForQualification(
		ctx context.Context,
		qualificationId int64,
		userId int32,
	) (bool, error)
	ListQualificationRequests(
		ctx context.Context,
		req *pbqualifications.ListQualificationRequestsRequest,
		userInfo *userinfo.UserInfo,
		where mysql.BoolExpression,
		includePhoneNumber bool,
	) (*pbqualifications.ListQualificationRequestsResponse, error)
	GetQualificationRequest(
		ctx context.Context,
		qualificationId int64,
		userId int32,
		includePhoneNumber bool,
	) (*resqualifications.QualificationRequest, error)
	ListQualificationsResults(
		ctx context.Context,
		req *pbqualifications.ListQualificationsResultsRequest,
		userInfo *userinfo.UserInfo,
		where mysql.BoolExpression,
		includePhoneNumber bool,
	) (*pbqualifications.ListQualificationsResultsResponse, error)
	GetQualificationResult(
		ctx context.Context,
		qualificationId int64,
		resultId int64,
		status []resqualifications.ResultStatus,
		userInfo *userinfo.UserInfo,
		userId int32,
		includePhoneNumber bool,
	) (*resqualifications.QualificationResult, error)
	GetExamUser(
		ctx context.Context,
		qualificationId int64,
		userId int32,
	) (*qualificationsexam.ExamUser, error)
	GetExamQuestions(
		ctx context.Context,
		q qrm.DB,
		qualificationId int64,
		withAnswers bool,
	) (*qualificationsexam.ExamQuestions, error)
	CountExamQuestions(ctx context.Context, qualificationId int64) (int64, error)
	GetExamResponses(
		ctx context.Context,
		qualificationId int64,
		userId int32,
	) (*qualificationsexam.ExamResponses, *qualificationsexam.ExamGrading, error)
	CreateQualification(ctx context.Context, tx qrm.DB, userInfo *userinfo.UserInfo) (int64, error)
	UpdateQualification(
		ctx context.Context,
		tx qrm.DB,
		quali *resqualifications.Qualification,
	) error
	DeleteQualification(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		deletedAt *timestamp.Timestamp,
	) error
	HandleQualificationRequirementsChanges(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		reqs []*resqualifications.QualificationRequirement,
	) error
	CreateQualificationResult(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		status resqualifications.ResultStatus,
		score *float32,
		summary string,
		creator *userinfo.UserInfo,
	) (int64, error)
	UpdateQualificationResult(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		resultId int64,
		userId int32,
		status resqualifications.ResultStatus,
		score *float32,
		summary string,
	) error
	UpdateExamResponseGrading(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		grading *qualificationsexam.ExamGrading,
	) error
	ApproveQualificationRequest(
		ctx context.Context,
		tx qrm.DB,
		req *resqualifications.QualificationRequest,
		userInfo *userinfo.UserInfo,
	) error
	UpsertQualificationRequest(
		ctx context.Context,
		tx qrm.DB,
		req *resqualifications.QualificationRequest,
	) error
	CreateExamUser(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		endsAt time.Time,
	) error
	UpsertExamResponses(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		responses *qualificationsexam.ExamResponses,
	) error
	UpsertExamUserEndedAt(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		endedAt time.Time,
	) error
	HandleExamQuestionsChanges(
		ctx context.Context,
		tx *sql.Tx,
		qualificationId int64,
		questions *qualificationsexam.ExamQuestions,
	) ([]*file.File, error)
	DeleteQualificationRequest(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
	) error
	UpdateRequestStatus(
		ctx context.Context,
		tx qrm.DB,
		qualificationId int64,
		userId int32,
		status resqualifications.RequestStatus,
	) error
	DeleteQualificationResult(ctx context.Context, tx qrm.DB, resultId int64) error
	DeleteExamUser(ctx context.Context, tx qrm.DB, qualificationId int64, userId int32) error
}

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) IStore {
	return &Store{db: db}
}
