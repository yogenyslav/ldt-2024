package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/internal/shared"
	"github.com/yogenyslav/ldt-2024/api/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PrepareData хендлер для подготовки данных.
func (h *Handler) PrepareData(c context.Context, in *pb.PrepareDataReq) (*emptypb.Empty, error) {
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
		"Handler.PrepareData",
	)
	defer span.End()

	organization, ok := ctx.Value(shared.OrganizationKey).(string)
	if !ok {
		log.Error().Msg("failed to get organization")
		return nil, status.Error(codes.Internal, "failed to get organization")
	}
	in.Organization = organization

	resp, err := h.predictor.PrepareData(ctx, in)
	if err != nil {
		log.Error().Err(err).Msg("failed to prepare data")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, status.Error(codes.OK, "data prepared")
}
