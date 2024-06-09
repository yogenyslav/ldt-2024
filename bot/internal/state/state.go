package state

import (
	"context"
	"fmt"

	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

// Machine is a state machine for the bot.
type Machine struct {
	client storage.Cache
}

// New creates a new state machine.
func New(client storage.Cache) *Machine {
	return &Machine{client: client}
}

func getKey(pref string, userID int64) string {
	return fmt.Sprintf("%s:%d", pref, userID)
}

// SetState sets the state for the user.
func (m *Machine) SetState(ctx context.Context, userID int64, state shared.State) error {
	return m.client.SetPrimitive(ctx, getKey("state", userID), int(state), shared.UserStateExp)
}

// GetState gets the state for the user.
func (m *Machine) GetState(ctx context.Context, userID int64) (shared.State, error) {
	state, err := m.client.GetInt(ctx, getKey("state", userID))
	return shared.State(state), err
}

// SetToken sets the JWT token for the user.
func (m *Machine) SetToken(ctx context.Context, userID int64, token string) error {
	return m.client.SetPrimitive(ctx, getKey("token", userID), token, shared.UserStateExp)
}

// GetToken gets the JWT token for the user.
func (m *Machine) GetToken(ctx context.Context, userID int64) (string, error) {
	return m.client.GetString(ctx, getKey("token", userID))
}
