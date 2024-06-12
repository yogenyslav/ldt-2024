package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/internal/api/prompter/model"
	"go.opentelemetry.io/otel/trace"
)

type prompterController interface {
	Extract(ctx context.Context, params model.ExtractReq) (*pb.ExtractedPrompt, error)
}

// Handler имплементация сервиса Prompter.
type Handler struct {
	pb.UnimplementedPrompterServer
	ctrl   prompterController
	tracer trace.Tracer
}

// New создает новый Handler.
func New(ctrl prompterController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		tracer: tracer,
	}
}
