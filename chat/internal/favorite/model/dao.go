package model

import (
	"encoding/json"
	"time"
)

// FavoriteDao представление избранных предиктов в базе данных.
type FavoriteDao struct {
	QueryID   int64
	Username  string
	Response  []byte
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDto конвертирует FavoriteDao в FavoriteDto.
func (f *FavoriteDao) ToDto() FavoriteDto {
	resp := make(map[string]any)
	_ = json.Unmarshal(f.Response, &resp)
	return FavoriteDto{
		QueryID:   f.QueryID,
		Response:  resp,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
