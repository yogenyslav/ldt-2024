package controller

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// FindOne находит организацию по username.
func (ctrl *Controller) FindOne(ctx context.Context, username string) (model.OrganizationDto, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.FindOne",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()

	org, err := ctrl.repo.FindOne(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrganizationDto{}, shared.ErrNoOrganization
		}
		log.Error().Err(err).Msg("failed to get organization")
		return model.OrganizationDto{}, shared.ErrGetOrganization
	}

	return org.ToDto(), nil
}
