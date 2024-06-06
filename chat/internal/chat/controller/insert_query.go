package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertQuery creates new query and returns its domain representation.
func (ctrl *Controller) InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertQuery",
		trace.WithAttributes(
			attribute.String("username", username),
			attribute.String("query", params.Prompt+params.Command),
			attribute.String("sessionID", sessionID.String()),
		),
	)
	defer span.End()

	queryID, err := ctrl.repo.InsertQuery(ctx, model.QueryDao{
		SessionID: sessionID,
		Prompt:    params.Prompt,
		Command:   params.Command,
		Username:  username,
	})
	if err != nil {
		return shared.ErrCreateQuery
	}

	log.Debug().Int64("queryID", queryID).Msg("query created")

	return nil
}
