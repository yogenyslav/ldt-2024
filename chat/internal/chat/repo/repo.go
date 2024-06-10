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
	insert into chat.query(prompt, username, session_id, status)
	values ($1, $2, $3, $4)
	returning id;
`

// InsertQuery creates new query and returns its id.
func (r *Repo) InsertQuery(ctx context.Context, params model.QueryDao) (int64, error) {
	var id int64

	err := r.pg.QueryTx(ctx, &id, insertQuery, params.Prompt, params.Username, params.SessionID, params.Status)
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
func (r *Repo) UpdateResponse(ctx context.Context, params model.ResponseDao) error {
	tag, err := r.pg.Exec(ctx, updateResponse, params.QueryID, params.Body, params.Status)
	if err != nil {
		return shared.ErrUpdateResponse
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoResponseWithID
	}
	return nil
}

const findQueryPrompt = `
	select prompt
	from chat.query
	where id = $1;
`

// FindQueryPrompt finds query prompt by id.
func (r *Repo) FindQueryPrompt(ctx context.Context, id int64) (string, error) {
	var prompt string
	err := r.pg.Query(ctx, &prompt, findQueryPrompt, id)
	return prompt, err
}

const updateQuery = `
	update chat.query
	set prompt = $2, status = $3, product = $4, type = $5
	where id = $1;
`

// UpdateQuery updates query by id.
func (r *Repo) UpdateQuery(ctx context.Context, params model.QueryDao) error {
	tag, err := r.pg.Exec(ctx, updateQuery, params.ID, params.Prompt, params.Status, params.Product, params.Type)
	if err != nil {
		return shared.ErrUpdateQuery
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoQueryWithID
	}
	return nil
}

const updateQueryStatus = `
	update chat.query
	set status = $2
	where id = $1;
`

// UpdateQueryStatus updates query status by id.
func (r *Repo) UpdateQueryStatus(ctx context.Context, id int64, status shared.QueryStatus) error {
	tag, err := r.pg.Exec(ctx, updateQueryStatus, id, status)
	if err != nil {
		return shared.ErrUpdateQuery
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoQueryWithID
	}
	return nil
}
