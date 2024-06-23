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
		trace.WithAttributes(attribute.String("email", params.Email)),
	)
	defer span.End()

	if params.Active {
		exists, err := ctrl.repo.CheckNotification(ctx, params.Email, params.OrganizationID)
		if err != nil {
			log.Error().Err(err).Msg("failed to check if notification exists")
			return err
		}
		if exists {
			return nil
		}

		return ctrl.repo.InsertOne(ctx, model.NotificationDao{
			Email:          params.Email,
			OrganizationID: params.OrganizationID,
		})
	}
	return ctrl.repo.DeleteOne(ctx, params.Email)
}
