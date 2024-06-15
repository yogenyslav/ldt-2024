package shared

const (
	// UsernameKey ключ для локального хранилища имени пользователя.
	UsernameKey = "username"
	// TraceCtxKey ключ для локального хранилища контекста с trace.
	TraceCtxKey = "trace"
)

// UserRole роль пользователя.
type UserRole int8

const (
	_ UserRole = iota
	// RoleAdmin администратор.
	RoleAdmin
	// RoleAnalyst аналитик.
	RoleAnalyst
	// RoleBuyer закупщик.
	RoleBuyer
)

// ToString возвращает строковое представление роли.
func (r UserRole) ToString() string {
	switch r {
	case RoleAdmin:
		return "ADMIN"
	case RoleAnalyst:
		return "ANALYST"
	case RoleBuyer:
		return "BUYER"
	default:
		return "UNDEFINED"
	}
}
