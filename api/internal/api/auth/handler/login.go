package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login logs in a user.
func (h *Handler) Login(c context.Context, params *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, err := pkg.GetTraceCtx(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to get trace context")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if ctx == nil {
		ctx = c
	}

	ctx, span := h.tracer.Start(ctx, "Handler.Login")
	defer span.End()

	req := model.LoginReq{
		Username: params.GetUsername(),
		Password: params.GetPassword(),
	}
	resp, err := h.ctrl.Login(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to login")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Info().Str("username", req.Username).Msg("user logged in")
	h.metrics.LoginCount.Inc()

	return &pb.LoginResponse{
		Token: resp.Token,
	}, status.Error(codes.OK, "login success")
}
