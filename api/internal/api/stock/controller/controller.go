package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/stock/model"
	"go.opentelemetry.io/otel/trace"
)

type stockRepo interface {
	ListProducts(ctx context.Context) ([]model.ProductDao, error)
}

// Controller is a struct that implements stock business logic.
type Controller struct {
	repo   stockRepo
	tracer trace.Tracer
}

// New creates new Controller.
func New(repo stockRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		tracer: tracer,
	}
}
