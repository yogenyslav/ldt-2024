package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"go.opentelemetry.io/otel/trace"
)

type organizationController interface {
	InsertOne(ctx context.Context, params model.OrganizationCreateReq, username string) (model.OrganizationCreateResp, error)
	FindOne(ctx context.Context, username string) (model.OrganizationDto, error)
	UpdateOne(ctx context.Context, params model.OrganizationUpdateReq, username string) error
}

// Handler имплементация сервера для организаций.
type Handler struct {
	ctrl   organizationController
	tracer trace.Tracer
}

// New создает новый Handler.
func New(ctrl organizationController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		tracer: tracer,
	}
}
