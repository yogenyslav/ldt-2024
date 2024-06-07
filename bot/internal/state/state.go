package state

import (
	"context"
	"fmt"

	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

type Machine struct {
	client storage.Cache
}

func New(client storage.Cache) *Machine {
	return &Machine{client: client}
}

func getKey(pref string, userId int64) string {
	return fmt.Sprintf("%s:%d", pref, userId)
}

func (m *Machine) SetState(ctx context.Context, userId int64, state shared.State) error {
	return m.client.SetPrimitive(ctx, getKey("state", userId), int(state), shared.UserStateExp)
}

func (m *Machine) GetState(ctx context.Context, userId int64) (shared.State, error) {
	state, err := m.client.GetInt(ctx, getKey("state", userId))
	return shared.State(state), err
}
