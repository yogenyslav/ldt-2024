package handler

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"go.opentelemetry.io/otel/trace"
)

// Handler имплементация сервиса Prompter.
type Handler struct {
	pb.UnimplementedPrompterServer
	prompter pb.PrompterClient
	tracer   trace.Tracer
}

// New создает новый Handler.
func New(prompter client.GrpcClient, tracer trace.Tracer) *Handler {
	return &Handler{
		prompter: pb.NewPrompterClient(prompter.GetConn()),
		tracer:   tracer,
	}
}
