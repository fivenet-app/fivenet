package mailer

import (
	"database/sql"

	pbmailer "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/mailer"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	mailersstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
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
	pbmailer.SettingsServiceServer
	pbmailer.ThreadServiceServer

	db       *sql.DB
	store    *mailersstore.Store
	ps       perms.Permissions
	enricher mstlystcdata.IUserAwareEnricher
	js       *events.JSWrapper

	access         *access.SubjectObjectAccess
	accessResolver *access.SubjectResolver
}

type Params struct {
	fx.In

	DB       *sql.DB
	P        perms.Permissions
	Enricher mstlystcdata.IUserAwareEnricher
	JS       *events.JSWrapper
	Store    *mailersstore.Store `optional:"true"`
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		store:    p.Store,
		ps:       p.P,
		enricher: p.Enricher,
		js:       p.JS,

		access:         access.NewMailerEmailsSubjectObjectAccess(p.DB),
		accessResolver: access.NewSubjectResolver(p.DB),
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbmailer.RegisterMailerServiceServer(srv, s)
	pbmailer.RegisterSettingsServiceServer(srv, s)
	pbmailer.RegisterThreadServiceServer(srv, s)
}
