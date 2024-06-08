package pkg

import (
	"context"

	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	tele "gopkg.in/telebot.v3"
)

// PushSpan pushes the span to the context
func PushSpan(ctx context.Context, span trace.Span) context.Context {
	traceID := span.SpanContext().TraceID().String()
	return metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID)
}

// GetTraceCtx returns the context.Context from tele.Context.
func GetTraceCtx(c tele.Context) context.Context {
	return c.Get(shared.TraceCtxKey).(context.Context)
}
