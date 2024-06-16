package handler

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/client"
	"github.com/yogenyslav/ldt-2024/api/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
)

// Handler имплементация сервиса Predictor.
type Handler struct {
	pb.UnimplementedPredictorServer
	predictor pb.PredictorClient
	m         *metrics.Metrics
	tracer    trace.Tracer
}

// New создает новый Handler.
func New(predictor client.GrpcClient, m *metrics.Metrics, tracer trace.Tracer) *Handler {
	return &Handler{
		predictor: pb.NewPredictorClient(predictor.GetConn()),
		m:         m,
		tracer:    tracer,
	}
}
