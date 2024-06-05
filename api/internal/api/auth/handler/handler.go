package handler

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/api/pkg/metrics"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type authController interface {
	Login(ctx context.Context, params model.LoginReq) (model.LoginResp, error)
}

// Handler is the struct that implements the AuthServiceServer interface.
type Handler struct {
	pb.UnimplementedAuthServiceServer
	ctrl    authController
	tracer  trace.Tracer
	metrics *metrics.Metrics
}

// New creates a new Handler.
func New(ctrl authController, tracer trace.Tracer, m *metrics.Metrics) *Handler {
	return &Handler{
		ctrl:    ctrl,
		tracer:  tracer,
		metrics: m,
	}
}

func getTraceCtx(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIDVal, ok := md["x-trace-id"]
	if !ok {
		log.Debug().Msg("trace id not found")
		return nil, nil
	}

	traceIDString := traceIDVal[0]
	traceID, err := trace.TraceIDFromHex(traceIDString)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to parse trace id")
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)
	return ctx, nil
}
