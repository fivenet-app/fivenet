package server

import (
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

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func NewTracerProvider(cfg *config.Config) (*tracesdk.TracerProvider, error) {
	if !cfg.Tracing.Enabled {
		return tracesdk.NewTracerProvider(), nil
	}

	var exporter tracesdk.SpanExporter
	var err error
	// If URL is set, setup jaeger trace exporter
	if cfg.Tracing.URL != "" {
		// Create the Jaeger exporter
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Tracing.URL)))
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	}

	return tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fivenet"),
			attribute.String("environment", cfg.Tracing.Environment),
		)),
	), nil
}
