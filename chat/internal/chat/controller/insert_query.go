package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertQuery creates new query.
//
//nolint:funlen // will be soon refactored
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

	tx, err := ctrl.repo.BeginTx(ctx)
	if err != nil {
		return query, shared.ErrBeginTx
	}
	defer func() {
		_ = ctrl.repo.RollbackTx(tx) //nolint:errcheck // transaction is either properly closed or nothing can be done
	}()

	queryID, err := ctrl.repo.InsertQuery(tx, model.QueryDao{
		SessionID: sessionID,
		Prompt:    params.Prompt,
		Username:  username,
		Status:    shared.StatusPending,
	})
	if err != nil {
		return query, shared.ErrCreateQuery
	}

	tx = pkg.PushSpan(tx, span)

	in := &pb.ExtractReq{Prompt: params.Prompt}
	meta, err := ctrl.prompter.Extract(tx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract meta from prompt")
		return query, err
	}

	if err := ctrl.repo.UpdateQueryMeta(tx, model.QueryMeta{
		Product: meta.GetProduct(),
		Type:    shared.QueryType(meta.GetType()),
	}, queryID); err != nil {
		log.Error().Err(err).Msg("failed to update query metadata")
		return query, err
	}

	if err := ctrl.InsertResponse(tx, queryID); err != nil {
		return query, err
	}
	if err := ctrl.repo.CommitTx(tx); err != nil {
		return query, shared.ErrCommitTx
	}

	query.ID = queryID
	query.Type = shared.QueryType(meta.GetType()).ToString()
	query.Product = meta.GetProduct()

	return query, nil
}
