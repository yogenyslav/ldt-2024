package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// FindOne godoc
// @Summary Возвращает избранный предикт по ID.
// @Description Возвращает избранный предикт по ID.
// @Tags favorite
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID предикта"
// @Success 200 {object} model.FavoriteDto "Избранный предикт"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Предикт не найден"
// @Router /favorite/{id} [get]
func (h *Handler) FindOne(c *fiber.Ctx) error {
	queryID, err := c.ParamsInt("id")
	if err != nil {
		return shared.ErrParseBody
	}

	resp, err := h.ctrl.FindOne(c.UserContext(), int64(queryID))
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
