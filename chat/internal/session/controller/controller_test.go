//go:build integration

package controller

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ldt-2024/chat/config"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/tests/database"
	"github.com/yogenyslav/ldt-2024/chat/tests/fixtures"
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
)

func init() {
	otel.SetTracerProvider(provider)
}

func fillDB(t *testing.T, ctx context.Context, repo sessionRepo, sessions ...model.SessionDao) {
	t.Helper()
	var err error
	for _, session := range sessions {
		err = repo.InsertOne(ctx, session)
		require.NoError(t, err)
	}
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
			username:  "user",
			sessionID: sessionID,
			wantErr:   nil,
		},
		{
			name:      "fail, duplicate id",
			username:  "user",
			sessionID: sessionID,
			wantErr:   shared.ErrSessionDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ctrl.NewSession(ctx, tt.sessionID, tt.username)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestController_List(t *testing.T) {
	ctx := context.Background()
	repo := sr.New(pg)
	ctrl := New(repo, tracer)

	toFill := []model.SessionDao{
		fixtures.SessionDao().New().V(),
		fixtures.SessionDao().New().V(),
		fixtures.SessionDao().New().Username("another_user").V(),
	}
	fillDB(t, ctx, repo, toFill...)

	defer database.TruncateTable(t, ctx, pg, "chat.session")

	tests := []struct {
		name     string
		username string
		wantLen  int
	}{
		{
			name:     "success, 2 sessions",
			username: "user",
			wantLen:  2,
		},
		{
			name:     "success, 1 session",
			username: "another_user",
			wantLen:  1,
		},
		{
			name:     "success, 0 sessions",
			username: "random_user",
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions, err := ctrl.List(ctx, tt.username)
			require.NoError(t, err)
			assert.Equal(t, tt.wantLen, len(sessions))
		})
	}
}

func TestController_Rename(t *testing.T) {
	ctx := context.Background()
	repo := sr.New(pg)
	ctrl := New(repo, tracer)

	sessionID := uuid.New()

	toFill := []model.SessionDao{
		fixtures.SessionDao().New().ID(sessionID).V(),
	}
	fillDB(t, ctx, repo, toFill...)

	defer database.TruncateTable(t, ctx, pg, "chat.session")

	tests := []struct {
		name    string
		id      uuid.UUID
		title   string
		wantErr error
	}{
		{
			name:    "success",
			id:      sessionID,
			title:   "new title",
			wantErr: nil,
		},
		{
			name:    "fail, no such id",
			id:      uuid.New(),
			title:   "new title",
			wantErr: shared.ErrNoSessionWithID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ctrl.Rename(ctx, model.RenameReq{
				ID:    tt.id,
				Title: tt.title,
			})

			if tt.wantErr == nil {
				assert.NoError(t, err)

				sessions, err := ctrl.List(ctx, "user")
				require.NoError(t, err)

				assert.Equal(t, tt.title, sessions[0].Title)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestController_Delete(t *testing.T) {
	ctx := context.Background()
	repo := sr.New(pg)
	ctrl := New(repo, tracer)

	sessionID := uuid.New()

	toFill := []model.SessionDao{
		fixtures.SessionDao().New().ID(sessionID).V(),
	}
	fillDB(t, ctx, repo, toFill...)

	defer database.TruncateTable(t, ctx, pg, "chat.session")

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr error
	}{
		{
			name:    "success",
			id:      sessionID,
			wantErr: nil,
		},
		{
			name:    "fail, no such id",
			id:      uuid.New(),
			wantErr: shared.ErrNoSessionWithID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ctrl.Delete(ctx, tt.id)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
