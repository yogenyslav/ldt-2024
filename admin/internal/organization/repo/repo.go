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
	insert into adm.organization (username, title, s3_bucket)
	values ($1, $2, $3)
	returning id;
`

// InsertOne вставляет новую организацию.
func (r *Repo) InsertOne(ctx context.Context, params model.OrganizationDao) (int64, error) {
	var id int64
	err := r.pg.Query(ctx, &id, insertOne, params.Username, params.Title, params.S3Bucket)
	return id, err
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

const updateOne = `
	update adm.organization
	set title = $3
	where id = $1 and username = $2;
`

// UpdateOne обновляет организацию.
func (r *Repo) UpdateOne(ctx context.Context, params model.OrganizationDao) error {
	tag, err := r.pg.Exec(ctx, updateOne, params.ID, params.Username, params.Title)
	if err != nil {
		return shared.ErrCreateOrganization
	}
	if tag.RowsAffected() == 0 {
		return shared.ErrNoOrganization
	}
	return err
}
