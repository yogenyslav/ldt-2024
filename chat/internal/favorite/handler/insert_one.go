package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/favorite/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// InsertOne godoc
// @Summary Добавляет новый предикт в избранное.
// @Description Добавляет новый предикт в избранное.
// @Tags favorite
// @Accept json
// @Produce json
// @Param body body model.FavoriteCreateReq true "Параметры запроса"
// @Success 201 {string} string "Предикт успешно добавлен в избранное"
// @Failure 400 {string} string "Некорректный запрос"
// @Router /favorite [post]
func (h *Handler) InsertOne(c *fiber.Ctx) error {
	var req model.FavoriteCreateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	if err := h.ctrl.InsertOne(c.UserContext(), req, username); err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}
