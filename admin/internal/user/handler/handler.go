package handler

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"github.com/yogenyslav/ldt-2024/admin/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
)

type userController interface {
	NewUser(ctx context.Context, req model.UserCreateReq) error
	InsertOrganization(ctx context.Context, req model.UserUpdateOrganizationReq) error
	DeleteOrganization(ctx context.Context, req model.UserUpdateOrganizationReq) error
	List(ctx context.Context, organizationID int64) ([]string, error)
}

// Handler имплементация сервера для работы с пользователями и организациями.
type Handler struct {
	ctrl      userController
	validator *validator.Validate
	m         *metrics.Metrics
	tracer    trace.Tracer
}

// New создает новый Handler.
func New(ctrl userController, m *metrics.Metrics, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:      ctrl,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		m:         m,
		tracer:    tracer,
	}
}
