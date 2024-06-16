package model

// LoginReq модель запроса для авторизации.
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResp модель ответа для авторизации.
type LoginResp struct {
	Token string   `json:"token"`
	Roles []string `json:"roles"`
}
