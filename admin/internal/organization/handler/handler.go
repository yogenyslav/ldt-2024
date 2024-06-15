package handler

import (
	"context"
	"mime/multipart"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"go.opentelemetry.io/otel/trace"
)

type organizationController interface {
	InsertOne(ctx context.Context, params model.OrganizationCreateReq, username string) error
	FindOne(ctx context.Context, username string) (model.OrganizationDto, error)
	UpdateOne(ctx context.Context, params model.OrganizationUpdateReq, username string) error
	ImportData(ctx context.Context, mpArchive *multipart.FileHeader, org string) error
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
