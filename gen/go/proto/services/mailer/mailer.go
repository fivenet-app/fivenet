package mailer

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetMailerEmails,
		TimestampColumn: table.FivenetMailerEmails.DeletedAt,
		MinDays:         60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetMailerThreads,
		TimestampColumn: table.FivenetMailerThreads.DeletedAt,
		MinDays:         60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetMailerMessages,
		TimestampColumn: table.FivenetMailerMessages.DeletedAt,
		MinDays:         60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetMailerTemplates,
		TimestampColumn: table.FivenetMailerTemplates.DeletedAt,
		MinDays:         60,
	})
}

type Server struct {
	MailerServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	js       *events.JSWrapper

	access *access.Grouped[mailer.JobAccess, *mailer.JobAccess, mailer.UserAccess, *mailer.UserAccess, mailer.QualificationAccess, *mailer.QualificationAccess, mailer.AccessLevel]
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
		ps:       p.P,
		enricher: p.Enricher,
		aud:      p.Aud,
		js:       p.JS,

		access: access.NewGrouped(
			p.DB,
			table.FivenetMailerEmails,
			&access.TargetTableColumns{
				ID:        table.FivenetMailerEmails.ID,
				DeletedAt: table.FivenetMailerEmails.DeletedAt,
				CreatorID: table.FivenetMailerEmails.UserID,
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
				table.FivenetMailerEmailsJobAccess.AS("job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsJobAccess.AS("job_access").ID,
						CreatedAt: table.FivenetMailerEmailsJobAccess.AS("job_access").CreatedAt,
						TargetID:  table.FivenetMailerEmailsJobAccess.AS("job_access").EmailID,
						Access:    table.FivenetMailerEmailsJobAccess.AS("job_access").Access,
					},
					Job:          table.FivenetMailerEmailsJobAccess.AS("job_access").Job,
					MinimumGrade: table.FivenetMailerEmailsJobAccess.AS("job_access").MinimumGrade,
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
				table.FivenetMailerEmailsUserAccess.AS("user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsUserAccess.AS("user_access").ID,
						CreatedAt: table.FivenetMailerEmailsUserAccess.AS("user_access").CreatedAt,
						TargetID:  table.FivenetMailerEmailsUserAccess.AS("user_access").EmailID,
						Access:    table.FivenetMailerEmailsUserAccess.AS("user_access").Access,
					},
					UserId: table.FivenetMailerEmailsUserAccess.AS("user_access").UserID,
				},
			),
			access.NewQualifications[mailer.QualificationAccess, *mailer.QualificationAccess, mailer.AccessLevel](
				table.FivenetMailerEmailsQualificationsAccess,
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsQualificationsAccess.ID,
						CreatedAt: table.FivenetMailerEmailsQualificationsAccess.CreatedAt,
						TargetID:  table.FivenetMailerEmailsQualificationsAccess.EmailID,
						Access:    table.FivenetMailerEmailsQualificationsAccess.Access,
					},
					QualificationId: table.FivenetMailerEmailsQualificationsAccess.QualificationID,
				},
				table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:        table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access").ID,
						CreatedAt: table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access").CreatedAt,
						TargetID:  table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access").EmailID,
						Access:    table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access").Access,
					},
					QualificationId: table.FivenetMailerEmailsQualificationsAccess.AS("qualification_access").QualificationID,
				},
			),
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterMailerServiceServer(srv, s)
}
