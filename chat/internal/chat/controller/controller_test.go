//go:build integration

//go:generate mockgen -source=../../api/pb/prompter_grpc.pb.go -destination=./mocks/prompter.go -package=mocks
package controller

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/controller/mocks"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	cr "github.com/yogenyslav/ldt-2024/chat/internal/chat/repo"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/tests/database"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/mock/gomock"
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

	ctrl := gomock.NewController(t)
	prompter := mocks.NewMockPrompterClient(ctrl)
	defer ctrl.Finish()

	chatRepository := cr.New(pg)
	sessionRepository := sr.New(pg)

	controller := New(chatRepository, sessionRepository, prompter, kc, cfg.KeyCloak.Realm, cfg.Server.CipherKey, tracer)

	defer database.TruncateTable(t, ctx, pg, "chat.query")

	tests := []struct {
		name      string
		params    model.QueryCreateReq
		username  string
		sessionID uuid.UUID
		wantErr   error
	}{
		{
			name:      "success, prompt",
			params:    model.QueryCreateReq{Prompt: "test prompt"},
			username:  "user",
			sessionID: uuid.New(),
			wantErr:   nil,
		},
		{
			name:      "fail, empty prompt",
			params:    model.QueryCreateReq{Command: "test command"},
			username:  "user",
			sessionID: uuid.New(),
			wantErr:   shared.ErrEmptyQueryHint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.params.Prompt != "" {
				prompter.EXPECT().Extract(gomock.Any(), gomock.Any())
			}
			_, err := controller.InsertQuery(ctx, tt.params, tt.username, tt.sessionID)
			t.Log(err)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestController_InsertResponse(t *testing.T) {
	ctx := context.Background()
	database.TruncateTable(t, ctx, pg, "chat.response")

	ctrl := gomock.NewController(t)
	prompter := mocks.NewMockPrompterClient(ctrl)
	defer ctrl.Finish()

	chatRepository := cr.New(pg)
	sessionRepository := sr.New(pg)
	controller := New(chatRepository, sessionRepository, prompter, kc, cfg.KeyCloak.Realm, cfg.Server.CipherKey, tracer)

	defer database.TruncateTable(t, ctx, pg, "chat.response")

	tests := []struct {
		name    string
		queryID int64
	}{
		{
			name:    "success",
			queryID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := chatRepository.BeginTx(ctx)
			require.NoError(t, err)

			err = controller.InsertResponse(tx, tt.queryID)
			if assert.NoError(t, err) {
				require.NoError(t, chatRepository.CommitTx(tx))
			} else {
				require.NoError(t, chatRepository.RollbackTx(tx))
			}
		})
	}
}
