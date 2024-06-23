package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/chat/internal/notification/model"
)

type notificationController interface {
	Switch(ctx context.Context, params model.NotificationUpdateReq) error
}

// Handler обработчик для уведомлений.
type Handler struct {
	ctrl      notificationController
	validator *validator.Validate
}

// New создает новый Handler.
func New(ctrl notificationController) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
