package model

import (
	"time"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// QueryDto модель запроса.
type QueryDto struct {
	CreatedAt time.Time `json:"created_at"`
	Prompt    string    `json:"prompt"`
	Period    string    `json:"period"`
	Product   string    `json:"product"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	ID        int64     `json:"id"`
}

// ResponseDto модель ответа.
type ResponseDto struct {
	CreatedAt time.Time      `json:"created_at"`
	Body      string         `json:"body"`
	Data      map[string]any `json:"data"`
	DataType  string         `json:"data_type"`
	Status    string         `json:"status"`
}

// QueryCreateReq модель запроса для создания.
type QueryCreateReq struct {
	Prompt  string              `json:"prompt,omitempty" validate:"gte=5"`
	Command shared.QueryCommand `json:"command,omitempty"`
}

// QueryMeta модель метаданных запроса.
type QueryMeta struct {
	Product string
	Period  string
	Type    shared.QueryType
}
