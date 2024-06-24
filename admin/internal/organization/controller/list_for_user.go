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

// ListForUser находит список организаций по username.
func (ctrl *Controller) ListForUser(ctx context.Context, username string) ([]model.OrganizationDto, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListForUser",
		trace.WithAttributes(attribute.String("username", username)),
	)
	defer span.End()

	organizationDB, err := ctrl.repo.ListForUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.ErrNoOrganization
		}
		log.Error().Err(err).Msg("failed to get organization")
		return nil, shared.ErrGetOrganization
	}

	organization := make([]model.OrganizationDto, len(organizationDB))
	for i := 0; i < len(organizationDB); i++ {
		organization[i] = organizationDB[i].ToDto()
	}

	return organization, nil
}
