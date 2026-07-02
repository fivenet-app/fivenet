package modules

import (
	"context"
	"fmt"
	"net"

	"buf.build/go/protovalidate"
	grpcsvc "github.com/fivenet-app/fivenet/v2026/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	grpc_auth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/auth"
	grpc_permission "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/permission"
	grpc_sanitizer "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/sanitizer"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type GRPCServerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger

	GRPCAuth *auth.GRPCAuth
	GRPCPerm *auth.GRPCPerm

	Services []grpcsvc.Service `group:"grpcservices"`
}

func TestGRPCServer(
	_ context.Context,
) (*grpc.ClientConn, func(p GRPCServerParams) (*grpc.Server, error), error) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	conn, err := grpc.NewClient("passthrough:///Non-Existent.Server:80",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to test grpc server. %w", err)
	}

	// Create a Protovalidate Validator
	validator, err := protovalidate.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create protovalidate validator. %w", err)
	}

	return conn, func(p GRPCServerParams) (*grpc.Server, error) {
		srv := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				grpc_auth.UnaryServerInterceptor(p.GRPCAuth.GRPCAuthFunc),
				grpc_permission.UnaryServerInterceptor(p.GRPCPerm.GRPCPermissionUnaryFunc),
				protovalidate_middleware.UnaryServerInterceptor(validator),
				grpc_sanitizer.UnaryServerInterceptor(),
				recovery.UnaryServerInterceptor(),
			),
			grpc.ChainStreamInterceptor(
				grpc_auth.StreamServerInterceptor(p.GRPCAuth.GRPCAuthFunc),
				grpc_permission.StreamServerInterceptor(p.GRPCPerm.GRPCPermissionStreamFunc),
				protovalidate_middleware.StreamServerInterceptor(validator),
				grpc_sanitizer.StreamServerInterceptor(),
				recovery.StreamServerInterceptor(),
			),
		)

		for _, service := range p.Services {
			if service == nil {
				continue
			}

			service.RegisterServer(srv)
		}

		p.LC.Append(fx.StartHook(func() error {
			go func() {
				if err := srv.Serve(lis); err != nil {
					p.Logger.Error("error serving test grpc server", zap.Error(err))
					return
				}
			}()

			return nil
		}))

		p.LC.Append(fx.StopHook(func() error {
			err := lis.Close()
			if err != nil {
				return err
			}
			srv.Stop()

			return nil
		}))

		return srv, nil
	}, nil
}
