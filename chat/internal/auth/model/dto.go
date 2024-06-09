package model

// LoginReq is the request model for the login endpoint.
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResp is the response model for the login endpoint.
type LoginResp struct {
	Token string `json:"token"`
}
