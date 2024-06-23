package repo

import (
	"context"

	"github.com/yogenyslav/ldt-2024/chat/internal/user/model"
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
	insert into chat.notification (email, first_name, last_name, organization_id)
	values ($1, $2, $3, $4);
`

// InsertOne добавляет новое уведомление.
func (r *Repo) InsertOne(ctx context.Context, params model.NotificationDao) error {
	_, err := r.pg.Exec(ctx, insertOne, params.Email, params.FirstName, params.LastName, params.OrganizationID)
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
