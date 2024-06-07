package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"go.opentelemetry.io/otel/trace"
)

type chatController interface {
	InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) error
	Authorize(ctx context.Context, token string) (string, error)
}

// Handler is the chat handler
type Handler struct {
	ctrl   chatController
	tracer trace.Tracer
}

// New creates a new chat handler
func New(ctrl chatController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		tracer: tracer,
	}
}
