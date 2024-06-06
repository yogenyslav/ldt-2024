package model

import (
	"time"

	"github.com/google/uuid"
)

// SessionDto is a session model for the domain layer.
type SessionDto struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	ID        uuid.UUID `json:"id"`
}

// NewSessionResp is a session create response model.
type NewSessionResp struct {
	ID uuid.UUID `json:"id"`
}

// ListResp is a struct for list sessions response.
type ListResp struct {
	Sessions []SessionDto `json:"sessions"`
}
