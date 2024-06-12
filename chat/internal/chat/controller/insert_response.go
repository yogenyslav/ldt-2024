package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertResponse добавляет ответ в базу данных.
func (ctrl *Controller) InsertResponse(ctx context.Context, queryID int64) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertResponse",
		trace.WithAttributes(attribute.Int64("queryID", queryID)),
	)
	defer span.End()

	if err := ctrl.cr.InsertResponse(ctx, model.ResponseDao{
		QueryID: queryID,
	}); err != nil {
		log.Error().Err(err).Msg("failed to insert response")
		return shared.ErrCreateResponse
	}

	return nil
}
