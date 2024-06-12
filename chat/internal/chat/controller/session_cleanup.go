package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// SessionCleanup checks if whether session is empty and deletes it.
func (ctrl *Controller) SessionCleanup(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.SessionCleanup")
	defer span.End()

	isEmpty, err := ctrl.sr.SessionContentEmpty(ctx, sessionID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check whether session is empty")
		return err
	}

	if isEmpty {
		if err := ctrl.sr.DeleteOne(ctx, sessionID); err != nil {
			log.Error().Err(err).Msg("failed to delete empty session")
			return err
		}
	}

	return nil
}
