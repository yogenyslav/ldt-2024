package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// FindOne godoc
// @Summary Получить данные о сессии
// @Description Получить все запросы и ответы для сессии по ID
// @Tags session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "UUID сессии"
// @Success 200 {object} model.FindOneResp "Информация о сессии"
// @Failure 400 {object} string "Неверное значение ID"
// @Router /session/{id} [get]
func (h *Handler) FindOne(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return shared.ErrInvalidUUID
	}

	resp, err := h.ctrl.FindOne(c.UserContext(), id, username)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
