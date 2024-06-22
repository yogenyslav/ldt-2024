package model

import (
	"time"
)

// OrganizationDao организация в базе данных.
type OrganizationDao struct {
	CreatedAt time.Time
	Username  string
	Title     string
	ID        int64
}

// ToDto конвертирует OrganizationDao в OrganizationDto.
func (o OrganizationDao) ToDto() OrganizationDto {
	return OrganizationDto{
		Title:     o.Title,
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
	}
}
