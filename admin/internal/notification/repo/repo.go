package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/admin/internal/notification/model"
	"github.com/yogenyslav/pkg/storage"
)

// Repo репозиторий для уведомлений.
type Repo struct {
	pg storage.SQLDatabase
}

// New создает новый Repo.
func New(pg storage.SQLDatabase) *Repo {
	return &Repo{pg: pg}
}

const insertOne = `
	insert into chat.notification (email, organization_id)
	values ($1, $2);
`

// InsertOne добавляет новое уведомление.
func (r *Repo) InsertOne(ctx context.Context, params model.NotificationDao) error {
	_, err := r.pg.Exec(ctx, insertOne, params.Email, params.OrganizationID)
	return err
}

const deleteOne = `
	delete from chat.notification 
	where email = $1;
`

// DeleteOne удаляет уведомление.
func (r *Repo) DeleteOne(ctx context.Context, email string) error {
	_, err := r.pg.Exec(ctx, deleteOne, email)
	return err
}

const checkNotification = `
	select exists (
		select 1
		from chat.notification
		where email = $1 and organization_id = $2
	);
`

// CheckNotification проверяет наличие уведомления.
func (r *Repo) CheckNotification(ctx context.Context, params model.NotificationDao) (bool, error) {
	var exists bool
	err := r.pg.Query(ctx, &exists, checkNotification, params.Email, params.OrganizationID)
	return exists, err
}

const getEmailByUsername = `
	select email 
	from adm.user_email
	where username = $1;
`

// GetEmailByUsername возвращает email по username.
func (r *Repo) GetEmailByUsername(ctx context.Context, username string) (string, error) {
	var email string
	err := r.pg.Query(ctx, &email, getEmailByUsername, username)
	return email, err
}
