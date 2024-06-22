package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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
// @Router /user/{organization_id} [get]
func (h *Handler) List(c *fiber.Ctx) error {
	organizationID, err := c.ParamsInt("organization_id")
	if err != nil {
		log.Error().Err(err).Msg("organizationID must be int")
		return err
	}

	users, err := h.ctrl.List(c.Context(), int64(organizationID))
	if err != nil {
		return err
	}

	return c.JSON(users)
}
