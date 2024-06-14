package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// List возвращает список избранных предиктов.
func (ctrl *Controller) List(ctx context.Context, username string) ([]model.FavoriteDto, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.List")
	defer span.End()

	favoritesDB, err := ctrl.repo.List(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("failed to list favorites")
		return nil, shared.ErrGetFavorite
	}

	favorites := make([]model.FavoriteDto, len(favoritesDB))
	for i := 0; i < len(favoritesDB); i++ {
		favorites[i] = favoritesDB[i].ToDto()
	}
	return favorites, nil
}
