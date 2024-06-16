package controller

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
)

// UpdateOne обновляет предикт в избранном.
func (ctrl *Controller) UpdateOne(ctx context.Context, params model.FavoriteUpdateReq, username string) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.UpdateOne")
	defer span.End()

	resp, err := json.Marshal(params.Response)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return err
	}

	return ctrl.repo.UpdateOne(ctx, params.ID, username, resp)
}
