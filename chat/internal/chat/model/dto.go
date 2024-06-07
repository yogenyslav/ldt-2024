package model

import (
	"time"

	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// QueryDto is a domain layer representation of query.
type QueryDto struct {
	CreatedAt time.Time `json:"created_at"`
	Prompt    string    `json:"prompt,omitempty"`
	Command   string    `json:"command,omitempty"`
	ID        int64     `json:"id"`
}

// ResponseDto is a domain layer representation of response.
type ResponseDto struct {
	CreatedAt time.Time             `json:"created_at"`
	Body      string                `json:"body"`
	Status    shared.ResponseStatus `json:"status"`
}

// QueryCreateReq is a struct for creating new query request.
type QueryCreateReq struct {
	Prompt  string `json:"prompt,omitempty" validate:"gte=5"`
	Command string `json:"command,omitempty"`
}
