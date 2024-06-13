package state

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

// Machine стейт машина для управления данными пользователя.
type Machine struct {
	client storage.Cache
}

// New создает новую Machine.
func New(client storage.Cache) *Machine {
	return &Machine{client: client}
}

func getKey(pref string, userID int64) string {
	return fmt.Sprintf("%s:%d", pref, userID)
}

// SetState задать state для пользователя.
func (m *Machine) SetState(ctx context.Context, userID int64, state shared.UserState) error {
	return m.client.SetPrimitive(ctx, getKey("state", userID), int(state), shared.UserStateExp)
}

// GetState получить state для пользователя.
func (m *Machine) GetState(ctx context.Context, userID int64) (shared.UserState, error) {
	state, err := m.client.GetInt(ctx, getKey("state", userID))
	return shared.UserState(state), err
}

// SetToken задать токен для пользователя.
func (m *Machine) SetToken(ctx context.Context, userID int64, token string) error {
	return m.client.SetPrimitive(ctx, getKey("token", userID), token, shared.UserStateExp)
}

// GetToken получить токен для пользователя.
func (m *Machine) GetToken(ctx context.Context, userID int64) (string, error) {
	return m.client.GetString(ctx, getKey("token", userID))
}

type redisRoles struct {
	Roles []shared.UserRole `json:"roles"`
}

// SetRolesFromStrings задать роли для пользователя.
func (m *Machine) SetRolesFromStrings(ctx context.Context, userID int64, r []string) error {
	roles := make([]shared.UserRole, len(r))
	for idx, role := range r {
		roles[idx] = shared.RoleFromString(role)
	}
	return m.client.SetStruct(ctx, getKey("roles", userID), redisRoles{roles}, shared.UserStateExp)
}

// GetRoles получить роли для пользователя.
func (m *Machine) GetRoles(ctx context.Context, userID int64) ([]shared.UserRole, error) {
	var roles redisRoles
	err := m.client.GetStruct(ctx, &roles, getKey("roles", userID))
	return roles.Roles, err
}

// SetSessionID задать sessionID для пользователя.
func (m *Machine) SetSessionID(ctx context.Context, userID int64, sessionID uuid.UUID) error {
	return m.client.SetPrimitive(ctx, getKey("session_id", userID), sessionID.String(), shared.UserStateExp)
}

// GetSessionID получить sessionID для пользователя.
func (m *Machine) GetSessionID(ctx context.Context, userID int64) (uuid.UUID, error) {
	sessionID, err := m.client.GetString(ctx, getKey("session_id", userID))
	if err != nil {
		return uuid.New(), err
	}
	return uuid.Parse(sessionID)
}
