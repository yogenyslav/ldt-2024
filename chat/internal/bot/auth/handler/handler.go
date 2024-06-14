package handler

import (
	"go.opentelemetry.io/otel/trace"
)

// Handler is the handler for the auth command.
type Handler struct {
	tracer trace.Tracer
}

// New creates a new Handler.
func New(tracer trace.Tracer) *Handler {
	return &Handler{
		tracer: tracer,
	}
}
