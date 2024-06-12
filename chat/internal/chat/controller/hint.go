package controller

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Hint добавляет подсказку к запросу.
func (ctrl *Controller) Hint(ctx context.Context, queryID int64, params model.QueryCreateReq) (model.QueryDto, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Hint",
		trace.WithAttributes(
			attribute.Int64("queryID", queryID),
			attribute.String("hint", params.Prompt),
		),
	)
	defer span.End()

	var query model.QueryDto
	if params.Prompt == "" {
		return query, shared.ErrEmptyQueryHint
	}

	prompt, err := ctrl.cr.FindQueryPrompt(ctx, queryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return query, shared.ErrNoQueryWithID
		}
		log.Error().Err(err).Msg("failed to get query prompt")
		return query, shared.ErrGetQuery
	}

	newPrompt := prompt + "\nhint: " + params.Prompt
	meta, err := ctrl.extractMeta(ctx, newPrompt)
	if err != nil {
		return query, err
	}

	if err := ctrl.cr.UpdateQuery(ctx, model.QueryDao{
		ID:      queryID,
		Prompt:  newPrompt,
		Status:  shared.StatusPending,
		Product: meta.GetProduct(),
		Type:    shared.QueryType(meta.GetType()),
		Period:  meta.GetPeriod(),
	}); err != nil {
		return query, err
	}

	query.ID = queryID
	query.Prompt = newPrompt
	query.Status = shared.StatusPending.ToString()
	query.Product = meta.GetProduct()
	query.Type = shared.QueryType(meta.GetType()).ToString()
	query.CreatedAt = time.Now()
	query.Period = meta.GetPeriod()

	return query, nil
}
