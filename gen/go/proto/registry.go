package proto

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	grpc_auth "github.com/galexrt/fivenet/pkg/grpc/interceptors/auth"
	grpc_permission "github.com/galexrt/fivenet/pkg/grpc/interceptors/permission"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// GRPC Services
	pbauth "github.com/galexrt/fivenet/gen/go/proto/services/auth"
	pbcitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore"
	pbcompletor "github.com/galexrt/fivenet/gen/go/proto/services/completor"
	pbdmv "github.com/galexrt/fivenet/gen/go/proto/services/dmv"
	pbdocstore "github.com/galexrt/fivenet/gen/go/proto/services/docstore"
	pbjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs"
	pblivemapper "github.com/galexrt/fivenet/gen/go/proto/services/livemapper"
	pbnotificator "github.com/galexrt/fivenet/gen/go/proto/services/notificator"
	pbrector "github.com/galexrt/fivenet/gen/go/proto/services/rector"
)

var (
	GenericInternalServerError = status.Error(codes.Internal, "Internal server error")
)

type RegisterFunc func() error

func NewGRPCServer(ctx context.Context, logger *zap.Logger, db *sql.DB, tm *auth.TokenMgr, p perms.Permissions, aud audit.IAuditer) (*grpc.Server, net.Listener) {
	// Create GRPC Server
	lis, err := net.Listen("tcp", config.C.GRPC.Listen)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	// Setup metrics.
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	prometheus.MustRegister(srvMetrics)

	// Setup metric for panic recoveries.
	panicsTotal := promauto.With(prometheus.DefaultRegisterer).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()

		logger.Error("recovered from panic", zap.Any("err", p), zap.Stack("stacktrace"))
		if e, ok := p.(error); ok {
			sentry.CaptureException(e)
		}

		return GenericInternalServerError
	}

	ui := userinfo.NewUIRetriever(ctx, db)
	grpcAuth := auth.NewGRPCAuth(ui, tm)
	grpcPerm := auth.NewGRPCPerms(p)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(InterceptorLogger(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			grpc_auth.UnaryServerInterceptor(grpcAuth.GRPCAuthFunc),
			validator.UnaryServerInterceptor(),
			grpc_permission.UnaryServerInterceptor(grpcPerm.GRPCPermissionUnaryFunc),
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler),
			),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(),
			logging.StreamServerInterceptor(InterceptorLogger(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			grpc_auth.StreamServerInterceptor(grpcAuth.GRPCAuthFunc),
			validator.StreamServerInterceptor(),
			grpc_permission.StreamServerInterceptor(grpcPerm.GRPCPermissionStreamFunc),
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler),
			),
		),
	)

	// "Mostly Static Data" Cache
	cache, err := mstlystcdata.NewCache(ctx,
		logger.Named("mstlystcdata"), db)
	if err != nil {
		logger.Fatal("failed to create mostly static data cache", zap.Error(err))
	}
	go cache.Start()

	// Data enricher helper
	enricher := mstlystcdata.NewEnricher(cache)

	// Notifier
	notif := notifi.New(db)

	// Attach our GRPC services
	pbauth.RegisterAuthServiceServer(grpcServer, pbauth.NewServer(db, grpcAuth, tm, p, enricher, aud, ui))
	pbcitizenstore.RegisterCitizenStoreServiceServer(grpcServer, pbcitizenstore.NewServer(db, p, enricher, aud))
	pbcompletor.RegisterCompletorServiceServer(grpcServer, pbcompletor.NewServer(db, p, cache))
	pbdocstore.RegisterDocStoreServiceServer(grpcServer, pbdocstore.NewServer(db, p, enricher, aud, notif))
	pbjobs.RegisterJobsServiceServer(grpcServer, pbjobs.NewServer())
	livemapper := pblivemapper.NewServer(ctx, logger.Named("grpc_livemap"), db, p, enricher)
	go livemapper.Start()

	pblivemapper.RegisterLivemapperServiceServer(grpcServer, livemapper)
	pbnotificator.RegisterNotificatorServiceServer(grpcServer, pbnotificator.NewServer(logger.Named("grpc_notificator"), db, p, tm))
	pbdmv.RegisterDMVServiceServer(grpcServer, pbdmv.NewServer(db, p, enricher, aud))
	pbrector.RegisterRectorServiceServer(grpcServer, pbrector.NewServer(logger, db, p, aud, enricher))

	// Only run the livemapper random user marker generator in debug mode
	if config.C.Mode == gin.DebugMode {
		go livemapper.GenerateRandomUserMarker()
		go livemapper.GenerateRandomDispatchMarker()
	}

	return grpcServer, lis
}

// InterceptorLogger adapts zap logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)
		iter := logging.Fields(fields).Iterator()
		for iter.Next() {
			k, v := iter.At()
			f = append(f, zap.Any(k, v))
		}
		l = l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			l.Error(fmt.Sprintf("unknown log level '%v' for message", lvl), zap.String("msg", msg))
		}
	})
}
