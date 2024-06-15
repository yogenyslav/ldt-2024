package shared

import (
	"errors"
)

// 400
var (
	// ErrParseBody ошибка парсинга тела запроса.
	ErrParseBody = errors.New("failed to parse body")
)

// 401
var (
	// ErrMissingJWT ошибка отсутствия JWT токена.
	ErrMissingJWT = errors.New("missing JWT token")
	// ErrInvalidJWT ошибка невалидного JWT токена.
	ErrInvalidJWT = errors.New("invalid JWT token")
	// ErrLoginFailed ошибка авторизации.
	ErrLoginFailed = errors.New("login failed")
)

// 403
var (
	// ErrForbidden ошибка доступа.
	ErrForbidden = errors.New("forbidden")
)

// 404
var (
	// ErrNoOrganization ошибка отсутствия организации.
	ErrNoOrganization = errors.New("organization not found")
)

// 500
var (
	// ErrCtxConvertType ошибка конвертации значения контекста в тип.
	ErrCtxConvertType = errors.New("failed to convert context value to type")
	// ErrCipherTooShort ошибка слишком короткого шифра.
	ErrCipherTooShort = errors.New("cipher too short")
	// ErrEncryption ошибка шифрования.
	ErrEncryption = errors.New("failed to encrypt")

	// ErrCreateOrganization ошибка создания организации.
	ErrCreateOrganization = errors.New("failed to create organization")
	// ErrUpdateOrganization ошибка обновления организации.
	ErrUpdateOrganization = errors.New("failed to update organization")
	// ErrGetOrganization ошибка получения организации.
	ErrGetOrganization = errors.New("failed to get organization")
)
