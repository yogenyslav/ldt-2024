package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/pkg/storage"
)

// Repo chat repository.
type Repo struct {
	pg storage.SQLDatabase
}

// New creates new Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

const insertQuery = `
	insert into chat.query(prompt, command, username, session_id)
	values ($1, $2, $3, $4)
	returning id;
`

// InsertQuery creates new query and returns its id.
func (r *Repo) InsertQuery(ctx context.Context, params model.QueryDao) (int64, error) {
	var id int64
	err := r.pg.Query(ctx, &id, insertQuery, params.Prompt, params.Command, params.Username, params.SessionID)
	return id, err
}
