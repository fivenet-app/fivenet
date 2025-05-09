package modules

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
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
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpointURL(p.Config.Tracing.URL),
			otlptracegrpc.WithTimeout(p.Config.Tracing.Timeout),
		}
		if p.Config.Tracing.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		exporter, err = otlptracegrpc.New(ctx, opts...)

	case config.TracingExporter_OTLPTraceHTTP:
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpointURL(p.Config.Tracing.URL),
			otlptracehttp.WithTimeout(p.Config.Tracing.Timeout),
		}
		if p.Config.Tracing.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		exporter, err = otlptracehttp.New(ctx, opts...)

	case config.TracingExporter_StdoutTrace:
		fallthrough
	default:
		exporter, err = stdouttrace.New()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to setup tracing provider. %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname for tracer attributes. %w", err)
	}

	customAttrs := []attribute.KeyValue{
		semconv.ServiceNamespace("fivenet"),
		semconv.HostName(hostname),
		attribute.String("environment", p.Config.Tracing.Environment),
	}

	for i := range p.Config.Tracing.Attributes {
		split := strings.SplitN(p.Config.Tracing.Attributes[i], "=", 2)
		if len(split) < 2 {
			return nil, fmt.Errorf("failed to parse tracing attribute (%q)", p.Config.Tracing.Attributes[i])
		}

		customAttrs = append(customAttrs, attribute.String(split[0], split[1]))
	}

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		customAttrs...,
	)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(p.Config.Tracing.Ratio)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(attrs),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

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
