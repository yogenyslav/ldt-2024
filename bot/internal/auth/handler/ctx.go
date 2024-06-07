package handler

import (
	"context"

	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	tele "gopkg.in/telebot.v3"
)

// GetTraceCtx returns the context.Context from tele.Context.
func GetTraceCtx(c tele.Context) context.Context {
	return c.Get(shared.TraceCtxKey).(context.Context)
}
