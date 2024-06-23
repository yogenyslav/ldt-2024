package model

// NotificationDao модель подписки на уведомления.
type NotificationDao struct {
	Email          string
	FirstName      string
	LastName       string
	OrganizationID int64
}
