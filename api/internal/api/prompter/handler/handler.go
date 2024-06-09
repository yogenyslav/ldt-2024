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

// Handler is a struct that implements prompter server.
type Handler struct {
	pb.UnimplementedPrompterServer
	ctrl   prompterController
	tracer trace.Tracer
}

// New creates new Handler.
func New(ctrl prompterController, tracer trace.Tracer) *Handler {
	return &Handler{
		ctrl:   ctrl,
		tracer: tracer,
	}
}
