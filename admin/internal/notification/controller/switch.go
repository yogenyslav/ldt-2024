package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/notification/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Switch включает или выключает уведомления.
func (ctrl *Controller) Switch(ctx context.Context, params model.NotificationUpdateReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Switch",
		trace.WithAttributes(attribute.String("username", params.Username)),
	)
	defer span.End()

	token, err := ctrl.kc.LoginAdmin(ctx, ctrl.cfg.User, ctrl.cfg.Password, ctrl.cfg.AdminRealm)
	if err != nil {
		log.Error().Err(err).Msg("failed to login admin")
		return err
	}

	userInfo, err := ctrl.kc.GetUserInfo(ctx, token.AccessToken, ctrl.cfg.Realm)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		return err
	}

	if params.Active {
		exists, err := ctrl.repo.CheckNotification(ctx, model.NotificationDao{
			Email:          *userInfo.Email,
			OrganizationID: params.OrganizationID,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to check if notification exists")
			return err
		}
		if exists {
			return nil
		}

		return ctrl.repo.InsertOne(ctx, model.NotificationDao{
			Email:          *userInfo.Email,
			OrganizationID: params.OrganizationID,
		})
	}
	return ctrl.repo.DeleteOne(ctx, *userInfo.Email)
}
