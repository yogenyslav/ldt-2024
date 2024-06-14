package model

import (
	"time"
)

// FavoriteDto представление избранных предиктов в API.
type FavoriteDto struct {
	QueryID   int64          `json:"id"`
	Response  map[string]any `json:"response"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// FavoriteCreateReq запрос на добавление предикта в избранное.
type FavoriteCreateReq struct {
	QueryID  int64          `json:"id"`
	Response map[string]any `json:"response"`
}

// FavoriteUpdateReq запрос на обновление предикта в избранном.
type FavoriteUpdateReq struct {
	QueryID  int64          `json:"id"`
	Response map[string]any `json:"response"`
}
