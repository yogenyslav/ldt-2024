package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// SessionCleanup проверяет, пуста ли сессия, и удаляет ее, если она пуста.
func (ctrl *Controller) SessionCleanup(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.SessionCleanup")
	defer span.End()

	isEmpty, err := ctrl.sr.SessionContentEmpty(ctx, sessionID)
	if err != nil {
		return shared.ErrGetSession
	}

	if isEmpty {
		if err := ctrl.sr.DeleteOne(ctx, sessionID); err != nil {
			return shared.ErrDeleteSession
		}
	}

	return nil
}
