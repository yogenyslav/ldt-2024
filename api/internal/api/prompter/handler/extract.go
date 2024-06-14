package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Extract хендлер для извлечения метаданных из промпта.
func (h *Handler) Extract(c context.Context, in *pb.ExtractReq) (*pb.ExtractedPrompt, error) {
	ctx, err := pkg.GetTraceCtx(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to get trace context")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if ctx == nil {
		ctx = c
	}

	ctx, span := h.tracer.Start(
		ctx,
		"Handler.Extract",
		trace.WithAttributes(attribute.String("prompt", in.GetPrompt())),
	)
	defer span.End()

	resp, err := h.prompter.Extract(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to extract prompt metadata")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, status.Error(codes.OK, "metadata extracted")
}
