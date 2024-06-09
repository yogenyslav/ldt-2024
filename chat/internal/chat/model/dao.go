package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// QueryDao is a data layer representation of users query.
type QueryDao struct { //nolint:govet // order is required for sql request
	CreatedAt time.Time        `db:"created_at"`
	Prompt    string           `db:"prompt"`
	Command   string           `db:"command"`
	Product   string           `db:"product"`
	Type      shared.QueryType `db:"type"`
	ID        int64            `db:"id"`
	Username  string           `db:"username"`
	SessionID uuid.UUID        `db:"session_id"`
}

// ToDto converts QueryDao to QueryDto.
func (q QueryDao) ToDto() QueryDto {
	return QueryDto{
		CreatedAt: q.CreatedAt,
		Prompt:    q.Prompt,
		Command:   q.Command,
		Product:   q.Product,
		Type:      q.Type,
		ID:        q.ID,
	}
}

// ResponseDao is a data layer representation of response for query.
type ResponseDao struct {
	CreatedAt time.Time             `db:"created_at"`
	Body      string                `db:"body"`
	Status    shared.ResponseStatus `db:"status"`
	QueryID   int64                 `db:"query_id"`
}

// ToDto converts ResponseDao to ResponseDto.
func (r ResponseDao) ToDto() ResponseDto {
	return ResponseDto{
		CreatedAt: r.CreatedAt,
		Body:      r.Body,
		Status:    r.Status,
	}
}
