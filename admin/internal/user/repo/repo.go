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
	insert into adm.user_organization(username, organization_id)
	values ($1, $2);
`

// InsertOrganization добавляет пользователя в организацию.
func (r *Repo) InsertOrganization(ctx context.Context, params model.UserOrganizationDao) error {
	_, err := r.pg.Exec(ctx, insertOrganization, params.Username, params.OrganizationID)
	return err
}

const list = `
	select username
	from adm.user_organization
	where organization_id = $1 and is_deleted = false;
`

// List возвращает список пользователей по организации.
func (r *Repo) List(ctx context.Context, organizationID int64) ([]string, error) {
	var users []string
	err := r.pg.QuerySlice(ctx, &users, list, organizationID)
	return users, err
}

const deleteOne = `
	delete from adm.user_organization
	where username = $1 and organization_id = $2;
`

// DeleteOrganization удаляет пользователя из организации.
func (r *Repo) DeleteOrganization(ctx context.Context, params model.UserOrganizationDao) error {
	_, err := r.pg.Exec(ctx, deleteOne, params.Username, params.OrganizationID)
	return err
}
