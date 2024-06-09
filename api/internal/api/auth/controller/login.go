package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
)

// Login logs in a user.
func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Login")
	defer span.End()

	var resp model.LoginResp

	token, err := ctrl.kc.Login(ctx, ctrl.cfg.ClientID, ctrl.cfg.ClientSecret, ctrl.cfg.Realm, params.Username, params.Password)
	if err != nil {
		return resp, err
	}

	resp.Token = token.AccessToken
	return resp, nil
}
