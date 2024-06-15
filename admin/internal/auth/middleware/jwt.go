package middleware

import (
	"context"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/ldt-2024/admin/pkg"
	"github.com/yogenyslav/ldt-2024/admin/pkg/secure"
)

type UserOrganizationRepo interface {
	FindOrganization(ctx context.Context, username string) (string, error)
}

// JWT валидирует jwt токен.
func JWT(kc *gocloak.GoCloak, realm, cipher string, repo UserOrganizationRepo) fiber.Handler {
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

		userInfo, err := kc.GetUserInfo(c.UserContext(), authToken, realm)
		if err != nil || userInfo.PreferredUsername == nil {
			return shared.ErrInvalidJWT
		}

		userID := *userInfo.Sub

		groups, err := kc.GetUserGroups(c.UserContext(), authToken, realm, userID, gocloak.GetGroupsParams{})
		if err != nil {
			return err
		}

		hasAccess := false
		for _, role := range groups {
			if *role.Name == strings.ToLower(shared.RoleAdmin.ToString()) {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return shared.ErrForbidden
		}

		org, err := repo.FindOrganization(c.UserContext(), *userInfo.PreferredUsername)
		if err != nil {
			return err
		}

		c.Locals(shared.UsernameKey, *userInfo.PreferredUsername)
		c.Locals(shared.OrganizationKey, org)
		c.SetUserContext(pkg.PushToken(c.UserContext(), authToken))
		return c.Next()
	}
}
