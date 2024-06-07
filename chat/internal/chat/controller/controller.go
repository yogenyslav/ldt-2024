package controller

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"go.opentelemetry.io/otel/trace"
)

type chatRepo interface {
	BeingTx(ctx context.Context) (context.Context, error)
	InsertQuery(ctx context.Context, params model.QueryDao) (int64, error)
	InsertResponse(ctx context.Context, params model.ResponseDao) error
}

// Controller is a struct that implements chat business logic.
type Controller struct {
	repo      chatRepo
	tracer    trace.Tracer
	kc        *gocloak.GoCloak
	realm     string
	cipherKey string
}

// New creates new Controller.
func New(repo chatRepo, kc *gocloak.GoCloak, realm, cipher string, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:      repo,
		kc:        kc,
		realm:     realm,
		cipherKey: cipher,
		tracer:    tracer,
	}
}
