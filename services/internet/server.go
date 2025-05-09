package internet

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(
		&housekeeper.Table{
			Table:           table.FivenetInternetPages,
			TimestampColumn: table.FivenetInternetPages.DeletedAt,
			MinDays:         60,
		},
		&housekeeper.Table{
			Table:           table.FivenetInternetAds,
			TimestampColumn: table.FivenetInternetAds.DeletedAt,
			MinDays:         60,
		},
	)
}

type Server struct {
	pbinternet.InternetServiceServer
	pbinternet.DomainServiceServer
	pbinternet.AdsServiceServer

	db  *sql.DB
	aud audit.IAuditer

	access *access.Grouped[internet.PageJobAccess, *internet.PageJobAccess, internet.PageUserAccess, *internet.PageUserAccess, access.DummyQualificationAccess[internet.AccessLevel], *access.DummyQualificationAccess[internet.AccessLevel], internet.AccessLevel]
}

type Params struct {
	fx.In

	DB  *sql.DB
	Aud audit.IAuditer
}

func NewServer(p Params) *Server {
	s := &Server{
		db:  p.DB,
		aud: p.Aud,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbinternet.RegisterDomainServiceServer(srv, s)
	pbinternet.RegisterInternetServiceServer(srv, s)
	pbinternet.RegisterAdsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbinternet.PermsRemap
}
