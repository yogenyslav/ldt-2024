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
	select uo.username, ue.email, exists(
		select 1
		from chat.notification n
		where n.email = ue.email and n.organization_id = $1
	) as notifications
	from adm.user_organization uo
	join adm.user_email ue 
	    on uo.username = ue.username
	where 
	    organization_id = $1 and is_deleted = false;
`

// List возвращает список пользователей по организации.
func (r *Repo) List(ctx context.Context, organizationID int64) ([]model.UserListResp, error) {
	var users []model.UserListResp
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

const checkUserOrganization = `
	select count(*)
	from adm.user_organization
	where username = $1 and organization_id = $2;
`

// CheckUserOrganization проверить, состоит ли пользователь в организации.
func (r *Repo) CheckUserOrganization(ctx context.Context, username string, organizationID int64) (bool, error) {
	var exists int
	err := r.pg.Query(ctx, &exists, checkUserOrganization, username, organizationID)
	return exists == 1, err
}

const insertEmailByUsername = `
	insert into adm.user_email(email, username)
	values ($1, $2);
`

// InsertEmailByUsername записывает email по username.
func (r *Repo) InsertEmailByUsername(ctx context.Context, email, username string) error {
	_, err := r.pg.Exec(ctx, insertEmailByUsername, email, username)
	return err
}
