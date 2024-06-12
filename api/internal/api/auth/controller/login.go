package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
)

// Login выполняет вход пользователя в систему.
func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Login")
	defer span.End()

	var resp model.LoginResp

	token, err := ctrl.kc.Login(ctx, ctrl.cfg.ClientID, ctrl.cfg.ClientSecret, ctrl.cfg.Realm, params.Username, params.Password)
	if err != nil {
		return resp, err
	}

	info, err := ctrl.kc.RetrospectToken(ctx, token.AccessToken, ctrl.cfg.ClientID, ctrl.cfg.ClientSecret, ctrl.cfg.Realm)
	if err != nil {
		return resp, err
	}

	log.Debug().Any("permissions", info.Permissions).Msg("token info")
	resp.Token = token.AccessToken
	resp.Role = "admin"
	return resp, nil
}
