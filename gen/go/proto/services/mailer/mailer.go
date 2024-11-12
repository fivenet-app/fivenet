package mailer

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var (
	tUsers   = table.Users.AS("usershort")
	tCreator = table.Users.AS("creator")

	tUserProps = table.FivenetUserProps

	tThreads           = table.FivenetMsgsThreads.AS("thread")
	tThreadsUserState  = table.FivenetMsgsThreadsUserState.AS("threaduserstate")
	tThreadsUserAccess = table.FivenetMsgsThreadsUserAccess

	tMessages = table.FivenetMsgsMessages.AS("message")
)

type Server struct {
	MailerServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	js       *events.JSWrapper

	access *access.Grouped[mailer.ThreadJobAccess, *mailer.ThreadJobAccess, mailer.ThreadUserAccess, *mailer.ThreadUserAccess, mailer.AccessLevel]
}

type Params struct {
	fx.In

	DB       *sql.DB
	P        perms.Permissions
	Enricher *mstlystcdata.UserAwareEnricher
	Aud      audit.IAuditer
	JS       *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		p:        p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		js:       p.JS,

		access: access.NewGrouped[mailer.ThreadJobAccess, *mailer.ThreadJobAccess, mailer.ThreadUserAccess, *mailer.ThreadUserAccess, mailer.AccessLevel](
			p.DB,
			table.FivenetMsgsThreads,
			&access.TargetTableColumns{
				ID:         table.FivenetMsgsThreads.ID,
				DeletedAt:  table.FivenetMsgsThreads.DeletedAt,
				CreatorID:  table.FivenetMsgsThreads.CreatorID,
				CreatorJob: table.FivenetMsgsThreads.CreatorJob,
			},
			nil,
			access.NewUsers[mailer.ThreadUserAccess, *mailer.ThreadUserAccess, mailer.AccessLevel](
				table.FivenetMsgsThreadsUserAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMsgsThreadsUserAccess.ID,
						CreatedAt: table.FivenetMsgsThreadsUserAccess.CreatedAt,
						TargetID:  table.FivenetMsgsThreadsUserAccess.ThreadID,
						Access:    table.FivenetMsgsThreadsUserAccess.Access,
					},
					UserId: table.FivenetMsgsThreadsUserAccess.UserID,
				},
				table.FivenetMsgsThreadsUserAccess.AS("thread_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMsgsThreadsUserAccess.AS("thread_user_access").ID,
						CreatedAt: table.FivenetMsgsThreadsUserAccess.AS("thread_user_access").CreatedAt,
						TargetID:  table.FivenetMsgsThreadsUserAccess.AS("thread_user_access").ThreadID,
						Access:    table.FivenetMsgsThreadsUserAccess.AS("thread_user_access").Access,
					},
					UserId: table.FivenetMsgsThreadsUserAccess.AS("thread_user_access").UserID,
				},
			),
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMailerServiceServer(srv, s)
}
