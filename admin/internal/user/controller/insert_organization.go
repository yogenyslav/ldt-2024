package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertOrganization добавляет организацию пользователю.
func (ctrl *Controller) InsertOrganization(ctx context.Context, params model.UserUpdateOrganizationReq) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertOrganization",
		trace.WithAttributes(
			attribute.Int64("organizationID", params.OrganizationID),
			attribute.String("username", params.Username),
		),
	)
	defer span.End()

	exists, err := ctrl.repo.CheckUserOrganization(ctx, params.Username, params.OrganizationID)
	if err != nil {
		log.Error().Err(err).Msg("failed to check user in organization")
		return err
	}
	if exists {
		log.Error().Err(err).Msg("user is already in organization")
		return nil
	}

	return ctrl.repo.InsertOrganization(ctx, model.UserOrganizationDao{
		Username:       params.Username,
		OrganizationID: params.OrganizationID,
	})
}
