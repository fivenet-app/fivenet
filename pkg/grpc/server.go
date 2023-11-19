package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	grpc_auth "github.com/galexrt/fivenet/pkg/grpc/interceptors/auth"
	grpc_permission "github.com/galexrt/fivenet/pkg/grpc/interceptors/permission"
	grpc_sanitizer "github.com/galexrt/fivenet/pkg/grpc/interceptors/sanitizer"
	"github.com/galexrt/fivenet/pkg/grpc/interceptors/tracing"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/getsentry/sentry-go"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalServer = status.Error(codes.Internal, "errors.general.internal_error.title;errors.general.internal_error.content")
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("grpc_server")
}

var ServerModule = fx.Module("grpcserver",
	fx.Provide(
		NewServer,
	),
	fx.Decorate(wrapLogger),
)

type ServerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	Config   *config.Config
	DB       *sql.DB
	TP       *tracesdk.TracerProvider
	Services []Service `group:"grpcservices"`
	TokenMgr *auth.TokenMgr
	UserInfo userinfo.UserInfoRetriever
	Perms    perms.Permissions
}

type ServerResult struct {
	fx.Out

	Server *grpc.Server
}

func NewServer(p ServerParams) (ServerResult, error) {
	extraLogFields := func(ctx context.Context) logging.Fields {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return logging.Fields{"traceID", span.TraceID().String()}
		}
		return nil
	}

	// Setup metrics
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	prometheus.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	// Setup GRPC tracing
	otel.SetTracerProvider(p.TP)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Setup metric for panic recoveries
	panicsTotal := promauto.With(prometheus.DefaultRegisterer).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	grpcPanicRecoveryHandler := func(pa any) (err error) {
		panicsTotal.Inc()

		if e, ok := pa.(error); ok {
			p.Logger.Error("recovered from panic", zap.Error(e))
			sentry.CaptureException(e)
		} else {
			p.Logger.Error("recovered from panic", zap.Any("err", pa), zap.Stack("stacktrace"))
		}

		return ErrInternalServer
	}

	grpcAuth := auth.NewGRPCAuth(p.UserInfo, p.TokenMgr)
	grpcPerm := auth.NewGRPCPerms(p.Perms)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logging.UnaryServerInterceptor(InterceptorLogger(p.Logger),
				logging.WithFieldsFromContext(extraLogFields),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			tracing.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(grpcAuth.GRPCAuthFunc),
			grpc_permission.UnaryServerInterceptor(grpcPerm.GRPCPermissionUnaryFunc),
			validator.UnaryServerInterceptor(),
			grpc_sanitizer.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler),
			),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logging.StreamServerInterceptor(InterceptorLogger(p.Logger),
				logging.WithFieldsFromContext(extraLogFields),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			tracing.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(grpcAuth.GRPCAuthFunc),
			grpc_permission.StreamServerInterceptor(grpcPerm.GRPCPermissionStreamFunc),
			validator.StreamServerInterceptor(),
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler),
			),
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             4 * time.Second, // If a client pings more than once every 4 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Minute, // If a client is idle for 15 minute, send a GOAWAY
			MaxConnectionAge:      20 * time.Minute, // If any connection is alive for more than 20 minutes, send a GOAWAY
			MaxConnectionAgeGrace: 15 * time.Second, // Allow 15 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  20 * time.Second, // Ping the client if it is idle for 20 seconds to ensure the connection is still active
			Timeout:               7 * time.Second,  // Wait 7 second for the ping ack before assuming the connection is dead
		}),
	)

	for _, service := range p.Services {
		service.RegisterServer(srv)
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", p.Config.GRPC.Listen)
			if err != nil {
				return err
			}
			p.Logger.Info("grpc server listening", zap.String("address", p.Config.GRPC.Listen))
			go func() {
				if err := srv.Serve(ln); err != nil {
					p.Logger.Error("failed to serve grpc server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			go func() {
				srv.GracefulStop()
			}()
			// Wait 3 seconds before "forceful stop
			time.Sleep(3 * time.Second)

			srv.Stop()
			return nil
		},
	})

	return ServerResult{
		Server: srv,
	}, nil
}

// InterceptorLogger adapts zap logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
