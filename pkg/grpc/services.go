package grpc

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Service interface {
	RegisterServer(srv *grpc.Server)
}

// AsService annotates the given constructor to state that
// it provides a GRPC service to the "grpcservices" group.
func AsService(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Service)),
		fx.ResultTags(`group:"grpcservices"`),
	)
}
