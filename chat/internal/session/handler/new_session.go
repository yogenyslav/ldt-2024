package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yogenyslav/ldt-2024/chat/internal/session/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// NewSession godoc
// @Summary Новая сессия
// @Description Создать новую сессию в чате для авторизованного пользователя
// @Tags session
// @Accept json
// @Produce json
// @Param token header string true "access token"
// @Success 200 {object} model.NewSessionResp "ID новой сессии"
// @Failure 400 {object} string "Сессия с таким ID уже существует"
// @Router /session/new [post]
func (h *Handler) NewSession(c *fiber.Ctx) error {
	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		return shared.ErrCtxConvertType
	}

	id := uuid.New()
	if err := h.ctrl.NewSession(c.UserContext(), id, username); err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(model.NewSessionResp{
		ID: id,
	})
}
