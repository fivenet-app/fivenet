package wiki

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
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
		TimestampColumn: table.FivenetWikiPages.DeletedAt,
		MinDays:         60,
	})

	housekeeper.AddJobTable(
		&housekeeper.JobTable{
			TargetTable:           table.FivenetWikiPages,
			TargetSourceIDColumn:  table.FivenetWikiPages.ID,
			TargetDeletedAtColumn: table.FivenetWikiPages.DeletedAt,
			TargetJobColumn:       table.FivenetWikiPages.Job,
		},
		&housekeeper.JobTable{
			Source: &housekeeper.JobTableSource{
				SourceTable:           table.FivenetWikiPages,
				SourceJobColumn:       table.FivenetWikiPages.Job,
				SourceDeletedAtColumn: table.FivenetWikiPages.DeletedAt,
				SourceIDColumn:        table.FivenetWikiPages.ID,
			},
			TargetTable:          table.FivenetWikiPageActivity,
			TargetSourceIDColumn: table.FivenetWikiPageActivity.PageID,
		},
	)
}

var (
	tPage      = table.FivenetWikiPages.AS("page")
	tPageShort = table.FivenetWikiPages.AS("pageshort")
	tPAccess   = table.FivenetWikiPagesAccess.AS("access")
	tJobProps  = table.FivenetJobProps
)

type Server struct {
	pbwiki.WikiServiceServer

	logger   *zap.Logger
	db       *sql.DB
	aud      audit.IAuditer
	perms    perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	htmlDiff *htmldiffer.Differ

	access *access.Grouped[wiki.PageJobAccess, *wiki.PageJobAccess, wiki.PageUserAccess, *wiki.PageUserAccess, access.DummyQualificationAccess[wiki.AccessLevel], *access.DummyQualificationAccess[wiki.AccessLevel], wiki.AccessLevel]
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
}

func NewServer(p Params) *Server {
	s := &Server{
		logger:   p.Logger.Named("wiki"),
		db:       p.DB,
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
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbwiki.RegisterWikiServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbwiki.PermsRemap
}
