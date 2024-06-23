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

	email, err := ctrl.repo.GetEmailByUsername(ctx, params.Username)
	if err != nil {
		log.Error().Err(err).Msg("failed to get email by username")
		return err
	}

	if params.Active {
		exists, err := ctrl.repo.CheckNotification(ctx, model.NotificationDao{
			Email:          email,
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
			Email:          email,
			OrganizationID: params.OrganizationID,
		})
	}
	return ctrl.repo.DeleteOne(ctx, email)
}
