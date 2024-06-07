package controller

import (
	"context"

	"github.com/google/uuid"
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

	// if params.Prompt != ""
	// pass the query through prompt detection
	// and update value in params

	tx, err := ctrl.repo.BeginTx(ctx)
	if err != nil {
		return shared.ErrBeginTx
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
		return shared.ErrCreateQuery
	}

	if err := ctrl.InsertResponse(tx, queryID); err != nil {
		return err
	}

	if err := ctrl.repo.CommitTx(tx); err != nil {
		return shared.ErrCommitTx
	}

	// make prediction and update

	return nil
}
