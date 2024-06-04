package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/model"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var (
		err error
		req model.LoginReq
	)

	if err = c.BodyParser(&req); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}
	if err = h.validator.Struct(req); err != nil {
		return errors.Wrap(err, "failed to validate request")
	}

	resp, err := h.ctrl.Login(c.Context(), req)
	if err != nil {
		return errors.Wrap(err, "failed to login")
	}

	return c.Status(http.StatusOK).JSON(resp)
}
