package shared

import (
	"time"
)

var (
	UserStateExp = 24 * time.Hour
)

const (
	TraceCtxKey = "traceCtx"
	UserIDKey   = "userId"

	ErrorMessage    = "Что-то пошло не так. Попробуйте еще раз"
	NeedAuthMessage = "Для начала работы с ботом необходимо авторизоваться /auth"
)

type State int8

const (
	_ State = iota
	StateWaitsAuth
)
