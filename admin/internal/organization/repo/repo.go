package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
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
	insert into adm.organization (username, title, s3_bucket)
	values ($1, $2, $3)
	returning id;
`

const addToOrganization = `
	insert into adm.user_organization (username, organization)
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

	if err := r.pg.QueryTx(tx, &id, insertOne, params.Username, params.Title, params.S3Bucket); err != nil {
		return id, err
	}

	if _, err := r.pg.ExecTx(tx, addToOrganization, params.Username, params.Title); err != nil {
		return id, err
	}
	return id, r.pg.CommitTx(tx)
}

const findOne = `
	select id, username, title, s3_bucket, created_at
	from adm.organization
	where username = $1;
`

// FindOne находит организацию по username.
func (r *Repo) FindOne(ctx context.Context, username string) (model.OrganizationDao, error) {
	var org model.OrganizationDao
	err := r.pg.Query(ctx, &org, findOne, username)
	return org, err
}