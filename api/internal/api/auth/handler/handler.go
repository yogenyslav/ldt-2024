package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
)

type authController interface {
	Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error)
}

// Handler имплементация сервиса авторизации.
type Handler struct {
	pb.UnimplementedAuthServiceServer
	ctrl    authController
	tracer  trace.Tracer
	metrics *metrics.Metrics
}

// New создает новый Handler.
func New(ctrl authController, tracer trace.Tracer, m *metrics.Metrics) *Handler {
	return &Handler{
		ctrl:    ctrl,
		tracer:  tracer,
		metrics: m,
	}
}
