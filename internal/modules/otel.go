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
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

var TracerProviderModule = fx.Module("tracer_provider",
	fx.Provide(NewTracerProvider),
)

type TracingParams struct {
	fx.In

	LC fx.Lifecycle

	Config *config.Config
}

type TracingResults struct {
	fx.Out

	TP            *tracesdk.TracerProvider
	MeterProvider otelmetric.MeterProvider
}

// NewTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the configured exporter type that will send spans to the provided url. The
// returned TracerProvider will also use a Resource configured with all the
// information about the application.
func NewTracerProvider(p TracingParams) (TracingResults, error) {
	if !p.Config.OTLP.Enabled {
		return TracingResults{
			TP: tracesdk.NewTracerProvider(),
		}, nil
	}

	ctx := context.Background()

	res, err := newAttributes(p.Config.OTLP)
	if err != nil {
		return TracingResults{}, fmt.Errorf("failed to create attributes for tracer provider. %w", err)
	}

	// Log Exporter - set in the logger module

	// Metrics
	metricExporter, err := newMetricsExporter(ctx, p.Config.OTLP.Type, p.Config.OTLP.URL, p.Config.OTLP.Insecure, p.Config.OTLP.Timeout, p.Config.OTLP.Headers)
	if err != nil {
		return TracingResults{}, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(20*time.Second)),
		),
	)
	otel.SetMeterProvider(meterProvider)

	// Tracing
	traceExporter, err := newTraceExporter(ctx, p.Config.OTLP.Type, p.Config.OTLP.URL, p.Config.OTLP.Insecure, p.Config.OTLP.Timeout, p.Config.OTLP.Headers, p.Config.OTLP.Compression)
	if err != nil {
		return TracingResults{}, fmt.Errorf("failed to create trace exporter. %w", err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(p.Config.OTLP.Ratio)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(traceExporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 7*time.Second)
		defer cancel()

		// Gracefully shutdown components
		if err := meterProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to cleanly shut down meter provider. %w", err)
		}

		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to cleanly shut down tracing provider. %w", err)
		}

		return nil
	}))

	return TracingResults{
		TP:            tp,
		MeterProvider: meterProvider,
	}, nil
}

func newAttributes(cfg config.OTLPConfig) (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname for tracer attributes. %w", err)
	}

	customAttrs := []attribute.KeyValue{
		semconv.ServiceNamespace("fivenet"),
		semconv.HostName(hostname),
		attribute.String("environment", cfg.Environment),
	}

	for i := range cfg.Attributes {
		split := strings.SplitN(cfg.Attributes[i], "=", 2)
		if len(split) < 2 {
			return nil, fmt.Errorf("failed to parse tracing attribute (%q)", cfg.Attributes[i])
		}

		customAttrs = append(customAttrs, attribute.String(split[0], split[1]))
	}

	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		customAttrs...,
	)

	return attrs, nil
}

func newMetricsExporter(ctx context.Context, tracingType config.OtelExporter, endpoint string, insecure bool, timeout time.Duration, headers map[string]string) (metric.Exporter, error) {
	switch tracingType {
	case config.TracingExporter_OTLPGRPC:
		secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			secureOption = otlpmetricgrpc.WithInsecure()
		}
		return otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpointURL(endpoint),
			secureOption,
			otlpmetricgrpc.WithCompressor(gzip.Name),
			otlpmetricgrpc.WithTimeout(timeout),
			otlpmetricgrpc.WithHeaders(headers),
		)

	case config.TracingExporter_OTLPHTTP:
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpointURL(endpoint),
			otlpmetrichttp.WithTimeout(timeout),
			otlpmetrichttp.WithHeaders(headers),
		}

		if insecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}

		return otlpmetrichttp.New(ctx, opts...)

	case config.TracingExporter_StdoutTrace:
		fallthrough
	default:
		return stdoutmetric.New()
	}
}

func newLoggerExporter(ctx context.Context, loggerType config.OtelExporter, endpoint string, insecure bool, timeout time.Duration, headers map[string]string, compression string) (log.Exporter, error) {
	switch loggerType {
	case config.TracingExporter_OTLPGRPC:
		secureOption := otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			secureOption = otlploggrpc.WithInsecure()
		}

		logExporter, err := otlploggrpc.New(ctx,
			otlploggrpc.WithEndpointURL(endpoint),
			secureOption,
			otlploggrpc.WithCompressor(compression),
			otlploggrpc.WithTimeout(timeout),
			otlploggrpc.WithHeaders(headers),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create log exporter. %w", err)
		}

		return logExporter, nil

	case config.TracingExporter_OTLPHTTP:
		opts := []otlploghttp.Option{
			otlploghttp.WithEndpointURL(endpoint),
			otlploghttp.WithTimeout(timeout),
			otlploghttp.WithHeaders(headers),
		}

		logExporter, err := otlploghttp.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create log exporter. %w", err)
		}

		return logExporter, nil

	default:
		return stdoutlog.New()
	}
}

func newTraceExporter(ctx context.Context, tracingType config.OtelExporter, endpoint string, insecure bool, timeout time.Duration, headers map[string]string, compression string) (tracesdk.SpanExporter, error) {
	var err error
	var exporter tracesdk.SpanExporter
	switch tracingType {
	case config.TracingExporter_OTLPGRPC:
		secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			secureOption = otlptracegrpc.WithInsecure()
		}

		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpointURL(endpoint),
			secureOption,
			otlptracegrpc.WithCompressor(compression),
			otlptracegrpc.WithTimeout(timeout),
			otlptracegrpc.WithHeaders(headers),
		}

		exporter, err = otlptracegrpc.New(ctx, opts...)

	case config.TracingExporter_OTLPHTTP:
		comp := otlptracehttp.NoCompression
		if compression == "gzip" {
			comp = otlptracehttp.GzipCompression
		}

		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpointURL(endpoint),
			otlptracehttp.WithTimeout(timeout),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithCompression(comp),
		}
		if insecure {
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

	return exporter, nil
}
