package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// List godoc
// @Summary Список сессий
// @Description Получить список сессий в порядке убывания момента создания от последней к первой для авторизованного пользователя
// @Tags session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.ListResp "Список сессий"
// @Router /session/list [get]
func (h *Handler) List(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	sessions, err := h.ctrl.List(c.UserContext(), username)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(model.ListResp{
		Sessions: sessions,
	})
}
