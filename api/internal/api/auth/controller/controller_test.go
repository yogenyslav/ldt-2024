//go:build integration

package controller

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/ldt-2024/api/config"
	"github.com/yogenyslav/ldt-2024/api/internal/api/auth/model"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	cfg      = config.MustNew("../../../../config/config_test.yaml")
	exporter = tracetest.NewNoopExporter()
	provider = tracing.MustNewTraceProvider(exporter, "test")
	tracer   = otel.Tracer("test")
	kc       = gocloak.NewClient(cfg.KeyCloak.URL)
)

func init() {
	otel.SetTracerProvider(provider)
}

func TestController_Login(t *testing.T) {
	ctx := context.Background()
	ctrl := New(kc, cfg.KeyCloak, tracer)

	token, err := kc.LoginAdmin(ctx, cfg.KeyCloak.User, cfg.KeyCloak.Password, cfg.KeyCloak.AdminRealm)
	require.NoError(t, err)

	user := gocloak.User{
		FirstName: gocloak.StringP("test_name"),
		LastName:  gocloak.StringP("test_last_name"),
		Email:     gocloak.StringP("test@test.com"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP("test_user"),
	}
	userID, err := kc.CreateUser(ctx, token.AccessToken, cfg.KeyCloak.Realm, user)
	require.NoError(t, err)

	err = kc.SetPassword(ctx, token.AccessToken, userID, cfg.KeyCloak.Realm, "test123456", false)
	require.NoError(t, err)

	defer func() {
		err = kc.DeleteUser(ctx, token.AccessToken, cfg.KeyCloak.Realm, userID)
		require.NoError(t, err)
	}()

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "success",
			username: "test@test.com",
			password: "test123456",
			wantErr:  false,
		},
		{
			name:     "success",
			username: "test_user",
			password: "test123456",
			wantErr:  false,
		},
		{
			name:     "fail, invalid username",
			username: "aaa@test.com",
			password: "test123456",
			wantErr:  true,
		},
		{
			name:     "fail, invalid password",
			username: "test@test.com",
			password: "test12345",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := model.LoginReq{
				Username: tt.username,
				Password: tt.password,
			}
			resp, err := ctrl.Login(ctx, req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}
