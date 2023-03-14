package proto

import (
	"net"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/config"
	grpc_auth "github.com/galexrt/arpanet/pkg/grpc/auth"
	grpc_permission "github.com/galexrt/arpanet/pkg/grpc/permission"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/getsentry/sentry-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// GRPC Services
	pbauth "github.com/galexrt/arpanet/proto/services/auth"
	pbcitizenstore "github.com/galexrt/arpanet/proto/services/citizenstore"
	pbcompletor "github.com/galexrt/arpanet/proto/services/completor"
	pbdispatcher "github.com/galexrt/arpanet/proto/services/dispatcher"
	pbdocstore "github.com/galexrt/arpanet/proto/services/docstore"
	pbjobs "github.com/galexrt/arpanet/proto/services/jobs"
	pblivemapper "github.com/galexrt/arpanet/proto/services/livemapper"
)

type RegisterFunc func() error

func NewGRPCServer(logger *zap.Logger, tm *auth.TokenManager, p *perms.Perms) (*grpc.Server, net.Listener) {
	// Create GRPC Server
	lis, err := net.Listen("tcp", config.C.GRPC.Listen)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	sentryRecoverFunc := func(p interface{}) (err error) {
		if e, ok := p.(error); ok {
			sentry.CaptureException(e)
		}
		return status.Errorf(codes.Internal, "%v", p)
	}
	grpcAuth := &auth.GRPCAuth{
		TM: tm,
	}
	grpcPerm := &auth.GRPCPerm{
		P: p,
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_validator.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(grpcAuth.GRPCAuthFunc),
			grpc_permission.UnaryServerInterceptor(grpcPerm.GRPCPermissionUnaryFunc),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(sentryRecoverFunc),
			),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_validator.StreamServerInterceptor(),
			grpc_auth.StreamServerInterceptor(grpcAuth.GRPCAuthFunc),
			grpc_permission.StreamServerInterceptor(grpcPerm.GRPCPermissionStreamFunc),
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(sentryRecoverFunc),
			),
		)),
	)

	// Attach our GRPC services
	pbauth.RegisterAuthServiceServer(grpcServer, pbauth.NewServer(grpcAuth, tm, p))
	pbcitizenstore.RegisterCitizenStoreServiceServer(grpcServer, pbcitizenstore.NewServer(p))
	pbcompletor.RegisterCompletorServiceServer(grpcServer, pbcompletor.NewServer(p))
	pbdispatcher.RegisterDispatcherServiceServer(grpcServer, pbdispatcher.NewServer())
	pbdocstore.RegisterDocStoreServiceServer(grpcServer, pbdocstore.NewServer(p))
	pbjobs.RegisterJobsServiceServer(grpcServer, pbjobs.NewServer())
	pblivemapper.RegisterLivemapperServiceServer(grpcServer, pblivemapper.NewServer(logger.Named("grpc_livemap"), p))

	return grpcServer, lis
}
