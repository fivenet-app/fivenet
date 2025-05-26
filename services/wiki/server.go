package wiki

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/html/htmldiffer"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const defaultWikiUpperLimit = 250

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

var (
	tPage      = table.FivenetWikiPages.AS("page")
	tPageShort = table.FivenetWikiPages.AS("page_short")
	tPAccess   = table.FivenetWikiPagesAccess.AS("access")
	tJobProps  = table.FivenetJobProps
)

type Server struct {
	pbwiki.WikiServiceServer
	pbwiki.CollabServiceServer

	logger *zap.Logger
	db     *sql.DB
	js     *events.JSWrapper

	aud      audit.IAuditer
	perms    perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	htmlDiff *htmldiffer.Differ

	access *access.Grouped[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.PageUserAccess, *wiki.PageUserAccess, access.DummyQualificationAccess[wiki.AccessLevel], *access.DummyQualificationAccess[wiki.AccessLevel], wiki.AccessLevel]

	collabServer *collab.CollabServer
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger     *zap.Logger
	DB         *sql.DB
	Audit      audit.IAuditer
	Perms      perms.Permissions
	Enricher   *mstlystcdata.UserAwareEnricher
	HTMLDiffer *htmldiffer.Differ
	JS         *events.JSWrapper
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	collabServer := collab.New(ctxCancel, p.Logger, p.JS, "wiki_pages")

	s := &Server{
		logger: p.Logger.Named("wiki"),
		db:     p.DB,
		js:     p.JS,

		aud:      p.Audit,
		perms:    p.Perms,
		enricher: p.Enricher,
		htmlDiff: p.HTMLDiffer,

		access: access.NewGrouped[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.PageUserAccess, *wiki.PageUserAccess, access.DummyQualificationAccess[wiki.AccessLevel]](
			p.DB,
			table.FivenetWikiPages,
			&access.TargetTableColumns{
				ID:         table.FivenetWikiPages.ID,
				DeletedAt:  table.FivenetWikiPages.DeletedAt,
				CreatorID:  table.FivenetWikiPages.CreatorID,
				CreatorJob: table.FivenetWikiPages.Job,
			},
			access.NewJobs[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.AccessLevel](
				table.FivenetWikiPagesAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetWikiPagesAccess.ID,
						TargetID: table.FivenetWikiPagesAccess.TargetID,
						Access:   table.FivenetWikiPagesAccess.Access,
					},
					Job:          table.FivenetWikiPagesAccess.Job,
					MinimumGrade: table.FivenetWikiPagesAccess.MinimumGrade,
				},
				table.FivenetWikiPagesAccess.AS("page_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetWikiPagesAccess.AS("page_job_access").ID,
						TargetID: table.FivenetWikiPagesAccess.AS("page_job_access").TargetID,
						Access:   table.FivenetWikiPagesAccess.AS("page_job_access").Access,
					},
					Job:          table.FivenetWikiPagesAccess.AS("page_job_access").Job,
					MinimumGrade: table.FivenetWikiPagesAccess.AS("page_job_access").MinimumGrade,
				},
			),
			access.NewUsers[wiki.PageUserAccess, *wiki.PageUserAccess, wiki.AccessLevel](
				table.FivenetWikiPagesAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetWikiPagesAccess.ID,
						TargetID: table.FivenetWikiPagesAccess.TargetID,
						Access:   table.FivenetWikiPagesAccess.Access,
					},
					UserId: table.FivenetWikiPagesAccess.UserID,
				},
				table.FivenetWikiPagesAccess.AS("page_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetWikiPagesAccess.AS("page_user_access").ID,
						TargetID: table.FivenetWikiPagesAccess.AS("page_user_access").TargetID,
						Access:   table.FivenetWikiPagesAccess.AS("page_user_access").Access,
					},
					UserId: table.FivenetWikiPagesAccess.AS("page_user_access").UserID,
				},
			),
			nil,
		),
		collabServer: collabServer,
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

func (s *Server) GetPermsRemap() map[string]string {
	return pbwiki.PermsRemap
}
