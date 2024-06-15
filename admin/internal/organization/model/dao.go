package model

import (
	"time"
)

// OrganizationDao организация в базе данных.
type OrganizationDao struct {
	Username  string
	Title     string
	S3Bucket  string
	ID        int64
	CreatedAt time.Time
}

// ToDto конвертирует OrganizationDao в OrganizationDto.
func (o OrganizationDao) ToDto() OrganizationDto {
	return OrganizationDto{
		Title:     o.Title,
		ID:        o.ID,
		CreatedAt: o.CreatedAt,
	}
}
