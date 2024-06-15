package controller

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteOrganization удаляет организацию пользователя.
func (ctrl *Controller) DeleteOrganization(ctx context.Context, username string) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteOrganization",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()

	return ctrl.repo.DeleteOrganization(ctx, username)
}
