package controller

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	cr "github.com/yogenyslav/ldt-2024/chat/internal/chat/repo"
	"github.com/yogenyslav/ldt-2024/chat/tests/database"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	cfg      = config.MustNew("../../../config/config_test.yaml")
	exporter = tracetest.NewNoopExporter()
	provider = tracing.MustNewTraceProvider(exporter, "test")
	tracer   = otel.Tracer("test")
	pg       = postgres.MustNew(cfg.Postgres, tracer)
	kc       = gocloak.NewClient(cfg.KeyCloak.URL)
)

func init() {
	otel.SetTracerProvider(provider)
}

func TestController_InsertQuery(t *testing.T) {
	ctx := context.Background()
	repo := cr.New(pg)
	ctrl := New(repo, kc, cfg.KeyCloak.Realm, cfg.Server.CipherKey, tracer)

	defer database.TruncateTable(t, ctx, pg, "chat.query")

	tests := []struct {
		name      string
		params    model.QueryCreateReq
		username  string
		sessionID uuid.UUID
	}{
		{
			name:      "success, prompt",
			params:    model.QueryCreateReq{Prompt: "test prompt"},
			username:  "user",
			sessionID: uuid.New(),
		},
		{
			name:      "success, command",
			params:    model.QueryCreateReq{Command: "test command"},
			username:  "user",
			sessionID: uuid.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ctrl.InsertQuery(ctx, tt.params, tt.username, tt.sessionID)
			assert.NoError(t, err)
		})
	}
}
