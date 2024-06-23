package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/notification/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// Switch godoc
// @Summary Включает или выключает уведомления.
// @Description Включает или выключает уведомления.
// @Tags notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.NotificationUpdateReq true "Параметры запроса"
// @Success 200 {string} string "Статус уведомлений изменен"
// @Failure 400 {object} string "Некорректный запрос"
// @Router /notification/switch [post]
func (h *Handler) Switch(c *fiber.Ctx) error {
	var req model.NotificationUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	email, ok := c.Locals(shared.EmailKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}
	firstName, ok := c.Locals(shared.FirstNameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}
	lastName, ok := c.Locals(shared.LastNameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	req.Email = email
	req.FirstName = firstName
	req.LastName = lastName

	if err := h.ctrl.Switch(c.UserContext(), req); err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
