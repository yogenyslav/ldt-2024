package controller

import (
	"context"
	"errors"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"

	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
	"github.com/yogenyslav/ldt-2024/api/internal/shared"
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

	userInfo, err := ctrl.kc.GetUserInfo(ctx, token.AccessToken, ctrl.cfg.Realm)
	if err != nil {
		return resp, err
	}
	if userInfo.Sub == nil {
		return resp, errors.New("userID is nil")
	}

	userID := *userInfo.Sub

	groups, err := ctrl.kc.GetUserGroups(ctx, token.AccessToken, ctrl.cfg.Realm, userID, gocloak.GetGroupsParams{})
	if err != nil {
		return resp, err
	}

	roles := make([]pb.UserRole, len(groups))
	for idx, role := range groups {
		roles[idx] = shared.RoleFromString(*role.Name)
	}

	resp.Token = token.AccessToken
	resp.Roles = roles
	return resp, nil
}
