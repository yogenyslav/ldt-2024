package shared

import (
	"github.com/pkg/errors"
)

// 400
var (
	// ErrLoginFailed ошибка при неудачной попытке входа.
	ErrLoginFailed = errors.New("login failed")
	// ErrParseBody ошибка при неудачной попытке парсинга тела запроса.
	ErrParseBody = errors.New("failed to parse body")
	// ErrInvalidUUID ошибка при неверном значении UUID.
	ErrInvalidUUID = errors.New("invalid uuid")
	// ErrWsProtocolRequired ошибка при необходимости обновления до websocket.
	ErrWsProtocolRequired = errors.New("upgrade to websocket is required")
	// ErrEmptyQueryHint ошибка при пустом query hint.
	ErrEmptyQueryHint = errors.New("query hint can't be empty")
)

// 401
var (
	// ErrMissingJWT ошибка при отсутствии JWT.
	ErrMissingJWT = errors.New("missing JWT")
	// ErrInvalidJWT ошибка при неверном JWT.
	ErrInvalidJWT = errors.New("invalid JWT")
)

// 404
var (
	// ErrNoSessionWithID ошибка при отсутствии сессии с запрошенным id.
	ErrNoSessionWithID = errors.New("no session with such id found")
	// ErrNoResponseWithID ошибка при отсутствии ответа с запрошенным id.
	ErrNoResponseWithID = errors.New("no response with such id found")
	// ErrNoQueryWithID ошибка при отсутствии запроса с запрошенным id.
	ErrNoQueryWithID = errors.New("no query with such id found")
)

// 500
var (
	// ErrCtxConvertType ошибка при неверном типе конвертации в контексте.
	ErrCtxConvertType = errors.New("wrong type converting in context")
	// ErrCipherTooShort ошибка при слишком коротком ключе шифрования.
	ErrCipherTooShort = errors.New("cipher text is too short")
	// ErrEncryption ошибка при неудачном шифровании.
	ErrEncryption = errors.New("encryption failed")
	// ErrBeginTx ошибка при неудачной попытке начать sql транзакцию.
	ErrBeginTx = errors.New("failed to being an sql transaction")
	// ErrCommitTx ошибка при неудачной попытке завершить sql транзакцию.
	ErrCommitTx = errors.New("failed to commit transaction")

	// ErrCreateSession ошибка при неудачной попытке создать сессию.
	ErrCreateSession = errors.New("failed to create session")
	// ErrSessionDuplicateID ошибка при дублировании id сессии.
	ErrSessionDuplicateID = errors.New("got duplicated session id")
	// ErrGetSession ошибка при неудачной попытке получить сессию.
	ErrGetSession = errors.New("failed to get session")
	// ErrUpdateSession ошибка при неудачной попытке обновить сессию.
	ErrUpdateSession = errors.New("failed to update session")
	// ErrDeleteSession ошибка при неудачной попытке удалить сессию.
	ErrDeleteSession = errors.New("failed to delete session")

	// ErrCreateQuery ошибка при неудачной попытке создать запрос.
	ErrCreateQuery = errors.New("failed to create query")
	// ErrCreateResponse ошибка при неудачной попытке создать ответ.
	ErrCreateResponse = errors.New("failed to create response")
	// ErrUpdateResponse ошибка при неудачной попытке обновить ответ.
	ErrUpdateResponse = errors.New("failed to update response")
	// ErrUpdateQuery ошибка при неудачной попытке обновить запрос.
	ErrUpdateQuery = errors.New("failed to update query")
	// ErrGetQuery ошибка при неудачной попытке получить запрос.
	ErrGetQuery = errors.New("failed to get query")
)
