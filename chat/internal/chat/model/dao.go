package model

import (
	"time"

	"github.com/google/uuid"
)

// QueryDao is a data layer representation of users query.
type QueryDao struct {
	CreatedAt time.Time `db:"created_at"`
	Prompt    string    `db:"prompt"`
	Command   string    `db:"command"`
	Username  string    `db:"username"`
	ID        int64     `db:"id"`
	SessionID uuid.UUID `db:"session_id"`
}
