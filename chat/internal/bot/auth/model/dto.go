package model

// AuthorizeReq модель запроса на авторизацию в боте.
type AuthorizeReq struct {
	Token string   `json:"token"`
	Roles []string `json:"roles"`
	TgID  int64    `json:"tg_id"`
}
