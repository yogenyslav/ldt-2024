package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
)

type sessionController interface {
	NewSession(ctx context.Context, id uuid.UUID, username string, tg bool, tgID int64) error
	List(ctx context.Context, username string) ([]model.SessionDto, error)
	Rename(ctx context.Context, params model.RenameReq) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindOne(ctx context.Context, id uuid.UUID, username string) (model.FindOneResp, error)
}

// Handler имплементация сервиса сессий.
type Handler struct {
	ctrl      sessionController
	validator *validator.Validate
}

// New создает новый Handler.
func New(ctrl sessionController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
