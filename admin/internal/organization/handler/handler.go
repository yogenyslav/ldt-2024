package handler

import (
	"context"
	"mime/multipart"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
)

type organizationController interface {
	InsertOne(ctx context.Context, params model.OrganizationCreateReq, username string) (model.OrganizationCreateResp, error)
	ListForUser(ctx context.Context, username string) ([]model.OrganizationDto, error)
	ImportData(ctx context.Context, mpArchive *multipart.FileHeader, id int64) error
	UpdateOne(ctx context.Context, params model.OrganizationUpdateReq, username string) error
}

// Handler имплементация сервера для организаций.
type Handler struct {
	ctrl   organizationController
	m      *metrics.Metrics
	tracer trace.Tracer
}

// New создает новый Handler.
func New(ctrl organizationController, m *metrics.Metrics, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		m:      m,
		tracer: tracer,
	}
}
