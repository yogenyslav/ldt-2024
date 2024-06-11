package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"go.opentelemetry.io/otel/trace"
)

type stockController interface {
	ListProducts(ctx context.Context) ([]*pb.Product, error)
}

// Handler is a struct that implements stock server.
type Handler struct {
	pb.UnimplementedStockServer
	ctrl   stockController
	tracer trace.Tracer
}

// New creates new Handler.
func New(ctrl stockController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		tracer: tracer,
	}
}