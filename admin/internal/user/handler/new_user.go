package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
	"github.com/yogenyslav/ldt-2024/admin/internal/user/model"
)

// NewUser godoc
// @Summary Создает нового пользователя.
// @Description Создает нового пользователя.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body model.UserCreateReq true "Параметры пользователя"
// @Success 201 {string} string "Пользователь создан"
// @Failure 400 {string} string "Ошибка в запросе"
// @Failure 422 {string} string "Неверный формат данных"
// @Router /user [post]
func (h *Handler) NewUser(c *fiber.Ctx) error {
	var req model.UserCreateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	if err := h.ctrl.NewUser(c.Context(), req); err != nil {
		return err
	}

	for _, role := range req.Roles {
		switch strings.ToLower(role) {
		case shared.RoleAdmin.ToString():
			h.m.AdminUserCount.Inc()
		case shared.RoleAnalyst.ToString():
			h.m.AnalystUserCount.Inc()
		case shared.RoleBuyer.ToString():
			h.m.BuyerUserCount.Inc()
		}
	}

	return c.SendStatus(http.StatusCreated)
}
