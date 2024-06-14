package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// DeleteOne godoc
// @Summary Удаляет избранный предикт по QueryID.
// @Description Удаляет избранный предикт по QueryID.
// @Tags favorite
// @Accept json
// @Produce json
// @Param id path int true "QueryID предикта"
// @Success 204 {string} string "Предикт успешно удален"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Предикт не найден"
// @Router /favorite/{id} [delete]
func (h *Handler) DeleteOne(c *fiber.Ctx) error {
	queryID, err := c.ParamsInt("id")
	if err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	if err := h.ctrl.DeleteOne(c.UserContext(), int64(queryID), username); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
