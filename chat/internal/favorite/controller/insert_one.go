package controller

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/pkg"
)

// InsertOne добавляет новый предикт в избранное.
func (ctrl *Controller) InsertOne(ctx context.Context, params model.FavoriteCreateReq, username string) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.InsertOne")
	defer span.End()

	response, err := json.Marshal(params.Response)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return err
	}

	if err = ctrl.repo.InsertOne(ctx, model.FavoriteDao{
		QueryID:  params.QueryID,
		Username: username,
		Response: response,
	}); err != nil {
		if pkg.CheckDuplicateKey(err) {
			if err = ctrl.repo.RestoreOne(ctx, model.FavoriteDao{
				QueryID:  params.QueryID,
				Username: username,
				Response: response,
			}); err != nil {
				log.Error().Err(err).Msg("failed to restore favorite")
				return err
			}
			return nil
		}
		log.Error().Err(err).Msg("failed to insert favorite")
		return shared.ErrCreateFavorite
	}
	return nil
}
