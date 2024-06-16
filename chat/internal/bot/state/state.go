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

// SetValidateQuery задать id Query, который необходимо провалидировать.
func (m *Machine) SetValidateQuery(ctx context.Context, userID, queryID int64) error {
	return m.client.SetPrimitive(ctx, getKey("validate_query", userID), queryID, shared.UserStateExp)
}

// GetValidateQuery получить id Query, который надо провалидировать.
func (m *Machine) GetValidateQuery(ctx context.Context, userID int64) (int64, error) {
	return m.client.GetInt64(ctx, getKey("validate_query", userID))
}

// SetHintQuery задать id Query, который надо подсказать.
func (m *Machine) SetHintQuery(ctx context.Context, userID, queryID int64) error {
	return m.client.SetPrimitive(ctx, getKey("hint_query", userID), queryID, shared.UserStateExp)
}

// GetHintQuery получить id Query, который надо подсказать.
func (m *Machine) GetHintQuery(ctx context.Context, userID int64) (int64, error) {
	return m.client.GetInt64(ctx, getKey("hint_query", userID))
}
