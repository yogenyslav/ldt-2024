package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"go.opentelemetry.io/otel/trace"
)

type sessionRepo interface {
	InsertOne(ctx context.Context, params model.SessionDao) error
}

// Controller a struct that implements session business logic.
type Controller struct {
	repo   sessionRepo
	tracer trace.Tracer
}

// New creates new session Controller.
func New(repo sessionRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		tracer: tracer,
	}
}
