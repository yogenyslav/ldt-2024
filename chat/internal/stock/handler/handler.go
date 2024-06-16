package handler

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
)

// Handler имплементация сервера регулярных товаров.
type Handler struct {
	predictor pb.PredictorClient
}

// New создает новый Handler.
func New(predictor pb.PredictorClient) *Handler {
	return &Handler{
		predictor: predictor,
	}
}
