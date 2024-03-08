package server

import (
	"github.com/gin-gonic/gin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InjectToHeaders(tracer *tracesdk.TracerProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		if span := trace.SpanContextFromContext(c.Request.Context()); span.IsSampled() {
			c.Header("X-Trace-Id", span.TraceID().String())
		}

		c.Next()
	}
}
