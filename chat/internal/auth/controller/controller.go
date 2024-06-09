package controller

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"go.opentelemetry.io/otel/trace"
)

// Controller is a struct that implements the business logic of the auth service.
type Controller struct {
	authService pb.AuthServiceClient
	tracer      trace.Tracer
	cipherKey   string
}

// New creates a new instance of the Controller.
func New(authConn pb.AuthServiceClient, cipher string, tracer trace.Tracer) *Controller {
	return &Controller{
		authService: authConn,
		cipherKey:   cipher,
		tracer:      tracer,
	}
}
