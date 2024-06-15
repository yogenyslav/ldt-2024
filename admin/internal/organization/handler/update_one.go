package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// UpdateOne godoc
// @Summary Обновить организацию
// @Description Обновить организацию
// @Tags organization
// @Accept json
// @Produce json
// @Param Authorization header string true "access token"
// @Param body body model.OrganizationUpdateReq true "Параметры обновления организации"
// @Success 204 "Организация обновлена"
// @Failure 400 {object} string "Неверные параметры запроса"
// @Router /organization [put]
func (h *Handler) UpdateOne(c *fiber.Ctx) error {
	var req model.OrganizationUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		log.Error().Msg("failed to get username from context")
		return shared.ErrCtxConvertType
	}

	if err := h.ctrl.UpdateOne(c.UserContext(), req, username); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
