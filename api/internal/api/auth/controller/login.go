package controller

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
)

func (ctrl *Controller) Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.Login")
	defer span.End()

	var resp model.LoginResp

	token, err := ctrl.kc.Login(ctx, ctrl.cfg.KeyCloakClientID, ctrl.cfg.KeyCloakClientSecret, ctrl.cfg.KeyCloakRealm, params.Email, params.Password)
	if err != nil {
		return resp, errors.Wrap(err, "keycloak login failed")
	}

	user, err := ctrl.kc.GetUserInfo(ctx, token.AccessToken, ctrl.cfg.KeyCloakRealm)
	if err != nil {
		return resp, errors.Wrap(err, "failed to get user info")
	}

	log.Debug().Str("user", user.String()).Msg("user info")
	resp.Token = token.AccessToken

	return resp, nil
}
