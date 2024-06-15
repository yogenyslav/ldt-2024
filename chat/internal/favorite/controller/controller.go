package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"go.opentelemetry.io/otel/trace"
)

type favoriteRepo interface {
	InsertOne(ctx context.Context, params model.FavoriteDao) error
	List(ctx context.Context, username string) ([]model.FavoriteDao, error)
	FindOne(ctx context.Context, queryID int64) (model.FavoriteDao, error)
	UpdateOne(ctx context.Context, queryID int64, username string, response []byte) error
	DeleteOne(ctx context.Context, queryID int64, username string) error
}

// Controller контроллер избранных предиктов.
type Controller struct {
	repo   favoriteRepo
	tracer trace.Tracer
}

// New создает новый Controller.
func New(repo favoriteRepo, tracer trace.Tracer) *Controller {
	return &Controller{repo: repo, tracer: tracer}
}
