package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// List godoc
// @Summary Возвращает список избранных предиктов.
// @Description Возвращает список избранных предиктов.
// @Tags favorite
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} model.FavoriteDto "Список избранных предиктов"
// @Router /favorite/list [get]
func (h *Handler) List(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	resp, err := h.ctrl.List(c.UserContext(), username)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
