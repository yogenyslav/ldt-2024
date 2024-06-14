package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// UpdateOne godoc
// @Summary Обновляет избранный предикт.
// @Description Обновляет избранный предикт.
// @Tags favorite
// @Accept json
// @Produce json
// @Param body body model.FavoriteUpdateReq true "Параметры запроса"
// @Success 200 {string} string "Предикт успешно обновлен"
// @Failure 400 {string} string "Некорректный запрос"
// @Router /favorite [put]
func (h *Handler) UpdateOne(c *fiber.Ctx) error {
	var req model.FavoriteUpdateReq
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

	return c.SendStatus(http.StatusOK)
}
