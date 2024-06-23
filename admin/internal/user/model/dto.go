package model

// UserCreateReq запрос на создание пользователя.
type UserCreateReq struct {
	Roles           []string `json:"roles" validate:"required"`
	Username        string   `json:"username" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	FirstName       string   `json:"first_name" validate:"required"`
	LastName        string   `json:"last_name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	OrganizationIDs []int64  `json:"organization_ids"`
}

// UserListReq запрос на получение списка пользователей по организации.
type UserListReq struct {
	OrganizationID int64 `json:"organization_id" validate:"required,gte=1"`
}

// UserUpdateOrganizationReq запрос на добавление/удаление пользователя в/из организацию(-и).
type UserUpdateOrganizationReq struct {
	Username       string `json:"username" validate:"required"`
	OrganizationID int64  `json:"organization_id" validate:"required,gte=1"`
}

// UserListResp ответ на запрос списка пользователей по организации.
type UserListResp struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Notifications bool   `json:"notifications"`
}
