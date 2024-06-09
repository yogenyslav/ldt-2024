package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
)

type authController interface {
	Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error)
}

// Handler is the auth handler
type Handler struct {
	ctrl      authController
	validator *validator.Validate
}

// New creates a new auth handler
func New(ctrl authController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
