package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
)

type favoriteController interface {
	InsertOne(ctx context.Context, params model.FavoriteCreateReq, username string) error
	List(ctx context.Context, username string) ([]model.FavoriteDto, error)
	FindOne(ctx context.Context, queryID int64) (model.FavoriteDto, error)
	UpdateOne(ctx context.Context, params model.FavoriteUpdateReq, username string) error
	DeleteOne(ctx context.Context, queryID int64, username string) error
}

// Handler обработчик избранных предиктов.
type Handler struct {
	ctrl favoriteController
}

// New создает новый Handler.
func New(ctrl favoriteController) *Handler {
	return &Handler{ctrl: ctrl}
}
