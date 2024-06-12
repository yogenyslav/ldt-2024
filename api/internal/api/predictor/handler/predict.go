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

// Predict хендлер для предсказания.
func (h *Handler) Predict(c context.Context, in *pb.PredictReq) (*pb.PredictResp, error) {
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
		"Handler.Predict",
		trace.WithAttributes(
			attribute.String("product", in.GetProduct()),
			attribute.String("period", in.GetPeriod()),
			attribute.Int("type", int(in.GetType())),
		),
	)
	defer span.End()

	resp, err := h.predictor.Predict(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to predict")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, status.Error(codes.OK, "predicted")
}
