package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

var TracerProviderModule = fx.Module("tracer_provider",
	fx.Provide(NewTracerProvider),
)

type TracingParams struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
}

// NewTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the configured exporter type that will send spans to the provided url. The
// returned TracerProvider will also use a Resource configured with all the
// information about the application.
func NewTracerProvider(p TracingParams) (*tracesdk.TracerProvider, error) {
	if !p.Config.Tracing.Enabled {
		return tracesdk.NewTracerProvider(), nil
	}

	ctx := context.Background()

	var exporter tracesdk.SpanExporter
	var err error
	switch p.Config.Tracing.Type {
	case config.TracingExporter_OTLPTraceGRPC:
		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithEndpointURL(p.Config.Tracing.URL))

	case config.TracingExporter_OTLPTraceHTTP:
		exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(p.Config.Tracing.URL))

	case config.TracingExporter_StdoutTrace:
		fallthrough
	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to setup tracing provider. %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(p.Config.Tracing.Ratio)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fivenet"),
			attribute.String("environment", p.Config.Tracing.Environment),
		)),
	)

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to cleanly shut down tracing provider: %w", err)
		}

		return nil
	}))

	return tp, nil
}
