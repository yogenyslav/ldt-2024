package handler

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthFuncOverride overrides auth interceptor for auth service.
func (h *Handler) AuthFuncOverride(ctx context.Context, fullMethod string) (context.Context, error) {
	log.Debug().Str("fullMethod", fullMethod).Msg("auth")
	if fullMethod == shared.LoginEndpoint {
		return ctx, nil
	}
	return nil, status.Error(codes.Unauthenticated, "wrong endpoint for login")
}
