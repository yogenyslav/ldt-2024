package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/auth/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	tele "gopkg.in/telebot.v3"
)

// Listen запуск сервера авторизации.
func Listen(machine *state.Machine, bot *tele.Bot, port int) *fiber.App {
	srv := fiber.New()
	srv.Use(cors.New(cors.Config{
		AllowOrigins: "https://hawk-handy-wolf.ngrok-free.app,http://localhost:5173",
	}))
	srv.Post("/bot/auth", func(c *fiber.Ctx) error {
		var req model.AuthorizeReq
		if err := c.BodyParser(&req); err != nil {
			log.Error().Err(err).Msg("failed to parse the request")
			return err
		}
		if err := machine.SetToken(c.UserContext(), req.TgID, req.Token); err != nil {
			log.Error().Err(err).Msg("failed to set the token")
			return err
		}
		if err := machine.SetRolesFromStrings(c.UserContext(), req.TgID, req.Roles); err != nil {
			log.Error().Err(err).Msg("failed to set roles")
			return err
		}
		if err := machine.SetState(c.UserContext(), req.TgID, shared.StatePending); err != nil {
			log.Error().Err(err).Msg("failed to set state")
			return err
		}
		if _, err := bot.Send(&tele.User{ID: req.TgID}, "Вы успешно авторизовались!"); err != nil {
			log.Error().Err(err).Msg("failed to send a message")
			return err
		}
		return c.SendStatus(http.StatusOK)
	})

	go func() {
		if err := srv.Listen(fmt.Sprintf(":%d", port)); err != nil {
			log.Error().Err(err).Msg("failed to start the auth server")
		}
	}()

	return srv
}
