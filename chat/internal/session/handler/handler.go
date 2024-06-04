package handler

import (
	"github.com/go-playground/validator/v10"
)

type sessionController interface {
}

type Handler struct {
	ctrl      sessionController
	validator *validator.Validate
}

func New(ctrl sessionController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
