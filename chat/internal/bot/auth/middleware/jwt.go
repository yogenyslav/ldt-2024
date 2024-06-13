package middleware

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	tele "gopkg.in/telebot.v3"
)

// JWT валидирует jwt токен.
func JWT(machine *state.Machine, kc *gocloak.GoCloak, realm string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := pkg.GetTraceCtx(c)

			token, err := machine.GetToken(ctx, c.Sender().ID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get token")
				if e := c.Send(shared.NeedAuthMessage); e != nil {
					return e
				}
				return nil
			}

			userInfo, err := kc.GetUserInfo(ctx, token, realm)
			if err != nil || userInfo.PreferredUsername == nil {
				log.Error().Err(err).Msg("failed to get user info")
				if e := c.Send(shared.NeedAuthMessage); e != nil {
					return e
				}
				return nil
			}

			c.Set(shared.UsernameKey, *userInfo.PreferredUsername)
			c.Set(shared.TokenKey, token)

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
