package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor returns a gRPC UnaryServerInterceptor that injects trace information into the response trailer
// and records errors on the current span if present. This enables trace propagation and error reporting for unary RPCs.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			grpc.SetTrailer(ctx, metadata.Pairs("X-Trace-Id", span.TraceID().String()))
		}

		resp, err := handler(ctx, req)
		if err != nil {
			if trace.SpanFromContext(ctx).SpanContext().IsValid() {
				trace.SpanFromContext(ctx).RecordError(err)
			}
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a gRPC StreamServerInterceptor that injects trace information into the response trailer
// and records errors on the current span if present. This enables trace propagation and error reporting for streaming RPCs.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if span := trace.SpanContextFromContext(stream.Context()); span.IsSampled() {
			grpc.SetTrailer(stream.Context(), metadata.Pairs("X-Trace-Id", span.TraceID().String()))
		}

		err := handler(srv, stream)
		if err != nil {
			if trace.SpanFromContext(stream.Context()).SpanContext().IsValid() {
				trace.SpanFromContext(stream.Context()).RecordError(err)
			}
		}

		return err
	}
}
