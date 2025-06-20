import { defineNuxtPlugin } from '#app';
import type { ReadableSpan, SpanProcessor } from '@opentelemetry/sdk-trace-base';

const samplePercentage = 0.01; // 1% sampling rate

export default defineNuxtPlugin(async (nuxtApp) => {
    await nuxtApp.$appConfigPromise;

    const { system } = useAppConfig();

    if (!system.otlp?.enabled || !system.otlp?.url) return;

    // Lazy load OpenTelemetry and set up the tracer
    const { WebTracerProvider, TraceIdRatioBasedSampler } = await import('@opentelemetry/sdk-trace-web');
    const { OTLPTraceExporter } = await import('@opentelemetry/exporter-trace-otlp-http');
    const { BatchSpanProcessor } = await import('@opentelemetry/sdk-trace-base');
    const { ZoneContextManager } = await import('@opentelemetry/context-zone');
    const { registerInstrumentations } = await import('@opentelemetry/instrumentation');
    const { FetchInstrumentation } = await import('@opentelemetry/instrumentation-fetch');
    const { DocumentLoadInstrumentation } = await import('@opentelemetry/instrumentation-document-load');
    const { trace } = await import('@opentelemetry/api');
    const { defaultResource, resourceFromAttributes } = await import('@opentelemetry/resources');
    const { ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION } = await import('@opentelemetry/semantic-conventions');

    const resource = defaultResource().merge(
        resourceFromAttributes({
            [ATTR_SERVICE_NAME]: 'fivenet',
            [ATTR_SERVICE_VERSION]: APP_VERSION,
        }),
    );

    // Configure your OTLP/Collector endpoint
    const exporter = new OTLPTraceExporter({
        url: system.otlp.url,
        concurrencyLimit: 5,
        headers: system.otlp.headers ?? {},
        timeoutMillis: 15000, // 15 seconds
    });

    const provider = new WebTracerProvider({
        resource: resource,
        spanProcessors: [new DynamicAttributesSpanProcessor(), new BatchSpanProcessor(exporter)],
        sampler: new TraceIdRatioBasedSampler(samplePercentage),
    });

    provider.register({
        contextManager: new ZoneContextManager(),
    });

    // Instrument fetch/XHR and document load
    registerInstrumentations({
        instrumentations: [
            new FetchInstrumentation({
                ignoreNetworkEvents: true,
                ignoreUrls: [/^\/images\/livemap\/.*/],
            }),
            new DocumentLoadInstrumentation(),
            // Add other instrumentations as needed
        ],
    });

    // Example: custom span on navigation
    nuxtApp.hook('page:finish', () => {
        const span = trace.getTracer('nuxt-app').startSpan('pageview', {
            attributes: {
                url: window.location.href,
                title: document.title,
            },
        });
        span.end();
    });
});

// Custom span processor that adds global attributes to every span
class DynamicAttributesSpanProcessor implements SpanProcessor {
    onStart(span: ReadableSpan): void {
        const { activeChar } = useAuth();

        // Attach job info
        if (activeChar.value) {
            span.attributes['userinfo.job'] = activeChar.value.job;
            span.attributes['userinfo.job_grade'] = activeChar.value.jobGrade;
        }
    }

    onEnd(_span: ReadableSpan): void {}

    shutdown() {
        return Promise.resolve();
    }

    forceFlush(): Promise<void> {
        return Promise.resolve();
    }
}
