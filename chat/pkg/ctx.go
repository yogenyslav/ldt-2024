package pkg

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	tele "gopkg.in/telebot.v3"
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

// GetTraceCtx получить контекст с trace.
func GetTraceCtx(c tele.Context) context.Context {
	return c.Get(shared.TraceCtxKey).(context.Context)
}
