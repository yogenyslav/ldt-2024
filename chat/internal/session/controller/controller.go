package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"go.opentelemetry.io/otel/trace"
)

type sessionRepo interface {
	InsertOne(ctx context.Context, params model.SessionDao) error
	List(ctx context.Context, username string) ([]model.SessionDao, error)
	UpdateTitle(ctx context.Context, params model.RenameReq) error
	DeleteOne(ctx context.Context, id uuid.UUID) error
	FindMeta(ctx context.Context, id uuid.UUID) (model.SessionMeta, error)
	FindContent(ctx context.Context, id uuid.UUID) ([]model.SessionContentDao, error)
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
