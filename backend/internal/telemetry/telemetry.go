// Package telemetry provides OpenTelemetry tracing setup and utilities.
package telemetry

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracer oteltrace.Tracer

// InitTracing initializes OpenTelemetry tracing.
func InitTracing(ctx context.Context, serviceName, serviceVersion string, log *slog.Logger) (func(), error) {
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(newResource(ctx, serviceName, serviceVersion)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = tp.Tracer("github.com/zorcal/sbgfit/backend")

	cleanup := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.ErrorContext(ctx, "Failed to shutdown OpenTelemetry tracer provider", "error", err)
		}
	}

	return cleanup, nil
}

// GetTraceID retrieves the trace ID from the current span in the context.
func GetTraceID(ctx context.Context) string {
	spanctx := oteltrace.SpanFromContext(ctx).SpanContext()
	if !spanctx.IsValid() {
		return ""
	}
	return spanctx.TraceID().String()
}

// GetSpanID retrieves the span ID from the current span in the context.
func GetSpanID(ctx context.Context) string {
	spanctx := oteltrace.SpanFromContext(ctx).SpanContext()
	if !spanctx.IsValid() {
		return ""
	}
	return spanctx.SpanID().String()
}

// StartSpan starts a new span with the given name and returns the context and span.
// If the tracer hasn't been initialized (e.g., in tests), returns a no-op span.
func StartSpan(ctx context.Context, name string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	if tracer == nil {
		return ctx, oteltrace.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, name, opts...)
}

// SpanFromContext returns the current span from the context.
func SpanFromContext(ctx context.Context) oteltrace.Span {
	return oteltrace.SpanFromContext(ctx)
}

// SetBaggage sets a key-value pair in baggage that will be automatically
// propagated to all child spans and across service boundaries.
func SetBaggage(ctx context.Context, key, value string) context.Context {
	member, err := baggage.NewMember(key, value)
	if err != nil {
		return ctx
	}

	bag, err := baggage.FromContext(ctx).SetMember(member)
	if err != nil {
		return ctx
	}

	return baggage.ContextWithBaggage(ctx, bag)
}

// GetBaggage retrieves a value from baggage.
func GetBaggage(ctx context.Context, key string) string {
	return baggage.FromContext(ctx).Member(key).Value()
}

// StartSpanWithBaggageAttrs starts a new span and automatically adds all
// baggage items as span attributes. If the tracer hasn't been initialized
// (e.g., in tests), returns a no-op span.
func StartSpanWithBaggageAttrs(ctx context.Context, name string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	if tracer == nil {
		return ctx, oteltrace.SpanFromContext(ctx)
	}

	ctx, span := tracer.Start(ctx, name, opts...)

	bag := baggage.FromContext(ctx)
	for _, member := range bag.Members() {
		span.SetAttributes(attribute.String(member.Key(), member.Value()))
	}

	return ctx, span
}
