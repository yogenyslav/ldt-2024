package shared

const (
	// UsernameKey ключ для получения имени пользователя из контекста.
	UsernameKey = "x-username"
	// LoginEndpoint эндпоинт для авторизации.
	LoginEndpoint = "/api.AuthService/Login"
)

// UserRole роль пользователя.
type UserRole int8

const (
	_ UserRole = iota
	Admin
	Analyst
	Buyer
)

// ToString возвращает строковое представление роли.
func (r UserRole) ToString() string {
	switch r {
	case Admin:
		return "ADMIN"
	case Analyst:
		return "ANALYST"
	case Buyer:
		return "BUYER"
	default:
		return "UNDEFINED"
	}
}
