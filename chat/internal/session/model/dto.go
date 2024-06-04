package model

import (
	"github.com/google/uuid"
)

type Session struct {
	ID uuid.UUID `json:"id"`
}
