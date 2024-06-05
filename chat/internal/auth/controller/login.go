package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Login is a method that implements the login logic.
func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Login",
		trace.WithAttributes(attribute.String("email", params.Email)),
	)
	defer span.End()

	ctx = pkg.PushSpan(ctx, span)

	var resp model.LoginResp

	in := &pb.LoginRequest{
		Email:    params.Email,
		Password: params.Password,
	}
	token, err := ctrl.authService.Login(ctx, in)
	if err != nil {
		return resp, shared.ErrLoginFailed
	}

	resp.Token = token.Token
	return resp, nil
}
