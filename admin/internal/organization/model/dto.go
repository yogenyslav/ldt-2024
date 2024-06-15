package model

import (
	"time"
)

// OrganizationDto организация для передачи внутри сервиса.
type OrganizationDto struct {
	Title     string    `json:"title"`
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// OrganizationCreateReq запрос на создание организации.
type OrganizationCreateReq struct {
	Title string `json:"title"`
}

// OrganizationCreateResp ответ на создание организации.
type OrganizationCreateResp struct {
	ID int64 `json:"id"`
}
