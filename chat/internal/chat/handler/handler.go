package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/trace"
)

type chatController interface {
	InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) (model.QueryDto, error)
	Authorize(ctx context.Context, token string) (string, error)
	Predict(ctx context.Context, out chan<- Response, cancel <-chan struct{}, queryID int64)
	Hint(ctx context.Context, queryID int64, params model.QueryCreateReq) (model.QueryDto, error)
	UpdateStatus(ctx context.Context, queryID int64, status shared.QueryStatus) error
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
