package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
)

// InsertOrganization godoc
// @Summary Добавляет пользователя в организацию.
// @Description Добавляет пользователя в организацию.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body model.UserUpdateOrganizationReq true "Параметры пользователя"
// @Success 200 {string} string "Пользователь добавлен в организацию"
// @Failure 400 {string} string "Ошибка в запросе"
// @Failure 422 {string} string "Неверный формат данных"
// @Router /user/organization [post]
func (h *Handler) InsertOrganization(c *fiber.Ctx) error {
	var req model.UserUpdateOrganizationReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	if err := h.ctrl.InsertOrganization(c.Context(), req); err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
