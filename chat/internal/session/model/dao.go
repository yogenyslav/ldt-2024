package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
)

// SessionDao модель для хранения сессии в базе данных.
type SessionDao struct {
	CreatedAt time.Time `db:"created_at"`
	Username  string    `db:"username"`
	Title     string    `db:"title"`
	IsDeleted bool      `db:"is_deleted"`
	Tg        bool      `db:"tg"`
	TgID      int64     `db:"tg_id"`
	ID        uuid.UUID `db:"id"`
}

// ToDto конвертирует SessionDao в SessionDto.
func (s SessionDao) ToDto() SessionDto {
	return SessionDto{
		CreatedAt: s.CreatedAt,
		Title:     s.Title,
		TgID:      s.TgID,
		Tg:        s.Tg,
		ID:        s.ID,
	}
}

// SessionContentDao модель сообщений и запросов в сессии.
type SessionContentDao struct {
	Response model.ResponseDao `db:"response"`
	Query    model.QueryDao    `db:"query"`
}

// SessionMeta метаинформация о сессии.
type SessionMeta struct {
	Username  string `db:"username"`
	Title     string `db:"title"`
	IsDeleted bool   `db:"is_deleted"`
	Tg        bool   `db:"tg"`
}
