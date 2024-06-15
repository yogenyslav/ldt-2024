package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
	"github.com/yogenyslav/pkg/storage"
)

// Repo репозиторий для работы с пользователями и организациями.
type Repo struct {
	pg storage.SQLDatabase
}

// New создает новый Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{
		pg: pg,
	}
}

const insertOrganization = `
	insert into adm.user_organization(username, organization)
	values ($1, $2);
`

// InsertOrganization добавляет пользователя в организацию.
func (r *Repo) InsertOrganization(ctx context.Context, params model.UserOrganizationDao) error {
	_, err := r.pg.Exec(ctx, insertOrganization, params.Username, params.Organization)
	return err
}

const list = `
	select username
	from adm.user_organization
	where organization = $1 and is_deleted = false;
`

// List возвращает список пользователей по организации.
func (r *Repo) List(ctx context.Context, organization string) ([]string, error) {
	var users []string
	err := r.pg.QuerySlice(ctx, &users, list, organization)
	return users, err
}

const deleteOne = `
	update adm.user_organization
	set is_deleted = true
	where username = $1;
`

// DeleteOrganization удаляет пользователя из организации.
func (r *Repo) DeleteOrganization(ctx context.Context, username string) error {
	_, err := r.pg.Exec(ctx, deleteOne, username)
	return err
}

const findOrganization = `
	select organization
	from adm.user_organization
	where username = $1;
`

// FindOrganization находит организацию по username.
func (r *Repo) FindOrganization(ctx context.Context, username string) (string, error) {
	var organization string
	err := r.pg.Query(ctx, &organization, findOrganization, username)
	return organization, err
}
