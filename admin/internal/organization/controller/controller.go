package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/pkg/storage/minios3"
	"go.opentelemetry.io/otel/trace"
)

type organizationRepo interface {
	InsertOne(ctx context.Context, params model.OrganizationDao) (int64, error)
	FindOne(ctx context.Context, username string) (model.OrganizationDao, error)
	UpdateOne(ctx context.Context, params model.OrganizationDao) error
}

// Controller имплементирует методы для работы с организациями.
type Controller struct {
	repo   organizationRepo
	s3     minios3.S3
	tracer trace.Tracer
}

// New создает новый Controller.
func New(repo organizationRepo, s3 minios3.S3, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		s3:     s3,
		tracer: tracer,
	}
}
