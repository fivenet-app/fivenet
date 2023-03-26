package proto

import (
	"database/sql"
	"net"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/config"
	grpc_auth "github.com/galexrt/arpanet/pkg/grpc/auth"
	grpc_permission "github.com/galexrt/arpanet/pkg/grpc/permission"
	"github.com/galexrt/arpanet/pkg/mstlystcdata"
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
	pbdmv "github.com/galexrt/arpanet/proto/services/dmv"
	pbdocstore "github.com/galexrt/arpanet/proto/services/docstore"
	pbjobs "github.com/galexrt/arpanet/proto/services/jobs"
	pblivemapper "github.com/galexrt/arpanet/proto/services/livemapper"
	pbnotificator "github.com/galexrt/arpanet/proto/services/notificator"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "Overview", Name: "View"},
	})
}

type RegisterFunc func() error

func NewGRPCServer(logger *zap.Logger, db *sql.DB, tm *auth.TokenManager, p *perms.Perms) (*grpc.Server, net.Listener) {
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
	grpcAuth := auth.NewGRPCAuth(tm)
	grpcPerm := auth.NewGRPCPerms(p)
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

	// "Mostly Static Data" Cache
	cache, err := mstlystcdata.NewCache(logger.Named("mstlystcdata"), db)
	if err != nil {
		logger.Fatal("failed to create mostly static data cache", zap.Error(err))
	}
	cache.Start()

	// Data enricher helper
	enricher := mstlystcdata.NewEnricher(cache)

	// Attach our GRPC services
	pbauth.RegisterAuthServiceServer(grpcServer, pbauth.NewServer(db, grpcAuth, tm, p))
	pbcitizenstore.RegisterCitizenStoreServiceServer(grpcServer, pbcitizenstore.NewServer(db, p, enricher))
	pbcompletor.RegisterCompletorServiceServer(grpcServer, pbcompletor.NewServer(db, p, cache))
	pbdocstore.RegisterDocStoreServiceServer(grpcServer, pbdocstore.NewServer(db, p, enricher))
	pbjobs.RegisterJobsServiceServer(grpcServer, pbjobs.NewServer())
	livemapper := pblivemapper.NewServer(logger.Named("grpc_livemap"), db, p)
	pblivemapper.RegisterLivemapperServiceServer(grpcServer, livemapper)
	go livemapper.GenerateRandomUserMarker()
	pbnotificator.RegisterNotificatorServiceServer(grpcServer, pbnotificator.NewServer(logger.Named("grpc_notificator"), db, p))
	pbdmv.RegisterDMVServiceServer(grpcServer, pbdmv.NewServer(db, p, enricher))

	return grpcServer, lis
}
