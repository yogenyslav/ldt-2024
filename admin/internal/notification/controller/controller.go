package controller

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/admin/config"
	"github.com/yogenyslav/ldt-2024/admin/internal/notification/model"
	"go.opentelemetry.io/otel/trace"
)

type notificationRepo interface {
	InsertOne(ctx context.Context, params model.NotificationDao) error
	DeleteOne(ctx context.Context, email string) error
	CheckNotification(ctx context.Context, params model.NotificationDao) (bool, error)
}

// Controller контроллер для уведомлений.
type Controller struct {
	repo   notificationRepo
	cfg    *config.KeyCloakConfig
	kc     *gocloak.GoCloak
	tracer trace.Tracer
}

// New создает новый Controller.
func New(repo notificationRepo, cfg *config.KeyCloakConfig, kc *gocloak.GoCloak, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		cfg:    cfg,
		kc:     kc,
		tracer: tracer,
	}
}
