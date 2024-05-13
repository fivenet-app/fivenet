package events

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceIdHeader = "X-Trace-Id"
	spanIdHeader  = "X-Span-Id"
)

func (j *JSWrapper) addSpanInfoToMsg(ctx context.Context, msg *nats.Msg) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsSampled() {
		if msg.Header == nil {
			msg.Header = nats.Header{}
		}

		msg.Header.Set(traceIdHeader, span.SpanContext().TraceID().String())
		msg.Header.Set(spanIdHeader, span.SpanContext().SpanID().String())
	}
}

func (j *JSWrapper) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return j.PublishMsg(ctx, &nats.Msg{Subject: subject, Data: payload}, opts...)
}

func (j *JSWrapper) PublishMsg(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	j.addSpanInfoToMsg(ctx, msg)
	return j.JetStream.PublishMsg(ctx, msg, opts...)
}

func (j *JSWrapper) PublishAsync(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	return j.PublishMsgAsync(ctx, &nats.Msg{Subject: subject, Data: payload}, opts...)
}

func (j *JSWrapper) PublishMsgAsync(ctx context.Context, msg *nats.Msg, opts ...jetstream.PublishOpt) (jetstream.PubAckFuture, error) {
	j.addSpanInfoToMsg(ctx, msg)
	return j.JetStream.PublishMsgAsync(msg, opts...)
}

func GetJetstreamMsgContext(msg jetstream.Msg) (spanContext trace.SpanContext, err error) {
	headers := msg.Headers()

	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(headers.Get(traceIdHeader))
	if err != nil {
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(headers.Get(spanIdHeader))
	if err != nil {
		return spanContext, err
	}

	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = true
	spanContext = trace.NewSpanContext(spanContextConfig)

	return spanContext, nil
}

func GetNatsMsgContext(msg *nats.Msg) (spanContext trace.SpanContext, err error) {
	var traceID trace.TraceID
	traceID, err = trace.TraceIDFromHex(msg.Header.Get(traceIdHeader))
	if err != nil {
		return spanContext, err
	}
	var spanID trace.SpanID
	spanID, err = trace.SpanIDFromHex(msg.Header.Get(spanIdHeader))
	if err != nil {
		return spanContext, err
	}

	var spanContextConfig trace.SpanContextConfig
	spanContextConfig.TraceID = traceID
	spanContextConfig.SpanID = spanID
	spanContextConfig.TraceFlags = 01
	spanContextConfig.Remote = true
	spanContext = trace.NewSpanContext(spanContextConfig)

	return spanContext, nil
}
