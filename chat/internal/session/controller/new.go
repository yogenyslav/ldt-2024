package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/pkg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NewSession создает новую сессию.
func (ctrl *Controller) NewSession(ctx context.Context, id uuid.UUID, username string) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.NewSession",
		trace.WithAttributes(
			attribute.String("id", id.String()),
			attribute.String("username", username),
		),
	)
	defer span.End()

	if err := ctrl.repo.InsertOne(ctx, model.SessionDao{
		ID:       id,
		Username: username,
	}); err != nil {
		if pkg.CheckDuplicateKey(err) {
			return shared.ErrSessionDuplicateID
		}
		return shared.ErrCreateSession
	}

	return nil
}
