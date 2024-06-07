package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
)

type SessionDaoBuilder struct {
	SessionDao *model.SessionDao
}

func SessionDao() *SessionDaoBuilder {
	return &SessionDaoBuilder{SessionDao: &model.SessionDao{}}
}

func (b *SessionDaoBuilder) New() *SessionDaoBuilder {
	return b.
		ID(uuid.New()).
		Title("title").
		Username("user").
		CreatedAt(time.Now())
}

func (b *SessionDaoBuilder) V() model.SessionDao {
	return *b.SessionDao
}

func (b *SessionDaoBuilder) P() *model.SessionDao {
	return b.SessionDao
}

func (b *SessionDaoBuilder) ID(v uuid.UUID) *SessionDaoBuilder {
	b.SessionDao.ID = v
	return b
}

func (b *SessionDaoBuilder) Title(v string) *SessionDaoBuilder {
	b.SessionDao.Title = v
	return b
}

func (b *SessionDaoBuilder) Username(v string) *SessionDaoBuilder {
	b.SessionDao.Username = v
	return b
}

func (b *SessionDaoBuilder) CreatedAt(v time.Time) *SessionDaoBuilder {
	b.SessionDao.CreatedAt = v
	return b
}
