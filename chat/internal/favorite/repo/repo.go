package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

// Repo репозиторий избранных предиктов.
type Repo struct {
	pg storage.SQLDatabase
}

// New создает новый Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

const insertOne = `
	insert into chat.favorite_responses (query_id, username, response)
	values ($1, $2, $3);
`

// InsertOne добавляет новый предикт в избранное.
func (r *Repo) InsertOne(ctx context.Context, params model.FavoriteDao) error {
	_, err := r.pg.Exec(ctx, insertOne, params.QueryID, params.Username, params.Response)
	return err
}

const list = `
	select query_id, response, created_at, updated_at
	from chat.favorite_responses
	where username = $1 and is_deleted = false;
`

// List возвращает список избранных предиктов.
func (r *Repo) List(ctx context.Context, username string) ([]model.FavoriteDao, error) {
	var favorites []model.FavoriteDao
	err := r.pg.QuerySlice(ctx, &favorites, list, username)
	return favorites, err
}

const findOne = `
	select response, created_at, updated_at
	from chat.favorite_responses
	where query_id = $1 and is_deleted = false;
`

// FindOne возвращает предикт из избранного.
func (r *Repo) FindOne(ctx context.Context, queryID int64) (model.FavoriteDao, error) {
	var favorite model.FavoriteDao
	err := r.pg.Query(ctx, &favorite, findOne, queryID)
	return favorite, err
}

const updateOne = `
	update chat.favorite_responses
	set response = $3, updated_at = current_timestamp
	where query_id = $1 and username = $2 and is_deleted = false;
`

// UpdateOne обновляет предикт в избранном.
func (r *Repo) UpdateOne(ctx context.Context, queryID int64, username string, response []byte) error {
	tag, err := r.pg.Exec(ctx, updateOne, queryID, username, response)
	if err != nil {
		return shared.ErrUpdateFavorite
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoFavoriteWithID
	}
	return nil
}

const deleteOne = `
	update chat.favorite_responses
	set is_deleted = true
	where query_id = $1 and username = $2 and is_deleted = false;
`

// DeleteOne удаляет предикт из избранного.
func (r *Repo) DeleteOne(ctx context.Context, queryID int64, username string) error {
	tag, err := r.pg.Exec(ctx, deleteOne, queryID, username)
	if err != nil {
		return shared.ErrDeleteFavorite
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoFavoriteWithID
	}
	return nil
}

const restoreOne = `
	update chat.favorite_responses
	set is_deleted = false
	where query_id = $1 and username = $2 and is_deleted = true;
`

// RestoreOne восстанавливает предикт в избранное.
func (r *Repo) RestoreOne(ctx context.Context, queryID int64, username string) error {
	tag, err := r.pg.Exec(ctx, restoreOne, queryID, username)
	if err != nil {
		return shared.ErrCreateSession
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoFavoriteWithID
	}
	return nil
}
