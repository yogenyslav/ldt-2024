package repo

import (
	"context"

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
