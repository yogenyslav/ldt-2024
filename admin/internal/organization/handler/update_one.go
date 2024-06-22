package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// UpdateOne godoc
// @Summary Изменить название
// @Description Изменить название организации
// @Tags organization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body model.OrganizationUpdateReq true "Новое название"
// @Success 204 "Название изменено"
// @Failure 404 {object} string "Организация не найдена"
// @Router /organization [put]
func (h *Handler) UpdateOne(c *fiber.Ctx) error {
	var req model.OrganizationUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	if err := h.ctrl.UpdateOne(c.UserContext(), req, username); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
