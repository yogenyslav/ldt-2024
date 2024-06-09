package middleware

import (
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg/secure"
)

func JWT(kc *gocloak.GoCloak, realm, cipher string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) != 2 {
			return shared.ErrMissingJWT
		}

		authToken, err := secure.Decrypt(t[1], cipher)
		if err != nil {
			return err
		}

		userInfo, err := kc.GetUserInfo(c.Context(), authToken, realm)
		if err != nil || userInfo.PreferredUsername == nil {
			return shared.ErrInvalidJWT
		}

		c.Locals(shared.UsernameKey, *userInfo.PreferredUsername)
		return c.Next()
	}
}
