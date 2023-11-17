package server

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

var TracerProviderModule = fx.Module("tracerProvider",
	fx.Provide(NewTracerProvider),
)

type TracingParams struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func NewTracerProvider(p TracingParams) (*tracesdk.TracerProvider, error) {
	if !p.Config.Tracing.Enabled {
		return tracesdk.NewTracerProvider(), nil
	}

	var exporter tracesdk.SpanExporter
	var err error
	// If URL is set, setup jaeger trace exporter
	if p.Config.Tracing.URL != "" {
		// Create the Jaeger exporter
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(p.Config.Tracing.URL)))
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
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
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to cleanly shut down tracing provider: %w", err)
		}

		return nil
	}))

	return tp, nil
}
