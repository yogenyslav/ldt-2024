package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertQuery добавляет запрос в базу данных.
func (ctrl *Controller) InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) (model.QueryDto, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertQuery",
		trace.WithAttributes(
			attribute.String("username", username),
			attribute.String("query", params.Prompt),
			attribute.String("sessionID", sessionID.String()),
		),
	)
	defer span.End()

	var query model.QueryDto
	if params.Prompt == "" {
		return query, shared.ErrEmptyQueryHint
	}

	tx, err := ctrl.cr.BeginTx(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return query, shared.ErrBeginTx
	}
	defer func() {
		err = ctrl.cr.RollbackTx(tx)
		if err != nil {
			log.Error().Err(err).Msg("failed to rollback transaction")
		}
	}()

	queryID, err := ctrl.cr.InsertQuery(tx, model.QueryDao{
		SessionID: sessionID,
		Prompt:    params.Prompt,
		Username:  username,
		Status:    shared.StatusPending,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to insert query")
		return query, shared.ErrCreateQuery
	}

	tx = pkg.PushSpan(tx, span)
	meta, err := ctrl.extractMeta(tx, params.Prompt)
	if err != nil {
		return query, err
	}

	if err := ctrl.cr.UpdateQueryMeta(tx, model.QueryMeta{
		Product: meta.GetProduct(),
		Type:    shared.QueryType(meta.GetType()),
		Period:  meta.GetPeriod(),
	}, queryID); err != nil {
		return query, err
	}

	if err := ctrl.InsertResponse(tx, queryID); err != nil {
		return query, err
	}
	if err := ctrl.cr.CommitTx(tx); err != nil {
		log.Error().Err(err).Msg("failed to commit transaction")
		return query, shared.ErrCommitTx
	}

	query.ID = queryID
	query.Type = shared.QueryType(meta.GetType()).ToString()
	query.Product = meta.GetProduct()
	query.Status = shared.StatusPending.ToString()
	query.CreatedAt = time.Now()
	query.Period = meta.GetPeriod()

	return query, nil
}
