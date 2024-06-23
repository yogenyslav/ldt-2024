package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/user/model"
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
		return ctrl.repo.InsertOne(ctx, model.NotificationDao{
			Email:          params.Email,
			FirstName:      params.FirstName,
			LastName:       params.LastName,
			OrganizationID: params.OrganizationID,
		})
	}
	return ctrl.repo.DeleteOne(ctx, params.Email)
}
