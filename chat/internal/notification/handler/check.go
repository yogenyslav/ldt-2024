package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/notification/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Check godoc
// @Summary Проверяет наличие уведомлений.
// @Description Проверяет наличие уведомлений.
// @Tags notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organization_id path int true "ID организации"
// @Success 200 {object} model.NotificationExistsResp "Статус уведомления"
// @Failure 400 {object} string "Некорректный запрос"
// @Router /notification/check/{organization_id} [get]
func (h *Handler) Check(c *fiber.Ctx) error {
	organizationID, err := c.ParamsInt("organization_id")
	if err != nil {
		log.Error().Err(err).Msg("failed to parse organization_id")
		return err
	}

	email, ok := c.Locals("email").(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	exists, err := h.ctrl.Check(c.UserContext(), email, int64(organizationID))
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(model.NotificationExistsResp{Exists: exists})
}
