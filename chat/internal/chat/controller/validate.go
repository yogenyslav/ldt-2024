package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateStatus updates query status by id.
func (ctrl *Controller) UpdateStatus(ctx context.Context, queryID int64, status shared.QueryStatus) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdateStatus",
		trace.WithAttributes(
			attribute.Int64("queryID", queryID),
			attribute.Int("status", int(status)),
		),
	)
	defer span.End()

	return ctrl.repo.UpdateQuery(ctx, model.QueryDao{
		ID:     queryID,
		Status: status,
	})
}
