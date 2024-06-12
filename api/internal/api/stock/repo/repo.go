package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/api/internal/api/stock/model"
	"github.com/yogenyslav/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
)

const coll = "product"

// Repo репозиторий для работы с продуктами.
type Repo struct {
	mongo storage.MongoDatabase
}

// New создает новый репозиторий.
func New(mongo storage.MongoDatabase) *Repo {
	return &Repo{
		mongo: mongo,
	}
}

// ListProducts найти все продукты.
func (r *Repo) ListProducts(ctx context.Context) ([]model.ProductDao, error) {
	var products []model.ProductDao
	err := r.mongo.FindMany(ctx, coll, bson.D{}, &products)
	return products, err
}
