package notifications

import (
	"database/sql"

	pbnotifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2025/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetNotifications,
		IDColumn:        table.FivenetNotifications.ID,
		TimestampColumn: table.FivenetNotifications.CreatedAt,

		MinDays: 90,
	})
}

type Server struct {
	pbnotifications.NotificationsServiceServer

	logger   *zap.Logger
	db       *sql.DB
	p        perms.Permissions
	tm       *auth.TokenMgr
	ui       userinfo.UserInfoRetriever
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	appCfg   appconfig.IConfig
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	DB        *sql.DB
	Perms     perms.Permissions
	TM        *auth.TokenMgr
	UI        userinfo.UserInfoRetriever
	JS        *events.JSWrapper
	Enricher  *mstlystcdata.Enricher
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	s := &Server{
		logger:   p.Logger,
		db:       p.DB,
		p:        p.Perms,
		tm:       p.TM,
		ui:       p.UI,
		js:       p.JS,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbnotifications.RegisterNotificationsServiceServer(srv, s)
}

func (s *Server) GetPermsRemap() map[string]string {
	return pbnotifications.PermsRemap
}
