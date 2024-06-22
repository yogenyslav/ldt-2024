package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/pkg"
	"github.com/yogenyslav/pkg/storage/minios3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InsertOne создает новую организацию.
func (ctrl *Controller) InsertOne(ctx context.Context, params model.OrganizationCreateReq, username string) (model.OrganizationCreateResp, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.InsertOrganization",
		trace.WithAttributes(
			attribute.String("title", params.Title),
			attribute.String("username", username),
		),
	)
	defer span.End()

	var resp model.OrganizationCreateResp

	id, err := ctrl.repo.InsertOne(ctx, model.OrganizationDao{
		Username: username,
		Title:    params.Title,
	})
	if err != nil {
		if pkg.CheckDuplicateKey(err) {
			return resp, shared.ErrDuplicateTitle
		}
		log.Error().Err(err).Msg("failed to insert organization")
		return resp, shared.ErrCreateOrganization
	}

	err = ctrl.s3.CreateBucket(ctx, &minios3.Bucket{
		Name:   getOrganizationTitle(id),
		Region: "eu-central-1",
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create bucket")
		return resp, shared.ErrCreateOrganization
	}

	resp.ID = id
	return resp, nil
}
