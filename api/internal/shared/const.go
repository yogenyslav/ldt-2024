package shared

import "strings"

const (
	// UsernameKey ключ для получения имени пользователя из контекста.
	UsernameKey = "x-username"
	// LoginEndpoint эндпоинт для авторизации.
	LoginEndpoint = "/api.AuthService/Login"
)

// UserRole роль пользователя.
type UserRole int8

const (
	RoleUndefined UserRole = iota
	RoleAdmin
	RoleAnalyst
	RoleBuyer
)

// RoleFromString конвертирует строку роли в число.
func RoleFromString(v string) UserRole {
	switch strings.ToLower(v) {
	case "admin":
		return RoleAdmin
	case "analyst":
		return RoleAnalyst
	case "buyer":
		return RoleBuyer
	default:
		return RoleUndefined
	}
}
