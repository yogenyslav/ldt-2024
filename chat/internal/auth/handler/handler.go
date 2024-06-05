package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
	"go.opentelemetry.io/otel/trace"
)

type authController interface {
	Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error)
}

// Handler is the auth handler
type Handler struct {
	ctrl      authController
	validator *validator.Validate
	tracer    trace.Tracer
}

// New creates a new auth handler
func New(ctrl authController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		tracer:    tracer,
	}
}
