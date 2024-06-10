package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
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

// BeginTx starts new transaction.
func (r *Repo) BeginTx(ctx context.Context) (context.Context, error) {
	ctx, err := r.pg.BeginSerializable(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// CommitTx commits transaction.
func (r *Repo) CommitTx(ctx context.Context) error {
	return r.pg.CommitTx(ctx)
}

// RollbackTx rollbacks transaction.
func (r *Repo) RollbackTx(ctx context.Context) error {
	return r.pg.RollbackTx(ctx)
}

const insertQuery = `
	insert into chat.query(prompt, command, username, session_id)
	values ($1, $2, $3, $4)
	returning id;
`

// InsertQuery creates new query and returns its id.
func (r *Repo) InsertQuery(ctx context.Context, params model.QueryDao) (int64, error) {
	var id int64

	err := r.pg.QueryTx(ctx, &id, insertQuery, params.Prompt, params.Command, params.Username, params.SessionID)
	return id, err
}

const insertResponse = `
	insert into chat.response(query_id)
	values ($1);
`

// InsertResponse create new response with processing status.
func (r *Repo) InsertResponse(ctx context.Context, params model.ResponseDao) error {
	_, err := r.pg.ExecTx(ctx, insertResponse, params.QueryID)
	return err
}

const updateQueryMeta = `
	update chat.query
	set product = $2, type = $3
	where id = $1;
`

// UpdateQueryMeta updates metadata for query after passing through prompter.
func (r *Repo) UpdateQueryMeta(ctx context.Context, params model.QueryMeta, id int64) error {
	_, err := r.pg.ExecTx(ctx, updateQueryMeta, id, params.Product, params.Type)
	return err
}

const updateResponse = `
	update chat.response
	set body = $2, status = $3
	where query_id = $1;
`

// UpdateResponse updates response status by query id.
func (r *Repo) UpdateResponse(ctx context.Context, id int64, params model.ResponseDao) error {
	tag, err := r.pg.Exec(ctx, updateResponse, id, params.Body, params.Status)
	if err != nil {
		return shared.ErrUpdateResponse
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoResponseWithID
	}
	return nil
}
