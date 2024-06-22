package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// FindOne godoc
// @Summary Получить список организаций
// @Description Получить список организаций для пользователя
// @Tags organization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} model.OrganizationDto "Список организаций пользователя"
// @Failure 404 {object} string "Организация не найдена"
// @Router /organization [get]
func (h *Handler) FindOne(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	resp, err := h.ctrl.ListForUser(c.UserContext(), username)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
