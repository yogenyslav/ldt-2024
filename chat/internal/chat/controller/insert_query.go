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
func (ctrl *Controller) InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) (int64, error) {
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

	tx, err := ctrl.repo.BeginTx(ctx)
	if err != nil {
		return 0, shared.ErrBeginTx
	}
	defer func() {
		_ = ctrl.repo.RollbackTx(tx) //nolint:errcheck // transaction is either properly closed or nothing can be done
	}()

	queryID, err := ctrl.repo.InsertQuery(tx, model.QueryDao{
		SessionID: sessionID,
		Prompt:    params.Prompt,
		Command:   params.Command,
		Username:  username,
	})
	if err != nil {
		return 0, shared.ErrCreateQuery
	}

	if params.Prompt != "" {
		tx = pkg.PushSpan(tx, span)

		in := &pb.ExtractReq{Prompt: params.Prompt}
		meta, err := ctrl.prompter.Extract(tx, in)
		if err != nil {
			log.Error().Err(err).Msg("failed to extract meta from prompt")
			return 0, err
		}

		if err := ctrl.repo.UpdateQueryMeta(tx, model.QueryMeta{
			Product: meta.GetProduct(),
			Type:    shared.QueryType(meta.GetType()),
		}, queryID); err != nil {
			log.Error().Err(err).Msg("failed to update query metadata")
			return 0, err
		}
	}

	if err := ctrl.InsertResponse(tx, queryID); err != nil {
		return 0, err
	}
	if err := ctrl.repo.CommitTx(tx); err != nil {
		return 0, shared.ErrCommitTx
	}

	return queryID, nil
}
