package controller

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"go.opentelemetry.io/otel/trace"
)

type chatRepo interface {
	BeginTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
	InsertQuery(ctx context.Context, params model.QueryDao) (int64, error)
	InsertResponse(ctx context.Context, params model.ResponseDao) error
	UpdateQueryMeta(ctx context.Context, params model.QueryMeta, id int64) error
	UpdateResponseStatus(ctx context.Context, id int64, status shared.ResponseStatus) error
}

// Controller is a struct that implements chat business logic.
type Controller struct {
	repo      chatRepo
	tracer    trace.Tracer
	kc        *gocloak.GoCloak
	prompter  pb.PrompterClient
	realm     string
	cipherKey string
}

// New creates new Controller.
func New(repo chatRepo, prompter pb.PrompterClient, kc *gocloak.GoCloak, realm, cipher string, tracer trace.Tracer) *Controller {
	return &Controller{
		repo:      repo,
		kc:        kc,
		realm:     realm,
		cipherKey: cipher,
		tracer:    tracer,
		prompter:  prompter,
	}
}
