package handler

import (
	"github.com/go-playground/validator/v10"
)

type sessionController interface {
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
