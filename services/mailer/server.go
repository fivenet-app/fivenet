package mailer

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetMailerEmails,
		IDColumn:        table.FivenetMailerEmails.ID,
		JobColumn:       table.FivenetMailerEmails.Job,
		DeletedAtColumn: table.FivenetMailerEmails.DeletedAt,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:           table.FivenetMailerThreads,
				IDColumn:        table.FivenetMailerThreads.ID,
				ForeignKey:      table.FivenetMailerThreads.CreatorEmailID,
				DeletedAtColumn: table.FivenetMailerThreads.DeletedAt,

				MinDays: 60,
			},
			{
				Table:           table.FivenetMailerMessages,
				IDColumn:        table.FivenetMailerMessages.ID,
				ForeignKey:      table.FivenetMailerMessages.SenderID,
				DeletedAtColumn: table.FivenetMailerMessages.DeletedAt,

				MinDays: 60,
			},
			{
				Table:           table.FivenetMailerTemplates,
				IDColumn:        table.FivenetMailerTemplates.ID,
				ForeignKey:      table.FivenetMailerTemplates.EmailID,
				DeletedAtColumn: table.FivenetMailerTemplates.DeletedAt,

				MinDays: 60,
			},
		},
	})
}

type Server struct {
	pbmailer.MailerServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
		js       *events.JSWrapper

	access *access.Grouped[mailer.JobAccess, *mailer.JobAccess, mailer.UserAccess, *mailer.UserAccess, mailer.QualificationAccess, *mailer.QualificationAccess, mailer.AccessLevel]
}

type Params struct {
	fx.In

	DB       *sql.DB
	P        perms.Permissions
	Enricher *mstlystcdata.UserAwareEnricher
		JS       *events.JSWrapper
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.P,
		enricher: p.Enricher,
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
				table.FivenetMailerEmailsAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.ID,
						TargetID: table.FivenetMailerEmailsAccess.TargetID,
						Access:   table.FivenetMailerEmailsAccess.Access,
					},
					Job:          table.FivenetMailerEmailsAccess.Job,
					MinimumGrade: table.FivenetMailerEmailsAccess.MinimumGrade,
				},
				table.FivenetMailerEmailsAccess.AS("job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.AS("job_access").ID,
						TargetID: table.FivenetMailerEmailsAccess.AS("job_access").TargetID,
						Access:   table.FivenetMailerEmailsAccess.AS("job_access").Access,
					},
					Job:          table.FivenetMailerEmailsAccess.AS("job_access").Job,
					MinimumGrade: table.FivenetMailerEmailsAccess.AS("job_access").MinimumGrade,
				},
			),
			access.NewUsers[mailer.UserAccess, *mailer.UserAccess, mailer.AccessLevel](
				table.FivenetMailerEmailsAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.ID,
						TargetID: table.FivenetMailerEmailsAccess.TargetID,
						Access:   table.FivenetMailerEmailsAccess.Access,
					},
					UserID: table.FivenetMailerEmailsAccess.UserID,
				},
				table.FivenetMailerEmailsAccess.AS("user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.AS("user_access").ID,
						TargetID: table.FivenetMailerEmailsAccess.AS("user_access").TargetID,
						Access:   table.FivenetMailerEmailsAccess.AS("user_access").Access,
					},
					UserID: table.FivenetMailerEmailsAccess.AS("user_access").UserID,
				},
			),
			access.NewQualifications[mailer.QualificationAccess, *mailer.QualificationAccess, mailer.AccessLevel](
				table.FivenetMailerEmailsAccess,
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.ID,
						TargetID: table.FivenetMailerEmailsAccess.TargetID,
						Access:   table.FivenetMailerEmailsAccess.Access,
					},
					QualificationID: table.FivenetMailerEmailsAccess.QualificationID,
				},
				table.FivenetMailerEmailsAccess.AS("qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID: table.FivenetMailerEmailsAccess.AS("qualification_access").ID,
						TargetID: table.FivenetMailerEmailsAccess.AS(
							"qualification_access",
						).TargetID,
						Access: table.FivenetMailerEmailsAccess.AS("qualification_access").Access,
					},
					QualificationID: table.FivenetMailerEmailsAccess.AS(
						"qualification_access",
					).QualificationID,
				},
			),
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbmailer.RegisterMailerServiceServer(srv, s)
}

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbmailer.PermsRemap
}
