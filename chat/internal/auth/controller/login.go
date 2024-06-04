package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
)

func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Login")
	defer span.End()

	var resp model.LoginResp

	in := &pb.LoginRequest{
		Email:    params.Email,
		Password: params.Password,
	}
	token, err := ctrl.authService.Login(ctx, in)
	if err != nil {
		return resp, err
	}

	resp.Token = token.Token
	return resp, nil
}
