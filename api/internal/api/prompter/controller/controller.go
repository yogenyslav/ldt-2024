package controller

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"go.opentelemetry.io/otel/trace"
)

// Controller имплементирует методы для работы с Prompter.
type Controller struct {
	prompter pb.PrompterClient
	tracer   trace.Tracer
}

// New создает новый Controller.
func New(prompter client.GrpcClient, tracer trace.Tracer) *Controller {
	return &Controller{
		prompter: pb.NewPrompterClient(prompter.GetConn()),
		tracer:   tracer,
	}
}
