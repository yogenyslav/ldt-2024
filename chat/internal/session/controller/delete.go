package controller

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Delete удаляет сессию по id.
func (ctrl *Controller) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Delete",
		trace.WithAttributes(attribute.String("id", id.String())),
	)
	defer span.End()

	if err := ctrl.repo.DeleteOne(ctx, id); err != nil {
		return err
	}
	return nil
}
