package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
)

// SessionDao is a session model for the data access layer.
type SessionDao struct {
	CreatedAt time.Time `db:"created_at"`
	Username  string    `db:"username"`
	Title     string    `db:"title"`
	IsDeleted bool      `db:"is_deleted"`
	Tg        bool      `db:"tg"`
	ID        uuid.UUID `db:"id"`
}

// ToDto converts a SessionDao to a SessionDto.
func (s SessionDao) ToDto() SessionDto {
	return SessionDto{
		CreatedAt: s.CreatedAt,
		Title:     s.Title,
		Tg:        s.Tg,
		ID:        s.ID,
	}
}

// SessionContentDao is a model that holds all data layer queries and responses for related session.
type SessionContentDao struct {
	Response model.ResponseDao `db:"response"`
	Query    model.QueryDao    `db:"query"`
}

// SessionStatus is a model for session status info.
type SessionStatus struct {
	Username  string `db:"username"`
	IsDeleted bool   `db:"is_deleted"`
	Tg        bool   `db:"tg"`
}
