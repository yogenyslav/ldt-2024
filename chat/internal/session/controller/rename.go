package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Rename переименовывает сессию.
func (ctrl *Controller) Rename(ctx context.Context, params model.RenameReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.Rename",
		trace.WithAttributes(
			attribute.String("id", params.ID.String()),
			attribute.String("title", params.Title),
		),
	)
	defer span.End()

	if err := ctrl.repo.UpdateTitle(ctx, params); err != nil {
		return err
	}
	return nil
}
