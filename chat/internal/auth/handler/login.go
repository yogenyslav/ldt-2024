package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Login godoc
// @Summary Авторизация
// @Description Авторизоваться в системе, используя почту и пароль
// @Tags auth
// @Accept json
// @Produce json
// @Param req body model.LoginReq true "Запрос на авторизацию"
// @Success 200 {object} model.LoginResp "Успешная авторизация"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 401 {object} string "Неверный данные для входа"
// @Failure 422 {object} string "Ошибка валидации данных"
// @Router /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req model.LoginReq

	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	resp, err := h.ctrl.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}
