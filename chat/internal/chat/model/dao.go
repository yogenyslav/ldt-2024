package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
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

// ToDto converts QueryDao to QueryDto.
func (q QueryDao) ToDto() QueryDto {
	return QueryDto{
		CreatedAt: q.CreatedAt,
		Prompt:    q.Prompt,
		Command:   q.Command,
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
