package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/stock/model"
	"github.com/yogenyslav/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
)

const coll = "product"

// Repo is a stock repository.
type Repo struct {
	mongo storage.MongoDatabase
}

// New creates new Repo.
func New(mongo storage.MongoDatabase) *Repo {
	return &Repo{
		mongo: mongo,
	}
}

// ListProducts finds all products.
func (r *Repo) ListProducts(ctx context.Context) ([]model.ProductDao, error) {
	var products []model.ProductDao
	err := r.mongo.FindMany(ctx, coll, bson.D{}, &products)
	return products, err
}
