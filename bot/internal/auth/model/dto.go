package model

// AuthorizeReq is the request model for the authorize endpoint.
type AuthorizeReq struct {
	Token string `json:"token"`
	TgID  int64  `json:"tg_id"`
}
