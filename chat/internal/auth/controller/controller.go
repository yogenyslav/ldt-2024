package controller

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Controller struct {
	authService pb.AuthServiceClient
	tracer      trace.Tracer
}

func New(authConn *grpc.ClientConn, tracer trace.Tracer) *Controller {
	return &Controller{
		authService: pb.NewAuthServiceClient(authConn),
		tracer:      tracer,
	}
}
