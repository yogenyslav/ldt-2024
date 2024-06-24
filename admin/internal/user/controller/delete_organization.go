package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteOrganization удаляет организацию пользователя.
func (ctrl *Controller) DeleteOrganization(ctx context.Context, params model.UserUpdateOrganizationReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteOrganization",
		trace.WithAttributes(attribute.String("username", params.Username)),
	)
	defer span.End()

	return ctrl.repo.DeleteOrganization(ctx, model.UserOrganizationDao{
		Username:       params.Username,
		OrganizationID: params.OrganizationID,
	})
}
