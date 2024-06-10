package controller

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Hint adds hint to existing prompt by id.
func (ctrl *Controller) Hint(ctx context.Context, queryID int64, params model.QueryCreateReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Hint",
		trace.WithAttributes(
			attribute.Int64("queryID", queryID),
			attribute.String("hint", params.Prompt),
		),
	)
	defer span.End()

	if params.Prompt == "" {
		return shared.ErrEmptyQueryHint
	}

	prompt, err := ctrl.repo.FindQueryPrompt(ctx, queryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrNoQueryWithID
		}
		return shared.ErrGetQuery
	}

	log.Debug().Str("initial prompt", prompt).Msg("adding hint")

	if err := ctrl.repo.UpdateQueryPrompt(ctx, queryID, prompt+"\nhint: "+params.Prompt); err != nil {
		return err
	}
	return ctrl.repo.UpdateQueryStatus(ctx, queryID, shared.StatusPending)
}