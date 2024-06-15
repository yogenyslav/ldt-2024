package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertOrganization добавляет организацию пользователю.
func (ctrl *Controller) InsertOrganization(ctx context.Context, params model.UserUpdateOrganizationReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertOrganization",
		trace.WithAttributes(
			attribute.String("organization", params.Organization),
			attribute.String("username", params.Username),
		),
	)
	defer span.End()

	return ctrl.repo.InsertOrganization(ctx, model.UserOrganizationDao{
		Username:     params.Username,
		Organization: params.Organization,
	})
}
