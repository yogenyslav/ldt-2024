package model

import "github.com/yogenyslav/ldt-2024/api/internal/shared"

// LoginReq внутренняя модель запроса для метода Login.
type LoginReq struct {
	Username string
	Password string
}

// LoginResp внутренняя модель ответа для метода Login.
type LoginResp struct {
	Token string
	Roles []shared.UserRole
}
