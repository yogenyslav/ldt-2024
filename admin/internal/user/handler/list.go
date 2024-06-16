package handler

import (
	"github.com/gofiber/fiber/v2"
)

// List godoc
// @Summary Возвращает список пользователей по организации.
// @Description Возвращает список пользователей по организации.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organization path string true "Название организации"
// @Success 200 {array} string "Список пользователей"
// @Router /user/{organization} [get]
func (h *Handler) List(c *fiber.Ctx) error {
	organization := c.Params("organization")
	users, err := h.ctrl.List(c.Context(), organization)
	if err != nil {
		return err
	}

	return c.JSON(users)
}
