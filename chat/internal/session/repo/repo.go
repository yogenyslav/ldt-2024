package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
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
	where username = $1 and is_deleted = false
	order by created_at;
`

// List get list of sessions by username.
func (r *Repo) List(ctx context.Context, username string) ([]model.SessionDao, error) {
	var sessions []model.SessionDao
	err := r.pg.QuerySlice(ctx, &sessions, list, username)
	return sessions, err
}

const updateTitle = `
	update chat.session
	set title = $2
	where id = $1 and is_deleted = false;
`

// UpdateTitle updates session title filtered by id.
func (r *Repo) UpdateTitle(ctx context.Context, params model.RenameReq) error {
	tag, err := r.pg.Exec(ctx, updateTitle, params.ID, params.Title)
	if err != nil {
		return shared.ErrUpdateSession
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoSessionWithID
	}
	return err
}

const deleteOne = `
	update chat.session
	set is_deleted = true
	where id = $1;
`

// DeleteOne deletes a session by id.
func (r *Repo) DeleteOne(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pg.Exec(ctx, deleteOne, id)
	if err != nil {
		return shared.ErrDeleteSession
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoSessionWithID
	}
	return nil
}
