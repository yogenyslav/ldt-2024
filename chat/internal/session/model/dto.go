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

// RenameReq is a model for session rename request.
type RenameReq struct {
	Title string    `json:"title" validate:"required,gte=3,max=25"`
	ID    uuid.UUID `json:"id" validate:"required"`
}
