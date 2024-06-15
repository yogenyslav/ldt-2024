package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// FindOne godoc
// @Summary Получить организацию
// @Description Получить организацию для пользователя
// @Tags organization
// @Accept json
// @Produce json
// @Param Authorization header string true "access token"
// @Success 200 {object} model.OrganizationDto "Информация об организации"
// @Failure 404 {object} string "Организация не найдена"
// @Router /organization [get]
func (h *Handler) FindOne(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	resp, err := h.ctrl.FindOne(c.UserContext(), username)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
