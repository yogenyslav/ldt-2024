package shared

import (
	"time"
)

var (
	UserStateExp = 24 * time.Hour
)

const (
	TraceCtxKey = "traceCtx"
)

type State int8

const (
	_ State = iota
	StateWaitsAuth
)
