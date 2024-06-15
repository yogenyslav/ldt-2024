package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateOne обновляет организацию.
func (ctrl *Controller) UpdateOne(ctx context.Context, params model.OrganizationUpdateReq, username string) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdateOne",
		trace.WithAttributes(
			attribute.String("title", params.Title),
			attribute.String("username", username),
		),
	)
	defer span.End()

	return ctrl.repo.UpdateOne(ctx, model.OrganizationDao{
		ID:       params.ID,
		Username: username,
		Title:    params.Title,
	})
}
