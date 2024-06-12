package handler

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"go.opentelemetry.io/otel/trace"
)

// Handler имплементация сервиса Predictor.
type Handler struct {
	pb.UnimplementedPredictorServer
	predictor pb.PredictorClient
	tracer    trace.Tracer
}

// New создает новый Handler.
func New(predictor client.GrpcClient, tracer trace.Tracer) *Handler {
	return &Handler{
		predictor: pb.NewPredictorClient(predictor.GetConn()),
		tracer:    tracer,
	}
}
