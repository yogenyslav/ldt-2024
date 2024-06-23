package model

// NotificationUpdateReq запрос на изменения состояния уведомлений для организации.
type NotificationUpdateReq struct {
	Email          string `json:"email"`
	Active         bool   `json:"active"`
	OrganizationID int64  `json:"organization_id" validate:"required,gte=1"`
}
