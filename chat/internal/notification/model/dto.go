package model

// NotificationUpdateReq запрос на изменения состояния уведомлений для организации.
type NotificationUpdateReq struct {
	OrganizationID int64  `json:"organization_id" validate:"required,gte=1"`
	Active         bool   `json:"active"`
	Email          string `json:"-"`
	FirstName      string `json:"-"`
	LastName       string `json:"-"`
}

// NotificationExistsResp ответ на запрос о наличии уведомления.
type NotificationExistsResp struct {
	Exists bool `json:"exists"`
}
