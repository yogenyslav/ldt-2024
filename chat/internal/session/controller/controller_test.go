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
	cc "github.com/yogenyslav/ldt-2024/chat/internal/chat/controller"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/controller/mocks"
	chatmodel "github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	cr "github.com/yogenyslav/ldt-2024/chat/internal/chat/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	sr "github.com/yogenyslav/ldt-2024/chat/internal/session/repo"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/tests/database"
	"github.com/yogenyslav/ldt-2024/chat/tests/fixtures"
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
			err := ctrl.NewSession(ctx, tt.sessionID, tt.username, false, 0)
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

func TestController_FindOne(t *testing.T) {
	kc := gocloak.NewClient(cfg.KeyCloak.URL)

	ctrl := gomock.NewController(t)
	prompter := mocks.NewMockPrompterClient(ctrl)
	defer ctrl.Finish()

	ctx := context.Background()
	qr := cr.New(pg)
	qc := cc.New(qr, prompter, kc, cfg.KeyCloak.Realm, cfg.Server.CipherKey, tracer)
	repo := sr.New(pg)
	controller := New(repo, tracer)

	firstSessionID := uuid.New()
	secondSessionID := uuid.New()

	defer database.TruncateTable(t, ctx, pg, "chat.session", "chat.query", "chat.response")

	sessionsToFill := []model.SessionDao{
		fixtures.SessionDao().New().ID(firstSessionID).V(),
		fixtures.SessionDao().New().ID(secondSessionID).V(),
	}
	fillDB(t, ctx, repo, sessionsToFill...)

	queriesToFill := []struct {
		params    chatmodel.QueryCreateReq
		username  string
		sessionID uuid.UUID
	}{
		{
			params:    chatmodel.QueryCreateReq{Prompt: "test prompt"},
			username:  "user",
			sessionID: firstSessionID,
		},
		{
			params:    chatmodel.QueryCreateReq{Prompt: "test prompt"},
			username:  "user",
			sessionID: firstSessionID,
		},
		{
			params:    chatmodel.QueryCreateReq{Prompt: "test prompt"},
			username:  "user",
			sessionID: firstSessionID,
		},
	}
	for _, query := range queriesToFill {
		prompter.EXPECT().Extract(gomock.Any(), gomock.Any())

		_, err := qc.InsertQuery(ctx, query.params, query.username, query.sessionID)
		require.NoError(t, err)
	}

	tests := []struct {
		name         string
		sessionID    uuid.UUID
		username     string
		wantLen      int
		wantEditable bool
		wantErr      error
	}{
		{
			name:         "success, len 3, editable",
			sessionID:    firstSessionID,
			username:     "user",
			wantLen:      3,
			wantEditable: true,
			wantErr:      nil,
		},
		{
			name:         "success, empty session, editable",
			sessionID:    secondSessionID,
			username:     "user",
			wantLen:      0,
			wantEditable: true,
			wantErr:      nil,
		},
		{
			name:         "success, len 3, not editable",
			sessionID:    firstSessionID,
			username:     "another_user",
			wantLen:      3,
			wantEditable: false,
			wantErr:      nil,
		},
		{
			name:         "success, empty session, not editable",
			sessionID:    secondSessionID,
			username:     "another_user",
			wantLen:      0,
			wantEditable: false,
			wantErr:      nil,
		},
		{
			name:      "fail, no such id",
			sessionID: uuid.New(),
			username:  "user",
			wantErr:   shared.ErrNoSessionWithID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := controller.FindOne(ctx, tt.sessionID, tt.username)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.sessionID, resp.ID)
			assert.Equal(t, tt.wantLen, len(resp.Content))
			assert.Equal(t, tt.wantEditable, resp.Editable)
		})
	}
}
