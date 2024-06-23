package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/notification/model"
	"go.opentelemetry.io/otel/trace"
)

type notificationRepo interface {
	InsertOne(ctx context.Context, params model.NotificationDao) error
	DeleteOne(ctx context.Context, email string) error
	CheckNotification(ctx context.Context, email string, organizationID int64) (bool, error)
}

// Controller контроллер для уведомлений.
type Controller struct {
	repo   notificationRepo
	tracer trace.Tracer
}

// New создает новый Controller.
func New(repo notificationRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		tracer: tracer,
	}
}
