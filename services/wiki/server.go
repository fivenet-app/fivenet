package wiki

import (
	"context"
	"database/sql"

	pbwiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/collab"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	wikistore "github.com/fivenet-app/fivenet/v2026/stores/wiki"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetWikiPages,
		IDColumn:        table.FivenetWikiPages.ID,
		DeletedAtColumn: table.FivenetWikiPages.DeletedAt,
		JobColumn:       table.FivenetWikiPages.Job,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:      table.FivenetWikiPagesActivity,
				ForeignKey: table.FivenetWikiPagesActivity.PageID,
				IDColumn:   table.FivenetWikiPagesActivity.ID,
			},
		},
	},
	)
}

type Server struct {
	pbwiki.WikiServiceServer
	pbwiki.CollabServiceServer

	logger *zap.Logger
	db     *sql.DB
	js     *events.JSWrapper

	perms    perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	notifi   notifi.INotifi

	access         *access.SubjectObjectAccess
	accessResolver *access.SubjectResolver

	collabServer *collab.CollabServer
	fHandler     *filestore.Handler[int64]
	store        wikistore.IStore
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	DB       *sql.DB
	Perms    perms.Permissions
	Enricher mstlystcdata.IUserAwareEnricher
	JS       *events.JSWrapper
	Storage  storage.IStorage
	Notifi   notifi.INotifi
	Store    wikistore.IStore `optional:"true"`
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	collabServer := collab.New(ctxCancel, p.Logger, p.JS, "wiki_pages")

	tPageFiles := table.FivenetWikiPagesFiles
	fHandler := filestore.NewHandler(
		p.Storage,
		p.DB,
		tPageFiles,
		tPageFiles.PageID,
		tPageFiles.FileID,
		3<<20,
		5,
		func(parentID int64) mysql.BoolExpression {
			return tPageFiles.PageID.EQ(mysql.Int64(parentID))
		},
		filestore.InsertJoinRow,
		false,
	).WithUploadFilter(filestore.NewImageUploadFilter())

	objAccess := access.NewWikiPageSubjectObjectAccess(p.DB)
	access.RegisterAccess("wiki_page", objAccess)

	s := &Server{
		logger: p.Logger.Named("wiki"),
		db:     p.DB,
		js:     p.JS,

		perms:    p.Perms,
		enricher: p.Enricher,
		notifi:   p.Notifi,

		access:         objAccess,
		accessResolver: access.NewSubjectResolver(p.DB),
		collabServer:   collabServer,
		fHandler:       fHandler,
		store:          p.Store,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		return s.collabServer.Start(ctxStartup)
	}))

	p.LC.Append(fx.StopHook(func(ctxStartup context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbwiki.RegisterWikiServiceServer(srv, s)
	pbwiki.RegisterCollabServiceServer(srv, s)
}
