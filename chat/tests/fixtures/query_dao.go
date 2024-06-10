package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

var queryID int64 = 0

type QueryDaoBuilder struct {
	QueryDao *model.QueryDao
}

func QueryDao() *QueryDaoBuilder {
	return &QueryDaoBuilder{QueryDao: &model.QueryDao{}}
}

func (b *QueryDaoBuilder) New() *QueryDaoBuilder {
	defer func() {
		queryID++
	}()
	return b.
		ID(queryID).
		SessionID(uuid.New()).
		Username("user").
		Prompt("some prompt").
		CreatedAt(time.Now()).
		Status(shared.StatusPending).
		Type(shared.TypePrediction).
		Product("test product")
}

func (b *QueryDaoBuilder) V() model.QueryDao {
	return *b.QueryDao
}

func (b *QueryDaoBuilder) P() *model.QueryDao {
	return b.QueryDao
}

func (b *QueryDaoBuilder) ID(v int64) *QueryDaoBuilder {
	b.QueryDao.ID = v
	return b
}

func (b *QueryDaoBuilder) SessionID(v uuid.UUID) *QueryDaoBuilder {
	b.QueryDao.SessionID = v
	return b
}

func (b *QueryDaoBuilder) Username(v string) *QueryDaoBuilder {
	b.QueryDao.Username = v
	return b
}

func (b *QueryDaoBuilder) Prompt(v string) *QueryDaoBuilder {
	b.QueryDao.Prompt = v
	return b
}

func (b *QueryDaoBuilder) CreatedAt(v time.Time) *QueryDaoBuilder {
	b.QueryDao.CreatedAt = v
	return b
}

func (b *QueryDaoBuilder) Status(v shared.QueryStatus) *QueryDaoBuilder {
	b.QueryDao.Status = v
	return b
}

func (b *QueryDaoBuilder) Type(v shared.QueryType) *QueryDaoBuilder {
	b.QueryDao.Type = v
	return b
}

func (b *QueryDaoBuilder) Product(v string) *QueryDaoBuilder {
	b.QueryDao.Product = v
	return b
}
