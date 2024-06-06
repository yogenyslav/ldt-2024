package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
)

type sessionController interface {
	NewSession(ctx context.Context, id uuid.UUID, username string) error
	List(ctx context.Context, username string) ([]model.SessionDto, error)
	Rename(ctx context.Context, params model.RenameReq) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Handler is the session handler
type Handler struct {
	ctrl      sessionController
	validator *validator.Validate
}

// New creates a new session handler
func New(ctrl sessionController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
