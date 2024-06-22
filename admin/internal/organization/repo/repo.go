package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/pkg/storage"
)

// Repo репозиторий организаций.
type Repo struct {
	pg storage.SQLDatabase
}

// New создает новый Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

const insertOne = `
	insert into adm.organization (username, title)
	values ($1, $2)
	returning id;
`

const addToOrganization = `
	insert into adm.user_organization (username, organization_id)
	values ($1, $2);
`

// InsertOne вставляет новую организацию.
func (r *Repo) InsertOne(ctx context.Context, params model.OrganizationDao) (int64, error) {
	var id int64
	tx, err := r.pg.BeginSerializable(ctx)
	if err != nil {
		return id, err
	}
	defer r.pg.RollbackTx(tx)

	if err := r.pg.QueryTx(tx, &id, insertOne, params.Username, params.Title); err != nil {
		return id, err
	}

	if _, err := r.pg.ExecTx(tx, addToOrganization, params.Username, id); err != nil {
		return id, err
	}
	return id, r.pg.CommitTx(tx)
}

const listForUser = `
	select created_at, username, title, id
	from adm.organization
	where username = $1;
`

// ListForUser находит список организаций по username.
func (r *Repo) ListForUser(ctx context.Context, username string) ([]model.OrganizationDao, error) {
	var org []model.OrganizationDao
	err := r.pg.QuerySlice(ctx, &org, listForUser, username)
	return org, err
}

const updateOne = `
	update adm.organization
	set title = $3
	where id = $1 and username = $2;
`

// UpdateOne обновить название организации.
func (r *Repo) UpdateOne(ctx context.Context, params model.OrganizationDao) error {
	tag, err := r.pg.Exec(ctx, updateOne, params.ID, params.Username, params.Title)
	if err != nil {
		return shared.ErrUpdateOrganization
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoOrganization
	}
	return nil
}
