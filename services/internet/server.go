package internet

import (
	"database/sql"

	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetInternetPages,
		IDColumn:        table.FivenetInternetPages.ID,
		DeletedAtColumn: table.FivenetInternetPages.DeletedAt,

		MinDays: 60,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetInternetAds,
		IDColumn:        table.FivenetInternetAds.ID,
		DeletedAtColumn: table.FivenetInternetAds.DeletedAt,

		MinDays: 60,
	})
}

type Server struct {
	pbinternet.InternetServiceServer
	pbinternet.DomainServiceServer
	pbinternet.AdsServiceServer

	db  *sql.DB
	aud audit.IAuditer
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

// GetPermsRemap returns the permissions re-mapping for the services.
func (s *Server) GetPermsRemap() map[string]string {
	return pbinternet.PermsRemap
}
