package controller

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/api/config"
	"go.opentelemetry.io/otel/trace"
)

// Controller имплементирует методы для работы с KeyCloak.
type Controller struct {
	kc     *gocloak.GoCloak
	cfg    *config.KeyCloakConfig
	tracer trace.Tracer
}

// New создает новый Controller.
func New(kc *gocloak.GoCloak, cfg *config.KeyCloakConfig, tracer trace.Tracer) *Controller {
	return &Controller{
		kc:     kc,
		cfg:    cfg,
		tracer: tracer,
	}
}
