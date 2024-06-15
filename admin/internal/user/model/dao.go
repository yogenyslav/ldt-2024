package model

import (
	"time"
)

// UserOrganizationDao структура для хранения данных о пользователе и организации.
type UserOrganizationDao struct {
	Username     string
	Organization string
	IsDeleted    bool
	CreatedAt    time.Time
}
