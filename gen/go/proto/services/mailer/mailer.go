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

	tEmails           = table.FivenetMailerEmails.AS("email")
	tEmailsUserAccess = table.FivenetMailerEmailsUserAccess
	tEmailsJobAccess  = table.FivenetMailerEmailsJobAccess

	tThreads           = table.FivenetMailerThreads.AS("thread")
	tThreadsEmailState = table.FivenetMailerThreadsStateEmail.AS("thread_state_email")
	tThreadsUserState  = table.FivenetMailerThreadsStateUser.AS("thread_state_user")

	tMessages = table.FivenetMailerMessages.AS("message")
)

type Server struct {
	MailerServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	js       *events.JSWrapper

	access *access.Grouped[mailer.JobAccess, *mailer.JobAccess, mailer.UserAccess, *mailer.UserAccess, mailer.AccessLevel]
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

		access: access.NewGrouped[mailer.JobAccess, *mailer.JobAccess, mailer.UserAccess, *mailer.UserAccess, mailer.AccessLevel](
			p.DB,
			table.FivenetMailerEmails,
			&access.TargetTableColumns{
				ID:         table.FivenetMailerEmails.ID,
				DeletedAt:  table.FivenetMailerEmails.DeletedAt,
				CreatorJob: table.FivenetMailerEmails.Job,
			},
			access.NewJobs[mailer.JobAccess, *mailer.JobAccess, mailer.AccessLevel](
				table.FivenetMailerEmailsJobAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsJobAccess.ID,
						CreatedAt: table.FivenetMailerEmailsJobAccess.CreatedAt,
						TargetID:  table.FivenetMailerEmailsJobAccess.EmailID,
						Access:    table.FivenetMailerEmailsJobAccess.Access,
					},
					Job:          table.FivenetMailerEmailsJobAccess.Job,
					MinimumGrade: table.FivenetMailerEmailsJobAccess.MinimumGrade,
				},
				table.FivenetMailerEmailsJobAccess.AS("email_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsJobAccess.AS("email_job_access").ID,
						CreatedAt: table.FivenetMailerEmailsJobAccess.AS("email_job_access").CreatedAt,
						TargetID:  table.FivenetMailerEmailsJobAccess.AS("email_job_access").EmailID,
						Access:    table.FivenetMailerEmailsJobAccess.AS("email_job_access").Access,
					},
					Job:          table.FivenetMailerEmailsJobAccess.AS("email_job_access").Job,
					MinimumGrade: table.FivenetMailerEmailsJobAccess.AS("email_job_access").MinimumGrade,
				},
			),
			access.NewUsers[mailer.UserAccess, *mailer.UserAccess, mailer.AccessLevel](
				table.FivenetMailerEmailsUserAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsUserAccess.ID,
						CreatedAt: table.FivenetMailerEmailsUserAccess.CreatedAt,
						TargetID:  table.FivenetMailerEmailsUserAccess.EmailID,
						Access:    table.FivenetMailerEmailsUserAccess.Access,
					},
					UserId: table.FivenetMailerEmailsUserAccess.UserID,
				},
				table.FivenetMailerEmailsUserAccess.AS("email_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsUserAccess.AS("email_user_access").ID,
						CreatedAt: table.FivenetMailerEmailsUserAccess.AS("email_user_access").CreatedAt,
						TargetID:  table.FivenetMailerEmailsUserAccess.AS("email_user_access").EmailID,
						Access:    table.FivenetMailerEmailsUserAccess.AS("email_user_access").Access,
					},
					UserId: table.FivenetMailerEmailsUserAccess.AS("email_user_access").UserID,
				},
			),
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMailerServiceServer(srv, s)
}
