package model

import (
	"time"
)

// UserOrganizationDao структура для хранения данных о пользователе и организации.
type UserOrganizationDao struct {
	Username       string
	OrganizationID int64
	IsDeleted      bool
	CreatedAt      time.Time
}
