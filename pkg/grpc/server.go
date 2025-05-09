package grpc

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	grpc_auth "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/auth"
	grpc_permission "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/permission"
	grpc_sanitizer "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/sanitizer"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/tracing"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

var ErrInternalServer = common.I18nErr(codes.Internal, &common.TranslateItem{Key: "errors.general.internal_error.content"}, &common.TranslateItem{Key: "errors.general.internal_error.title"})

// Setup metric for panic recoveries
var panicsTotal = promauto.With(prometheus.DefaultRegisterer).NewCounter(prometheus.CounterOpts{
	Name: "grpc_req_panics_recovered_total",
	Help: "Total number of gRPC requests recovered from internal panic.",
})

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

	Logger   *zap.Logger
	Config   *config.Config
	TP       *tracesdk.TracerProvider
	TokenMgr *auth.TokenMgr
	UserInfo userinfo.UserInfoRetriever
	Perms    perms.Permissions

	GRPCAuth *auth.GRPCAuth
	GRPCPerm *auth.GRPCPerm

	Services []Service `group:"grpcservices"`
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

	// Setup GRPC server with custom options interceptors, and tracing
	srv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     0,
			MaxConnectionAge:      time.Duration(math.MaxInt64),
			MaxConnectionAgeGrace: time.Duration(math.MaxInt64),
			Time:                  60 * time.Minute,
			Timeout:               20 * time.Second,
		}),
		grpc.MaxConcurrentStreams(128),

		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(p.TP))),

		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logging.UnaryServerInterceptor(InterceptorLogger(p.Logger),
				logging.WithFieldsFromContext(extraLogFields),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			tracing.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(p.GRPCAuth.GRPCAuthFunc),
			grpc_permission.UnaryServerInterceptor(p.GRPCPerm.GRPCPermissionUnaryFunc),
			validator.UnaryServerInterceptor(),
			grpc_sanitizer.UnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler(p.Logger)),
			),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			logging.StreamServerInterceptor(InterceptorLogger(p.Logger),
				logging.WithFieldsFromContext(extraLogFields),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
			tracing.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(p.GRPCAuth.GRPCAuthFunc),
			grpc_permission.StreamServerInterceptor(p.GRPCPerm.GRPCPermissionStreamFunc),
			validator.StreamServerInterceptor(),
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandler(grpcPanicRecoveryHandler(p.Logger)),
			),
		),
	)
	grpclog.SetLoggerV2(zapgrpc.NewLogger(p.Logger))

	for _, service := range p.Services {
		if service == nil {
			continue
		}

		service.RegisterServer(srv)
	}

	return ServerResult{
		Server: srv,
	}, nil
}

func grpcPanicRecoveryHandler(logger *zap.Logger) recovery.RecoveryHandlerFunc {
	return func(pa any) (err error) {
		panicsTotal.Inc()

		if e, ok := pa.(error); ok {
			logger.Error("recovered from panic", zap.Error(e))
		} else {
			logger.Error("recovered from panic", zap.Any("err", pa), zap.Stack("stacktrace"))
		}

		return ErrInternalServer
	}
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
