package model

// NotificationUpdateReq запрос на изменения состояния уведомлений для организации.
type NotificationUpdateReq struct {
	OrganizationID int64  `json:"organization_id" validate:"required,gte=1"`
	Active         bool   `json:"active" validate:"required"`
	Email          string `json:"-"`
	FirstName      string `json:"-"`
	LastName       string `json:"-"`
}
