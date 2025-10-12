// Package sanitizer provides a gRPC interceptor that sanitizes incoming requests.
package sanitizer

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ISanitize should be implemented by any type that wants to provide custom sanitization logic.
type ISanitize interface {
	Sanitize() error
}

// UnaryServerInterceptor returns a gRPC UnaryServerInterceptor that applies the sanitize logic
// to every incoming request that implements ISanitize. If sanitization fails, the request is rejected.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := sanitize(ctx, req); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC StreamServerInterceptor that applies the sanitize logic
// to every incoming stream request/message that implements ISanitize. If sanitization fails, the
// request is rejected.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := sanitize(stream.Context(), srv); err != nil {
			return err
		}

		// Wrap to intercept RecvMsg to be able to sanitize incoming messages
		wrapped := &sanitizeStream{ServerStream: stream}
		return handler(srv, wrapped)
	}
}

// sanitize checks if the request or response implements ISanitize and calls its Sanitize method.
// If an error is returned, it is converted to a gRPC InvalidArgument error.
func sanitize(_ context.Context, reqOrRes any) error {
	var err error

	switch v := reqOrRes.(type) {
	case ISanitize:
		err = v.Sanitize()
	}

	if err == nil {
		return nil
	}

	return status.Error(codes.InvalidArgument, err.Error())
}

type sanitizeStream struct {
	grpc.ServerStream
}

func (s *sanitizeStream) RecvMsg(m any) error {
	// Receive the next message from the client
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	// Sanitize/validate it before the handler sees it
	return sanitize(s.Context(), m)
}
