package auth

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/bot/internal/auth/model"
	"github.com/yogenyslav/ldt-2024/bot/internal/state"
	tele "gopkg.in/telebot.v3"
)

// Listen starts the auth server.
func Listen(machine *state.Machine, bot *tele.Bot, port int) *fiber.App {
	srv := fiber.New()
	srv.Use(cors.New(cors.Config{
		AllowOrigins: "https://hawk-handy-wolf.ngrok-free.app",
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
