package model

// LoginReq is the internal request model for the Login method.
type LoginReq struct {
	Email    string
	Password string
}

// LoginResp is the internal response model for the Login method.
type LoginResp struct {
	Token string
}
