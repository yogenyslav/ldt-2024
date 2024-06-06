package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/pkg/storage"
)

// Repo is a repository for the session model.
type Repo struct {
	pg storage.SQLDatabase
}

// New creates a new session repository.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

const insertOne = `
	insert into chat.session(id, username, title)
	values ($1, $2, $3);
`

// InsertOne create new session.
func (r *Repo) InsertOne(ctx context.Context, params model.SessionDao) error {
	_, err := r.pg.Exec(ctx, insertOne, params.ID, params.Username, params.Title)
	return err
}

const list = `
	select id, title, created_at
	from chat.session
	where username = $1
	order by created_at;
`

// List get list of sessions by username.
func (r *Repo) List(ctx context.Context, username string) ([]model.SessionDao, error) {
	var sessions []model.SessionDao
	err := r.pg.QuerySlice(ctx, &sessions, list, username)
	return sessions, err
}
