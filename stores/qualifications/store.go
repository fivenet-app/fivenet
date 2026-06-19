package qualificationsstore

import (
	"context"
	"database/sql"
	"time"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type ListQualificationsOptions struct {
	Pagination *database.PaginationRequest
	Sort       *database.Sort
	Search     string
}

type ListQualificationRequestsOptions struct {
	Pagination      *database.PaginationRequest
	Sort            *database.Sort
	QualificationID int64
	Status          []resqualifications.RequestStatus
	UserIDs         []int32
}

type ListQualificationsResultsOptions struct {
	Pagination      *database.PaginationRequest
	Sort            *database.Sort
	QualificationID int64
	Status          []resqualifications.ResultStatus
	UserIDs         []int32
}

type IStore interface {
	ListQualifications(
		ctx context.Context,
		opts ListQualificationsOptions,
		userInfo *userinfo.UserInfo,
		includePhoneNumber bool,
	) (*pbqualifications.ListQualificationsResponse, error)
	GetQualification(
		ctx context.Context,
		qualificationId int64,
		userInfo *userinfo.UserInfo,
		selectContent bool,
		includePhoneNumber bool,
	) (*resqualifications.Qualification, error)
	GetQualificationShort(
		ctx context.Context,
		qualificationId int64,
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
		opts ListQualificationRequestsOptions,
		userInfo *userinfo.UserInfo,
		includePhoneNumber bool,
	) (*pbqualifications.ListQualificationRequestsResponse, error)
	GetQualificationRequest(
		ctx context.Context,
		qualificationId int64,
		userId int32,
		userInfo *userinfo.UserInfo,
		includePhoneNumber bool,
	) (*resqualifications.QualificationRequest, error)
	ListQualificationsResults(
		ctx context.Context,
		opts ListQualificationsResultsOptions,
		userInfo *userinfo.UserInfo,
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
	db                  *sql.DB
	access              *access.SubjectObjectAccess
	resultSorter        *database.SorterBuilder
	requestSorter       *database.SorterBuilder
	qualificationSorter *database.SorterBuilder
}

func New(db *sql.DB) IStore {
	return &Store{
		db:     db,
		access: access.NewQualificationsSubjectObjectAccess(db),
		resultSorter: database.New(
			database.SpecMap{
				"status":    database.Column{Col: tQualiResult.Status},
				"createdAt": database.Column{Col: tQualiResult.CreatedAt},
			},
			[]mysql.OrderByClause{tQualiResult.CreatedAt.DESC()},
			nil,
			"createdAt",
			3,
		),
		requestSorter: database.New(
			database.SpecMap{
				"status":     database.Column{Col: tQualiReq.Status},
				"approvedAt": database.Column{Col: tQualiReq.ApprovedAt},
				"createdAt":  database.Column{Col: tQualiReq.CreatedAt},
			},
			[]mysql.OrderByClause{tQualiReq.CreatedAt.DESC()},
			nil,
			"createdAt",
			3,
		),
		qualificationSorter: database.New(
			database.SpecMap{
				"abbreviation": database.Column{Col: tQuali.Abbreviation},
				"id":           database.Column{Col: tQualiResult.ID},
			},
			[]mysql.OrderByClause{tQualiResult.ID.DESC()},
			nil,
			"id",
			3,
		),
	}
}
