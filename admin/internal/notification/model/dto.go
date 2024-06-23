package model

// NotificationUpdateReq запрос на изменения состояния уведомлений для организации.
type NotificationUpdateReq struct {
	Username       string `json:"username" validate:"required"`
	Active         bool   `json:"active"`
	OrganizationID int64  `json:"organization_id" validate:"required,gte=1"`
}
