package pkg

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// PushSpan pushes the span to the context.
func PushSpan(ctx context.Context, span trace.Span) context.Context {
	traceID := span.SpanContext().TraceID().String()
	return metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID)
}

// PushToken pushes access token to the context.
func PushToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", token)
}
