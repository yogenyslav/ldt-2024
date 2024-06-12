package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
)

type authController interface {
	Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error)
}

// Handler имплементация сервиса авторизации.
type Handler struct {
	ctrl      authController
	validator *validator.Validate
}

// New создает новый Handler.
func New(ctrl authController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
