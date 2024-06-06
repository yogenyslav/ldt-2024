package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Delete godoc
// @Summary Удалить сессию
// @Description Удалить сессию по ID
// @Tags session
// @Accept json
// @Produce json
// @Param token header string true "access token"
// @Param id path uuid.UUID true "ID сессии"
// @Success 204 {object} string "Сессия удалена"
// @Failure 400 {object} string "Неверное значение ID"
// @Failure 404 {object} string "Сессия с таким ID не найдена"
// @Router /session/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return shared.ErrInvalidUUID
	}

	if err := h.ctrl.Delete(c.UserContext(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
