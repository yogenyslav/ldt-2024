package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
	"go.opentelemetry.io/otel/trace"
)

type chatController interface {
	InsertQuery(ctx context.Context, params model.QueryCreateReq, username string, sessionID uuid.UUID) (model.QueryDto, error)
	Authorize(ctx context.Context, token string) (context.Context, string, error)
	Predict(ctx context.Context, out chan<- chatresp.Response, cancel <-chan struct{}, query int64)
	Hint(ctx context.Context, queryID int64, params model.QueryCreateReq) (model.QueryDto, error)
	UpdateStatus(ctx context.Context, queryID int64, status shared.QueryStatus) error
	SessionCleanup(ctx context.Context, sessionID uuid.UUID) error
}

type sessionController interface {
	NewSession(ctx context.Context, id uuid.UUID, username string) error
}

// Handler обработчик чата.
type Handler struct {
	machine *state.Machine
	cc      chatController
	sc      sessionController
	tracer  trace.Tracer
}

// New создает новый Handler.
func New(cc chatController, sc sessionController, machine *state.Machine, tracer trace.Tracer) *Handler {
	return &Handler{
		cc:      cc,
		sc:      sc,
		machine: machine,
		tracer:  tracer,
	}
}
