package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// FindOne возвращает предикт из избранного.
func (ctrl *Controller) FindOne(ctx context.Context, queryID int64) (model.FavoriteDto, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.FindOne")
	defer span.End()

	favorite, err := ctrl.repo.FindOne(ctx, queryID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find favorite")
		return model.FavoriteDto{}, shared.ErrGetFavorite
	}

	return favorite.ToDto(), nil
}
