package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
)

func InterceptorLogger() logging.LoggerFunc {
	return func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l := log.With().Fields(fields).Logger()

		switch level {
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		default:
			l.Debug().Msg(msg)
		}
	}
}
