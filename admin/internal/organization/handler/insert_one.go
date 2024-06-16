package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/organization/model"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// InsertOne godoc
// @Summary Создает новую организацию.
// @Description Создает новую организацию.
// @Tags organization
// @Accept json
// @Produce json
// @Param Authorization header string true "access token"
// @Param body body model.OrganizationCreateReq true "Параметры создания организации"
// @Success 201 {object} model.OrganizationCreateResp "ID созданной организации"
// @Failure 400 {object} string "Неверные параметры запроса"
// @Router /organization [post]
func (h *Handler) InsertOne(c *fiber.Ctx) error {
	var req model.OrganizationCreateReq
	if err := c.BodyParser(&req); err != nil {
		return shared.ErrParseBody
	}

	username, ok := c.Locals(shared.UsernameKey).(string)
	if !ok {
		log.Error().Msg("failed to get username from context")
		return shared.ErrCtxConvertType
	}

	resp, err := h.ctrl.InsertOne(c.UserContext(), req, username)
	if err != nil {
		return err
	}

	h.m.NumberOfActivatedCompanies.Inc()

	return c.Status(http.StatusCreated).JSON(resp)
}
