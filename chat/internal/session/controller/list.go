package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (ctrl *Controller) List(ctx context.Context, username string) ([]model.SessionDto, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.List",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()

	sessionsDB, err := ctrl.repo.List(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("list sessions internal repo error")
		return nil, shared.ErrGetSession
	}

	sessions := make([]model.SessionDto, len(sessionsDB))
	for idx, s := range sessionsDB {
		sessions[idx] = s.ToDto()
	}
	return sessions, nil
}
