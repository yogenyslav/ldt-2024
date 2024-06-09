package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Rename godoc
// @Summary Обновить заголовок
// @Description Обновить заголовок сессии, который отображается в интерфейсе
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "access token"
// @Param req body model.RenameReq true "ID и новое название сессии"
// @Success 204 {object} string "Сессия переименована"
// @Failure 400 {object} string "Сессия с таким ID уже существует"
// @Failure 404 {object} string "Сессия с таким ID не найдена"
// @Router /session/rename [put]
func (h *Handler) Rename(c *fiber.Ctx) error {
	var req model.RenameReq

	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	if err := h.ctrl.Rename(c.UserContext(), req); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
