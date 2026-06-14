package qualifications

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
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetQualifications,
		IDColumn:        table.FivenetQualifications.ID,
		JobColumn:       table.FivenetQualifications.Job,
		DeletedAtColumn: table.FivenetQualifications.DeletedAt,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:           table.FivenetQualificationsExamUsers,
				DeletedAtColumn: table.FivenetQualificationsExamUsers.EndsAt,
				ForeignKey:      table.FivenetQualificationsExamUsers.QualificationID,

				MinDays: 60,
			},
			{
				Table:           table.FivenetQualificationsRequests,
				DeletedAtColumn: table.FivenetQualificationsRequests.DeletedAt,
				ForeignKey:      table.FivenetQualificationsRequests.QualificationID,

				MinDays: 60,
			},
			{
				Table:           table.FivenetQualificationsResults,
				IDColumn:        table.FivenetQualificationsResults.ID,
				DeletedAtColumn: table.FivenetQualificationsResults.DeletedAt,
				ForeignKey:      table.FivenetQualificationsResults.QualificationID,

				MinDays: 60,
			},
		},
	})
}

var (
	tQuali      = table.FivenetQualifications.AS("qualification")
	tQualiFiles = table.FivenetQualificationsFiles
)

type Server struct {
	pbqualifications.QualificationsServiceServer
	pbqualifications.ExamServiceServer

	logger   *zap.Logger
	db       *sql.DB
	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	notif    notifi.INotifi
	storage  storage.IStorage

	access         *access.SubjectObjectAccess
	accessResolver *access.SubjectResolver

	fHandler *filestore.Handler[int64]
	store    qualificationsStore
}

type qualificationsStore interface {
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

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger            *zap.Logger
	DB                *sql.DB
	Perms             perms.Permissions
	UserAwareEnricher mstlystcdata.IUserAwareEnricher
	Config            *config.Config
	Notif             notifi.INotifi
	Storage           storage.IStorage
	Store             qualificationsStore `optional:"true"`
}

func NewServer(p Params) *Server {
	// 3 MiB limit
	qualiFileHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tQualiFiles,
		tQualiFiles.QualificationID,
		tQualiFiles.FileID,
		3<<20,
		5,
		func(parentId int64) mysql.BoolExpression {
			return tQualiFiles.QualificationID.EQ(mysql.Int64(parentId))
		},
		filestore.InsertJoinRow,
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	s := &Server{
		logger: p.Logger.Named("jobs"),

		db:       p.DB,
		perms:    p.Perms,
		enricher: p.UserAwareEnricher,
		notif:    p.Notif,
		storage:  p.Storage,

		access:         access.NewQualificationsSubjectObjectAccess(p.DB),
		accessResolver: access.NewSubjectResolver(p.DB),

		fHandler: qualiFileHandler,
		store:    p.Store,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbqualifications.RegisterQualificationsServiceServer(srv, s)
	pbqualifications.RegisterExamServiceServer(srv, s)
}
