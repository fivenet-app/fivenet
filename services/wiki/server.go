package wiki

import (
	"context"
	"database/sql"

	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
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
	"github.com/go-jet/jet/v2/qrm"
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
	store        wikiStore
}

type wikiStore interface {
	ListPages(ctx context.Context, q wikistore.ListPagesQuery) (*wikistore.ListPagesResult, error)
	GetPage(ctx context.Context, pageID int64, withContent bool) (*reswiki.Page, error)
	GetPageOrderInfo(ctx context.Context, q qrm.DB, pageID int64) (*wikistore.PageOrderInfo, error)
	NextPageGroupRank(
		ctx context.Context,
		q qrm.DB,
		job string,
		parentID *int64,
		startpage bool,
		excludeID int64,
	) (string, error)
	InsertPageGroupRank(
		ctx context.Context,
		q qrm.DB,
		job string,
		parentID *int64,
		startpage bool,
		excludeID int64,
		beforeID, afterID *int64,
	) (string, error)
	CountPageActivity(ctx context.Context, pageID int64) (int64, error)
	ListPageActivity(
		ctx context.Context,
		pageID int64,
		offset int64,
		limit int64,
	) ([]*wikiactivity.PageActivity, error)
	AddPageActivity(
		ctx context.Context,
		tx qrm.DB,
		activity *wikiactivity.PageActivity,
	) (int64, error)
	CountPageChildren(ctx context.Context, pageID int64) (int64, error)
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
	Store    wikiStore `optional:"true"`
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

	store := p.Store
	if store == nil {
		store = wikistore.New(p.DB)
	}

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
		store:          store,
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
