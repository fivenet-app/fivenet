package proto

import (
	"context"
	"database/sql"
	"net"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/config"
	grpc_auth "github.com/galexrt/fivenet/pkg/grpc/auth"
	grpc_permission "github.com/galexrt/fivenet/pkg/grpc/permission"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
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
	pbauth "github.com/galexrt/fivenet/proto/services/auth"
	pbcitizenstore "github.com/galexrt/fivenet/proto/services/citizenstore"
	pbcompletor "github.com/galexrt/fivenet/proto/services/completor"
	pbdmv "github.com/galexrt/fivenet/proto/services/dmv"
	pbdocstore "github.com/galexrt/fivenet/proto/services/docstore"
	pbjobs "github.com/galexrt/fivenet/proto/services/jobs"
	pblivemapper "github.com/galexrt/fivenet/proto/services/livemapper"
	pbnotificator "github.com/galexrt/fivenet/proto/services/notificator"
	pbrector "github.com/galexrt/fivenet/proto/services/rector"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "SuperUser", Name: "AnyAccess", Description: "Super User any access to view, edit and delete."},
	})
}

type RegisterFunc func() error

func NewGRPCServer(ctx context.Context, logger *zap.Logger, db *sql.DB, tm *auth.TokenManager, p perms.Permissions) (*grpc.Server, net.Listener) {
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
	cache, err := mstlystcdata.NewCache(ctx,
		logger.Named("mstlystcdata"), db)
	if err != nil {
		logger.Fatal("failed to create mostly static data cache", zap.Error(err))
	}
	cache.Start()

	// Data enricher helper
	enricher := mstlystcdata.NewEnricher(cache)

	// Audit Storer
	audit := audit.New(db)
	_ = audit

	// Attach our GRPC services
	pbauth.RegisterAuthServiceServer(grpcServer, pbauth.NewServer(db, grpcAuth, tm, p, enricher))
	pbcitizenstore.RegisterCitizenStoreServiceServer(grpcServer, pbcitizenstore.NewServer(db, p, enricher))
	pbcompletor.RegisterCompletorServiceServer(grpcServer, pbcompletor.NewServer(db, p, cache))
	pbdocstore.RegisterDocStoreServiceServer(grpcServer, pbdocstore.NewServer(db, p, enricher))
	pbjobs.RegisterJobsServiceServer(grpcServer, pbjobs.NewServer())
	livemapper := pblivemapper.NewServer(ctx, logger.Named("grpc_livemap"), db, p, enricher)
	livemapper.Start()

	pblivemapper.RegisterLivemapperServiceServer(grpcServer, livemapper)
	pbnotificator.RegisterNotificatorServiceServer(grpcServer, pbnotificator.NewServer(logger.Named("grpc_notificator"), db, p))
	pbdmv.RegisterDMVServiceServer(grpcServer, pbdmv.NewServer(db, p, enricher))
	pbrector.RegisterRectorServiceServer(grpcServer, pbrector.NewServer(logger, db, p))

	// Only run the livemapper random user marker generator in debug mode
	if config.C.Mode == gin.DebugMode {
		go livemapper.GenerateRandomUserMarker()
	}

	return grpcServer, lis
}
