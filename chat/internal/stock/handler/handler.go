package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
)

// Handler имплементация сервера регулярных товаров.
type Handler struct {
	predictor pb.PredictorClient
	validator *validator.Validate
}

// New создает новый Handler.
func New(predictor pb.PredictorClient) *Handler {
	return &Handler{
		predictor: predictor,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
