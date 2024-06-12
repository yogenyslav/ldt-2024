package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

// Repo репозиторий чата.
type Repo struct {
	pg storage.SQLDatabase
}

// New создает новый Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

// BeginTx начинает транзакцию.
func (r *Repo) BeginTx(ctx context.Context) (context.Context, error) {
	ctx, err := r.pg.BeginSerializable(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// CommitTx коммитит транзакцию.
func (r *Repo) CommitTx(ctx context.Context) error {
	return r.pg.CommitTx(ctx)
}

// RollbackTx откатывает транзакцию.
func (r *Repo) RollbackTx(ctx context.Context) error {
	return r.pg.RollbackTx(ctx)
}

const insertQuery = `
	insert into chat.query(prompt, username, session_id, status)
	values ($1, $2, $3, $4)
	returning id;
`

// InsertQuery создает новый запрос и возвращает его id.
func (r *Repo) InsertQuery(ctx context.Context, params model.QueryDao) (int64, error) {
	var id int64

	err := r.pg.QueryTx(ctx, &id, insertQuery, params.Prompt, params.Username, params.SessionID, params.Status)
	return id, err
}

const insertResponse = `
	insert into chat.response(query_id)
	values ($1);
`

// InsertResponse создает новый ответ со статусом "processing".
func (r *Repo) InsertResponse(ctx context.Context, params model.ResponseDao) error {
	_, err := r.pg.ExecTx(ctx, insertResponse, params.QueryID)
	return err
}

const updateQueryMeta = `
	update chat.query
	set product = $2, type = $3, period = $4
	where id = $1;
`

// UpdateQueryMeta обновляет метаданные запроса по id.
func (r *Repo) UpdateQueryMeta(ctx context.Context, params model.QueryMeta, id int64) error {
	_, err := r.pg.ExecTx(ctx, updateQueryMeta, id, params.Product, params.Type, params.Period)
	return err
}

const updateResponse = `
	update chat.response
	set body = $2, status = $3
	where query_id = $1;
`

// UpdateResponse обновляет ответ по id.
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

// FindQueryPrompt возвращает текст запроса по id.
func (r *Repo) FindQueryPrompt(ctx context.Context, id int64) (string, error) {
	var prompt string
	err := r.pg.Query(ctx, &prompt, findQueryPrompt, id)
	return prompt, err
}

const updateQuery = `
	update chat.query
	set prompt = $2, status = $3, product = $4, type = $5, period = $6
	where id = $1;
`

// UpdateQuery обновляет запрос по id.
func (r *Repo) UpdateQuery(ctx context.Context, params model.QueryDao) error {
	tag, err := r.pg.Exec(ctx, updateQuery, params.ID, params.Prompt, params.Status, params.Product, params.Type, params.Period)
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

// UpdateQueryStatus обновляет статус запроса по id.
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

const findQueryMeta = `
	select product, period, type
	from chat.query
	where id = $1;
`

// FindQueryMeta получить метаданные запроса по id.
func (r *Repo) FindQueryMeta(ctx context.Context, id int64) (model.QueryMeta, error) {
	var meta model.QueryMeta
	err := r.pg.Query(ctx, &meta, findQueryMeta, id)
	return meta, err
}
