package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// DeleteOrganization godoc
// @Summary Удаляет организацию.
// @Description Удаляет организацию.
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "access token"
// @Param username path string true "Имя пользователя"
// @Success 204 {string} string "Организация удалена"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /user/{username} [delete]
func (h *Handler) DeleteOrganization(c *fiber.Ctx) error {
	username := c.Params("username")
	if err := h.ctrl.DeleteOrganization(c.Context(), username); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
