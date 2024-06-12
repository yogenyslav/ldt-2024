package model

import "github.com/yogenyslav/ldt-2024/api/internal/shared"

// LoginReq is the internal request model for the Login method.
type LoginReq struct {
	Username string
	Password string
}

// LoginResp is the internal response model for the Login method.
type LoginResp struct {
	Token string
	Roles []shared.UserRole
}
