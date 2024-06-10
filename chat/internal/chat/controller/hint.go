package controller

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
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

	newPrompt := prompt + "\nhint: " + params.Prompt

	in := &pb.ExtractReq{Prompt: newPrompt}
	meta, err := ctrl.prompter.Extract(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract meta from prompt")
		return err
	}

	return ctrl.repo.UpdateQuery(ctx, model.QueryDao{
		ID:      queryID,
		Prompt:  newPrompt,
		Status:  shared.StatusPending,
		Product: meta.GetProduct(),
		Type:    shared.QueryType(meta.GetType()),
	})
}
