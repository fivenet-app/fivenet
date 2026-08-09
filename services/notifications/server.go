package notifications

import (
	"context"
	"database/sql"
	"sync"

	pbnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/notifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	mailerstore "github.com/fivenet-app/fivenet/v2026/stores/mailer"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	"github.com/prometheus/client_golang/prometheus"
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

	logger      *zap.Logger
	ctx         context.Context //nolint:containedctx // Server keeps lifecycle context for stream consumer cleanup.
	db          *sql.DB
	js          *events.JSWrapper
	auth        *grpcauth.GRPCAuth
	store       notificationsstore.IStore
	mailerStore mailerstore.IStore
	metrics     *notificationMetrics
}

type notificationMetrics struct {
	activeSessions prometheus.Gauge
	lastSession    prometheus.Gauge
}

var (
	notificationMetricsOnce sync.Once
	notificationMetricsInst *notificationMetrics
)

func getNotificationMetrics() *notificationMetrics {
	notificationMetricsOnce.Do(func() {
		notificationMetricsInst = &notificationMetrics{
			activeSessions: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "user",
				Name:      "active_session_count",
				Help:      "Number of active user sessions.",
			}),
			lastSession: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "user",
				Name:      "last_session_time",
				Help:      "Timestamp of the last started user session.",
			}),
		}

		prometheus.MustRegister(
			notificationMetricsInst.activeSessions,
			notificationMetricsInst.lastSession,
		)
	})

	return notificationMetricsInst
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger      *zap.Logger
	DB          *sql.DB
	JS          *events.JSWrapper
	Auth        *grpcauth.GRPCAuth
	Store       notificationsstore.IStore
	MailerStore mailerstore.IStore
}

func NewServer(p Params) *Server {
	ctxCancel, cancel := context.WithCancel(context.Background())

	s := &Server{
		logger:      p.Logger,
		ctx:         ctxCancel,
		db:          p.DB,
		js:          p.JS,
		auth:        p.Auth,
		store:       p.Store,
		mailerStore: p.MailerStore,
		metrics:     getNotificationMetrics(),
	}

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbnotifications.RegisterNotificationsServiceServer(srv, s)
}

// AuthFuncOverride allows notifications streams to start with account-only auth
// and upgrade to full char auth when the user selects a character.
func (s *Server) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	switch fullMethod {
	case pbnotifications.NotificationsService_Stream_FullMethodName:
		if hasUserTokenInContext(ctx) {
			return s.auth.GRPCAuthFunc(ctx, fullMethod)
		}
		return s.auth.GRPCAuthFuncWithoutUserInfo(ctx, fullMethod)

	default:
		return s.auth.GRPCAuthFunc(ctx, fullMethod)
	}
}

func hasUserTokenInContext(ctx context.Context) bool {
	token, err := grpcauth.GetUserTokenFromGRPCContext(ctx)
	return err == nil && token != ""
}
