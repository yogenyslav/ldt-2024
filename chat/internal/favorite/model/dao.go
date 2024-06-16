package model

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

// FavoriteDao представление избранных предиктов в базе данных.
type FavoriteDao struct {
	ID        int64
	Username  string
	Response  []byte
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDto конвертирует FavoriteDao в FavoriteDto.
func (f *FavoriteDao) ToDto() FavoriteDto {
	resp := make(map[string]any)
	if err := json.Unmarshal(f.Response, &resp); err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return FavoriteDto{
			ID:        f.ID,
			Response:  nil,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		}
	}
	return FavoriteDto{
		ID:        f.ID,
		Response:  resp,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
