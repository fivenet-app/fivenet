package qualifications

import (
	"database/sql"

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
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/mysql"
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

var tQualiFiles = table.FivenetQualificationsFiles

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
	store    qualificationsstore.IStore
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
	Store             qualificationsstore.IStore
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
