package model

import (
	"time"

	"github.com/google/uuid"
)

// SessionDao is a session model for the data access layer.
type SessionDao struct {
	CreatedAt time.Time `db:"created_at"`
	Username  string    `db:"username"`
	Title     string    `db:"title"`
	IsDeleted bool      `db:"is_deleted"`
	ID        uuid.UUID `db:"id"`
}

// ToDto converts a SessionDao to a SessionDto.
func (s SessionDao) ToDto() SessionDto {
	return SessionDto{
		CreatedAt: s.CreatedAt,
		Title:     s.Title,
		ID:        s.ID,
	}
}
