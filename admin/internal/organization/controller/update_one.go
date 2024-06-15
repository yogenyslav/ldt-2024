package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
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

	if err := ctrl.repo.UpdateOne(ctx, model.OrganizationDao{
		ID:       params.ID,
		Username: username,
		Title:    params.Title,
	}); err != nil {
		log.Error().Err(err).Msg("failed to update organization")
		return shared.ErrUpdateOrganization
	}
	return nil
}
