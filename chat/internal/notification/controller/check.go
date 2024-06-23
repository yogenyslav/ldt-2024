package controller

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Check проверяет уведомления.
func (ctrl *Controller) Check(ctx context.Context, email string, organizationID int64) (bool, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Check",
		trace.WithAttributes(attribute.String("email", email)),
		trace.WithAttributes(attribute.Int64("organizationID", organizationID)),
	)
	defer span.End()

	return ctrl.repo.CheckNotification(ctx, email, organizationID)
}
