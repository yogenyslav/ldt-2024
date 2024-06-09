package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/pkg/storage"
)

func TruncateTable(t *testing.T, ctx context.Context, pg storage.SQLDatabase, tables ...string) {
	t.Helper()
	_, err := pg.Exec(ctx, fmt.Sprintf(`
		truncate %s restart identity
	`, strings.Join(tables, ",")))
	require.NoError(t, err)
}
