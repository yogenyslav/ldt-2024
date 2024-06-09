package middleware

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/bot/internal/shared"
	"github.com/yogenyslav/ldt-2024/bot/internal/state"
	"github.com/yogenyslav/ldt-2024/bot/pkg"
	tele "gopkg.in/telebot.v3"
)

// JWT middleware for authorizing users.
func JWT(machine *state.Machine, kc *gocloak.GoCloak, realm string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := pkg.GetTraceCtx(c)

			token, err := machine.GetToken(ctx, c.Sender().ID)
			if err != nil {
				if e := c.Send(shared.NeedAuthMessage); e != nil {
					return e
				}
				return nil
			}

			userInfo, err := kc.GetUserInfo(ctx, token, realm)
			if err != nil || userInfo.PreferredUsername == nil {
				if e := c.Send(shared.NeedAuthMessage); e != nil {
					return e
				}
				return nil
			}

			if c.Callback() != nil {
				log.Debug().Any("data", c.Callback().Data).Msg("callback")
			}

			err = next(c)
			if err == nil {
				log.Info().Int("updateId", c.Update().ID).Int64("userId", c.Sender().ID).Msg("success")
			}
			return err
		}
	}
}
