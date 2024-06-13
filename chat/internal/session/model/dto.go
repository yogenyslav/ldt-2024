package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
)

// SessionDto модель для передачи сессии внутри сервиса.
type SessionDto struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Tg        bool      `json:"tg"`
	TgID      int64     `json:"tg_id"`
	ID        uuid.UUID `json:"id"`
}

// NewSessionResp модель для ответа на запрос создания сессии.
type NewSessionResp struct {
	ID uuid.UUID `json:"id"`
}

// ListResp модель для ответа на запрос списка сессий.
type ListResp struct {
	Sessions []SessionDto `json:"sessions"`
}

// RenameReq модель для запроса на переименование сессии.
type RenameReq struct {
	Title string    `json:"title" validate:"required"`
	ID    uuid.UUID `json:"id" validate:"required"`
}

// FindOneResp модель для ответа на запрос получения сессии по id.
type FindOneResp struct {
	Title    string              `json:"title"`
	Content  []SessionContentDto `json:"content"`
	Editable bool                `json:"editable"`
	Tg       bool                `json:"tg"`
	ID       uuid.UUID           `json:"id"`
}

// SessionContentDto модель для передачи контента сессии внутри сервиса.
type SessionContentDto struct {
	Response model.ResponseDto `json:"response"`
	Query    model.QueryDto    `json:"query"`
}
