package pkg

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// PushSpan прокидывает span в контекст.
func PushSpan(ctx context.Context, span trace.Span) context.Context {
	traceID := span.SpanContext().TraceID().String()
	return metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID)
}

// PushToken производит авторизацию.
func PushToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
}
