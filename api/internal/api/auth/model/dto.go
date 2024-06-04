package model

type LoginReq struct {
	Email    string
	Password string
}

type LoginResp struct {
	Token string
}
