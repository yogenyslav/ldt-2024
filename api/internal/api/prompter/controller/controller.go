package controller

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"go.opentelemetry.io/otel/trace"
)

// Controller is a struct that implements prompter business logic.
type Controller struct {
	prompter pb.PrompterClient
	tracer   trace.Tracer
}

// New creates new Controller.
func New(prompter client.GrpcClient, tracer trace.Tracer) *Controller {
	return &Controller{
		prompter: pb.NewPrompterClient(prompter.GetConn()),
		tracer:   tracer,
	}
}
