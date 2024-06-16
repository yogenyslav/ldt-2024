package controller

import (
	"context"
)

// DeleteOne удаляет предикт из избранного.
func (ctrl *Controller) DeleteOne(ctx context.Context, queryID int64, username string) error {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.DeleteOne")
	defer span.End()

	if err := ctrl.repo.DeleteOne(ctx, queryID, username); err != nil {
		return err
	}

	return nil
}
