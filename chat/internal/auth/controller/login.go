package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"github.com/yogenyslav/ldt-2024/chat/pkg/secure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Login is a method that implements the login logic.
func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Login",
		trace.WithAttributes(attribute.String("username", params.Username)),
	)
	defer span.End()

	ctx = pkg.PushSpan(ctx, span)

	var resp model.LoginResp

	in := &pb.LoginRequest{
		Username: params.Username,
		Password: params.Password,
	}
	token, err := ctrl.authService.Login(ctx, in)
	if err != nil {
		return resp, shared.ErrLoginFailed
	}

	tokenEncrypted, err := secure.Encrypt(token.Token, ctrl.cipherKey)
	if err != nil {
		log.Error().Err(err).Msg("encryption internal error")
		return resp, shared.ErrEncryption
	}
	resp.Token = tokenEncrypted

	return resp, nil
}
