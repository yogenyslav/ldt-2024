package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// QueryDao модель запроса в базе данных.
type QueryDao struct { //nolint:govet // order is required for sql request
	CreatedAt time.Time          `db:"created_at"`
	Prompt    string             `db:"prompt"`
	Product   string             `db:"product"`
	Period    string             `db:"period"`
	Status    shared.QueryStatus `db:"status"`
	Type      shared.QueryType   `db:"type"`
	ID        int64              `db:"id"`
	Username  string             `db:"username"`
	SessionID uuid.UUID          `db:"session_id"`
}

// ToDto конвертирует QueryDao в QueryDto.
func (q QueryDao) ToDto() QueryDto {
	return QueryDto{
		CreatedAt: q.CreatedAt,
		Prompt:    q.Prompt,
		Status:    q.Status.ToString(),
		Period:    q.Period,
		Product:   q.Product,
		Type:      q.Type.ToString(),
		ID:        q.ID,
	}
}

// ResponseDao модель ответа в базе данных.
type ResponseDao struct {
	CreatedAt time.Time             `db:"created_at"`
	Body      string                `db:"body"`
	Data      []byte                `db:"data"`
	DataType  shared.QueryType      `db:"data_type"`
	Status    shared.ResponseStatus `db:"status"`
	QueryID   int64                 `db:"query_id"`
}

// ToDto конвертирует ResponseDao в ResponseDto.
func (r ResponseDao) ToDto() ResponseDto {
	resp := make(map[string]any)
	if err := json.Unmarshal(r.Data, &resp); err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return ResponseDto{
			CreatedAt: r.CreatedAt,
			Body:      r.Body,
			DataType:  r.DataType.ToString(),
			Data:      nil,
			Status:    r.Status.ToString(),
		}
	}

	return ResponseDto{
		CreatedAt: r.CreatedAt,
		DataType:  r.DataType.ToString(),
		Data:      resp,
		Body:      r.Body,
		Status:    r.Status.ToString(),
	}
}
