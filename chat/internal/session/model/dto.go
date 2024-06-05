package model

import (
	"github.com/google/uuid"
)

// SessionDto is a session model for the domain layer.
type SessionDto struct {
	ID uuid.UUID `json:"id"`
}
