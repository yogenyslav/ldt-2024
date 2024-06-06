package controller

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ldt-2024/chat/config"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/test/database"
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
	kc       = gocloak.NewClient(cfg.KeyCloak.URL)
	pg       = postgres.MustNew(cfg.Postgres, tracer)
)

func init() {
	otel.SetTracerProvider(provider)
}

func TestController_NewSession(t *testing.T) {
	ctx := context.Background()
	repo := sr.New(pg)
	ctrl := New(repo, tracer)
	sessionID := uuid.New()

	defer database.TruncateTable(t, ctx, pg, "chat.session")

	tests := []struct {
		name      string
		username  string
		password  string
		sessionID uuid.UUID
		wantErr   error
	}{
		{
			name:      "success",
			username:  cfg.KeyCloak.User,
			password:  cfg.KeyCloak.Password,
			sessionID: sessionID,
			wantErr:   nil,
		},
		{
			name:      "fail, duplicate id",
			username:  cfg.KeyCloak.User,
			password:  cfg.KeyCloak.Password,
			sessionID: sessionID,
			wantErr:   shared.ErrSessionDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kc.LoginAdmin(ctx, cfg.KeyCloak.User, cfg.KeyCloak.Password, cfg.KeyCloak.AdminRealm)
			require.NoError(t, err)

			err = ctrl.NewSession(ctx, tt.sessionID, tt.username)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
