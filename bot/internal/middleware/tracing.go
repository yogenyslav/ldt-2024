package middleware

import (
	"context"
	"fmt"

	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	tele "gopkg.in/telebot.v3"
)

// Tracing middleware for starting traces.
func Tracing(tracer trace.Tracer) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx, span := tracer.Start(
				context.Background(),
				fmt.Sprintf("%d/%d", c.Update().ID, c.Sender().ID),
				trace.WithAttributes(
					attribute.String("action", c.Text()),
					attribute.Int("updateId", c.Update().ID),
					attribute.Int64("userId", c.Sender().ID),
				),
			)
			defer span.End()

			ctx = context.WithValue(ctx, shared.UserIDKey, c.Sender().ID)
			c.Set(shared.TraceCtxKey, ctx)

			return next(c)
		}
	}
}
