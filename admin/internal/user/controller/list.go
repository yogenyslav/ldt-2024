package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// List возвращает список пользователей в организации.
func (ctrl *Controller) List(ctx context.Context, organizationID int64) ([]model.UserListResp, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.List",
		trace.WithAttributes(attribute.Int64("organizationID", organizationID)),
	)
	defer span.End()
	return ctrl.repo.List(ctx, organizationID)
}
