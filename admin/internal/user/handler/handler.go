package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/trace"
)

type userController interface {
	NewUser(ctx context.Context, req model.UserCreateReq) error
	InsertOrganization(ctx context.Context, req model.UserUpdateOrganizationReq) error
	DeleteOrganization(ctx context.Context, username string) error
	List(ctx context.Context, organization string) ([]string, error)
}

// Handler имплементация сервера для работы с пользователями и организациями.
type Handler struct {
	ctrl      userController
	validator *validator.Validate
	tracer    trace.Tracer
}

// New создает новый Handler.
func New(ctrl userController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		tracer:    tracer,
	}
}
