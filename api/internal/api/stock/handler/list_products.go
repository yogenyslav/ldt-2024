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

// ListProducts хендлер для получения списка всех продуктов.
func (h *Handler) ListProducts(c context.Context, _ *emptypb.Empty) (*pb.ListProductsResp, error) {
	ctx, err := pkg.GetTraceCtx(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to get trace context")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if ctx == nil {
		ctx = c
	}

	ctx, span := h.tracer.Start(ctx, "Handler.ListProducts")
	defer span.End()

	products, err := h.ctrl.ListProducts(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list products")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListProductsResp{
		Products: products,
	}, status.Error(codes.OK, "list products")
}
