package middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	tele "gopkg.in/telebot.v3"
)

// Tracing трассировка.
func Tracing(machine *state.Machine, tracer trace.Tracer) tele.MiddlewareFunc {
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

			_, err := machine.GetState(ctx, c.Sender().ID)
			if errors.Is(err, redis.Nil) {
				if err = machine.SetState(ctx, c.Sender().ID, shared.StateWaitsAuth); err != nil {
					log.Error().Err(err).Msg("failed to set state")
					return err
				}
				return next(c)
			}
			if err != nil {
				log.Error().Err(err).Msg("failed to get state")
				return err
			}

			return next(c)
		}
	}
}
