package pkg

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// PushSpan прокинуть span в контекст.
func PushSpan(ctx context.Context, span trace.Span) context.Context {
	traceID := span.SpanContext().TraceID().String()
	return metadata.AppendToOutgoingContext(ctx, "x-trace-id", traceID)
}

// GetTraceCtx получить контекст, содержащий span.
func GetTraceCtx(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIDVal, ok := md["x-trace-id"]
	if !ok {
		log.Debug().Msg("trace id not found")
		return nil, nil
	}

	traceIDString := traceIDVal[0]
	traceID, err := trace.TraceIDFromHex(traceIDString)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to parse trace id")
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)
	return ctx, nil
}
