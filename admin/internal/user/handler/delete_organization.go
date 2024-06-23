package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
)

// DeleteOrganization godoc
// @Summary Удаляет организацию.
// @Description Удаляет организацию.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.UserUpdateOrganizationReq true "Параметры на удаление"
// @Success 204 {string} string "Организация удалена"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /user [delete]
func (h *Handler) DeleteOrganization(c *fiber.Ctx) error {
	var req model.UserUpdateOrganizationReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	if req.Username == username {
		return errors.New("can't delete yourself")
	}

	if err := h.ctrl.DeleteOrganization(c.Context(), req); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
