package shared

import (
	"strings"
	"time"
)

var (
	// UserStateExp время, через которое стейт инвалидируется в redis.
	UserStateExp = 24 * time.Hour
)

const (
	// UsernameKey ключ для получения имени пользователя из контекста.
	UsernameKey = "x-username"
	// TraceCtxKey ключ для trace внутри контекста.
	TraceCtxKey = "traceCtx"
	// UserIDKey ключ для получения userID.
	UserIDKey = "userID"
	// TokenKey ключ для получения токена.
	TokenKey = "token"
	// UsernameKey ключ для получения username.
	enumsUndefined = "UNDEFINED"

	// ErrorMessage сообщение об ошибке внутри бота.
	ErrorMessage = "Что-то пошло не так. Попробуйте еще раз"
	// NeedAuthMessage сообщение о необходимости авторизоваться.
	NeedAuthMessage = "Для начала работы с ботом необходимо авторизоваться /auth"
)

// ResponseStatus статус ответа.
type ResponseStatus int8

const (
	_ ResponseStatus = iota
	// StatusCreated статус создания.
	StatusCreated
	// StatusProcessing статус обработки.
	StatusProcessing
	// StatusSuccess статус успешного завершения.
	StatusSuccess
	// StatusError статус ошибки.
	StatusError
	// StatusCanceled статус отмены.
	StatusCanceled
)

// ToString возвращает строковое представление статуса ответа.
func (s ResponseStatus) ToString() string {
	switch s {
	case StatusCreated:
		return "CREATED"
	case StatusProcessing:
		return "PROCESSING"
	case StatusSuccess:
		return "SUCCCESS"
	case StatusError:
		return "ERROR"
	case StatusCanceled:
		return "CANCELED"
	}
	return enumsUndefined
}

// QueryType тип запроса.
type QueryType int8

const (
	_ QueryType = iota
	// TypePrediction предсказание.
	TypePrediction
	// TypeStock акция.
	TypeStock
)

// ToString возвращает строковое представление типа запроса.
func (t QueryType) ToString() string {
	switch t {
	case TypePrediction:
		return "PREDICTION"
	case TypeStock:
		return "STOCK"
	}
	return enumsUndefined
}

// QueryCommand команда запроса.
type QueryCommand string

const (
	// CommandValid команда валидации.
	CommandValid QueryCommand = "valid"
	// CommandInvalid команда невалидации.
	CommandInvalid QueryCommand = "invalid"
	// CommandCancel команда отмены.
	CommandCancel QueryCommand = "cancel"
)

// QueryStatus статус запроса.
type QueryStatus int8

const (
	_ QueryStatus = iota
	// StatusPending статус ожидания.
	StatusPending
	// StatusValid статус валидности.
	StatusValid
	// StatusInvalid статус невалидного запроса.
	StatusInvalid
)

// ToString возвращает строковое представление статуса запроса.
func (s QueryStatus) ToString() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusValid:
		return "VALID"
	case StatusInvalid:
		return "INVALID"
	}
	return enumsUndefined
}

// UserRole роль пользователя.
type UserRole int8

const (
	RoleUndefined UserRole = iota
	// RoleAdmin администратор.
	RoleAdmin
	// RoleAnalyst аналитик.
	RoleAnalyst
	// RoleBuyer закупщик.
	RoleBuyer
)

// ToString возвращает строковое представление роли пользователя.
func (r UserRole) ToString() string {
	switch r {
	case RoleAdmin:
		return "ADMIN"
	case RoleAnalyst:
		return "ANALYST"
	case RoleBuyer:
		return "BUYER"
	default:
		return enumsUndefined
	}
}

// RoleFromString возвращает роль по ее строковому представлению.
func RoleFromString(v string) UserRole {
	switch strings.Split(strings.ToLower(v), "_")[1] {
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

// UserState тип для состояний пользователя.
type UserState int8

const (
	_ UserState = iota
	// StateWaitsAuth ожидает авторизации.
	StateWaitsAuth
	// StatePending пользователь авторизован, ожидание действия.
	StatePending
	// StateValidate ожидает подтверждения результата prompter.
	StateValidate
	// StateHint ожидает подсказки.
	StateHint
)
