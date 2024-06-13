package middleware

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	tele "gopkg.in/telebot.v3"
)

// State проверяет, что у пользователя верный state.
func State(machine *state.Machine) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := pkg.GetTraceCtx(c)

			s, err := machine.GetState(ctx, c.Sender().ID)
			if err != nil {
				return err
			}
			if s != shared.StatePending && s != shared.StateValidate {
				return nil
			}
			return next(c)
		}
	}
}
