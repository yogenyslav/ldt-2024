package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
)

// SessionDto is a session model for the domain layer.
type SessionDto struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Tg        bool      `json:"tg"`
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

// FindOneResp is a model for find one session request.
type FindOneResp struct {
	Content  []SessionContentDto `json:"content"`
	Editable bool                `json:"editable"`
	Tg       bool                `json:"tg"`
	ID       uuid.UUID           `json:"id"`
}

// SessionContentDto is a model that holds all domain layer queries and responses for related session.
type SessionContentDto struct {
	Response model.ResponseDto `json:"response"`
	Query    model.QueryDto    `json:"query"`
}
