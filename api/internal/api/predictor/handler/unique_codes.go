package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UniqueCodes хендлер для подготовки данных.
func (h *Handler) UniqueCodes(c context.Context, _ *emptypb.Empty) (*pb.UniqueCodesResp, error) {
	ctx, err := pkg.GetTraceCtx(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to get trace context")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if ctx == nil {
		ctx = c
	}

	ctx, span := h.tracer.Start(ctx, "Handler.UniqueCodes")
	defer span.End()

	resp, err := h.predictor.UniqueCodes(ctx, &emptypb.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("failed to get unique codes")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, status.Error(codes.OK, "unique codes fetched")
}
