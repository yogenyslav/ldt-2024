package controller

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
)

// ListProducts получает список всех продуктов.
func (ctrl *Controller) ListProducts(ctx context.Context) ([]*pb.Product, error) {
	ctx, span := ctrl.tracer.Start(ctx, "Controller.ListProducts")
	defer span.End()

	productsDB, err := ctrl.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	products := make([]*pb.Product, len(productsDB))

	for i := 0; i < len(productsDB); i++ {
		products[i] = productsDB[i].ToPb()
	}
	return products, nil
}
