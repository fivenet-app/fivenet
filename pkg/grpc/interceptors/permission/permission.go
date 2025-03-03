package grpc_permission

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
)

// Permission*Func is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.PermissionInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type (
	PermissionUnaryFunc  func(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error)
	PermissionStreamFunc func(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error)
)

// ServiceUnaryPermissionFuncOverride
type ServiceUnaryPermissionFuncOverride interface {
	PermissionUnaryFuncOverride(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error)
}

// ServiceStreamPermissionFuncOverride
type ServiceStreamPermissionFuncOverride interface {
	PermissionStreamFuncOverride(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error)
}

// GetPermsRemap
type GetPermsRemapFunc interface {
	GetPermsRemap() map[string]string
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request permission checks.
func UnaryServerInterceptor(permissionFunc PermissionUnaryFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceUnaryPermissionFuncOverride); ok {
			newCtx, err = overrideSrv.PermissionUnaryFuncOverride(ctx, info)
		} else {
			newCtx, err = permissionFunc(ctx, info)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request permission checks.
func StreamServerInterceptor(permissionFunc PermissionStreamFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceStreamPermissionFuncOverride); ok {
			newCtx, err = overrideSrv.PermissionStreamFuncOverride(stream.Context(), srv, info)
		} else {
			newCtx, err = permissionFunc(stream.Context(), srv, info)
		}
		if err != nil {
			return err
		}
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
