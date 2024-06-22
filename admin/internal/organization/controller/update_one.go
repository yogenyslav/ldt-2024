package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
)

// UpdateOne обновить название организации.
func (ctrl *Controller) UpdateOne(ctx context.Context, params model.OrganizationUpdateReq, username string) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.UpdateOne")
	defer span.End()

	return ctrl.repo.UpdateOne(ctx, model.OrganizationDao{
		Username: username,
		Title:    params.Title,
		ID:       params.ID,
	})
}
