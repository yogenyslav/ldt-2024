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

const findMeta = `
	select username, title, is_deleted, tg
	from chat.session
	where id = $1;
`

// FindMeta returns session meta info.
func (r *Repo) FindMeta(ctx context.Context, id uuid.UUID) (model.SessionMeta, error) {
	var status model.SessionMeta
	err := r.pg.Query(ctx, &status, findMeta, id)
	return status, err
}

const findContent = `
	select 
		(r.created_at, r.body, r.status) as response,
		(q.created_at, q.prompt, q.product, q.status, q.type, q.id) as query
	from chat.query q
	join
	    chat.response r
		on q.id = r.query_id
	join 
	    chat.session s
		on q.session_id = s.id
	where
	    q.session_id = $1 
	  	and s.is_deleted = false;
`

// FindContent finds session content by id.
func (r *Repo) FindContent(ctx context.Context, id uuid.UUID) ([]model.SessionContentDao, error) {
	var content []model.SessionContentDao
	err := r.pg.QuerySlice(ctx, &content, findContent, id)
	return content, err
}

const sessionClenaup = `
	select count(id)
	from chat.query
	where session_id = $1;
`

// SessionContentEmpty checks if session content is empty.
func (r *Repo) SessionContentEmpty(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	var count int64
	err := r.pg.Query(ctx, &count, sessionClenaup, sessionID)
	return count == 0, err
}
