package shared

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
	"strings"
)

const (
	// UsernameKey ключ для получения имени пользователя из контекста.
	UsernameKey = "x-username"
	// LoginEndpoint эндпоинт для авторизации.
	LoginEndpoint = "/api.AuthService/Login"
)

// RoleFromString конвертирует строку роли в число.
func RoleFromString(v string) pb.UserRole {
	switch strings.ToLower(v) {
	case "admin":
		return pb.UserRole_ROLE_ADMIN
	case "analyst":
		return pb.UserRole_ROLE_ANALYST
	case "buyer":
		return pb.UserRole_ROLE_BUYER
	default:
		return pb.UserRole_ROLE_UNDEFINED
	}
}
