package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
)

// NewUser создает нового пользователя.
func (ctrl *Controller) NewUser(ctx context.Context, params model.UserCreateReq) error {
	ctx, span := ctrl.tracer.Start(ctx, "controller.NewUser")
	defer span.End()

	token, err := ctrl.kc.LoginAdmin(ctx, ctrl.cfg.User, ctrl.cfg.Password, ctrl.cfg.AdminRealm)
	if err != nil {
		log.Error().Err(err).Msg("failed to login admin")
		return err
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(params.FirstName),
		LastName:  gocloak.StringP(params.LastName),
		Email:     gocloak.StringP(params.Email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(params.Username),
	}
	userID, err := ctrl.kc.CreateUser(ctx, token.AccessToken, ctrl.cfg.Realm, user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user in kc")
		return err
	}

	if err := ctrl.kc.SetPassword(ctx, token.AccessToken, userID, ctrl.cfg.Realm, params.Password, false); err != nil {
		log.Error().Err(err).Msg("failed to set password")
		return err
	}

	for _, role := range params.Roles {
		path := fmt.Sprintf("/%s", strings.ToLower(role))
		group, err := ctrl.kc.GetGroupByPath(ctx, token.AccessToken, ctrl.cfg.Realm, path)
		if err != nil {
			log.Error().Err(err).Msg("failed to get group by path")
			return err
		}
		if err := ctrl.kc.AddUserToGroup(ctx, token.AccessToken, ctrl.cfg.Realm, userID, *group.ID); err != nil {
			log.Error().Err(err).Msg("failed to add user to group")
			return err
		}
	}

	if params.Organization != "" {
		if err := ctrl.repo.InsertOrganization(ctx, model.UserOrganizationDao{
			Username:     params.Username,
			Organization: params.Organization,
		}); err != nil {
			log.Error().Err(err).Msg("failed to insert organization")
			return err
		}
	}

	return nil
}
