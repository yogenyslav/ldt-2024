package controller

import (
	"github.com/yogenyslav/ldt-2024/admin/internal/api/pb"
	"go.opentelemetry.io/otel/trace"
)

// Controller имплементирует методы для работы с сервисом авторизации.
type Controller struct {
	authService pb.AuthServiceClient
	cipherKey   string
	tracer      trace.Tracer
}

// New создает новый Controller.
func New(authConn pb.AuthServiceClient, cipher string, tracer trace.Tracer) *Controller {
	return &Controller{
		authService: authConn,
		cipherKey:   cipher,
		tracer:      tracer,
	}
}
