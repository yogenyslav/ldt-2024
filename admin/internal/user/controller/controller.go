package controller

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/admin/config"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/trace"
)

type userOrgRepo interface {
	InsertOrganization(ctx context.Context, params model.UserOrganizationDao) error
	List(ctx context.Context, organizationID int64) ([]model.UserListResp, error)
	DeleteOrganization(ctx context.Context, params model.UserOrganizationDao) error
	CheckUserOrganization(ctx context.Context, username string, organizationID int64) (bool, error)
}

// Controller имплементирует методы для работы с пользователями и организациями.
type Controller struct {
	repo   userOrgRepo
	cfg    *config.KeyCloakConfig
	kc     *gocloak.GoCloak
	tracer trace.Tracer
}

// New создает новый Controller.
func New(repo userOrgRepo, cfg *config.KeyCloakConfig, kc *gocloak.GoCloak, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:   repo,
		cfg:    cfg,
		kc:     kc,
		tracer: tracer,
	}
}
