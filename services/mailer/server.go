package mailer

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	pbmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer"
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
	pbmailer.MailerServiceServer

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
					UserId: table.FivenetMailerEmailsAccess.UserID,
				},
				table.FivenetMailerEmailsAccess.AS("user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.AS("user_access").ID,
						TargetID: table.FivenetMailerEmailsAccess.AS("user_access").TargetID,
						Access:   table.FivenetMailerEmailsAccess.AS("user_access").Access,
					},
					UserId: table.FivenetMailerEmailsAccess.AS("user_access").UserID,
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
					QualificationId: table.FivenetMailerEmailsAccess.QualificationID,
				},
				table.FivenetMailerEmailsAccess.AS("qualification_access"),
				&access.QualificationAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetMailerEmailsAccess.AS("qualification_access").ID,
						TargetID: table.FivenetMailerEmailsAccess.AS("qualification_access").TargetID,
						Access:   table.FivenetMailerEmailsAccess.AS("qualification_access").Access,
					},
					QualificationId: table.FivenetMailerEmailsAccess.AS("qualification_access").QualificationID,
				},
			),
		),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbmailer.RegisterMailerServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbmailer.PermsRemap
}
