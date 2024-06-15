package controller

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// List возвращает список пользователей в организации.
func (ctrl *Controller) List(ctx context.Context, organization string) ([]string, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.List",
		trace.WithAttributes(attribute.String("organization", organization)),
	)
	defer span.End()
	return ctrl.repo.List(ctx, organization)
}
