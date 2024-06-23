package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/admin/internal/shared"
)

// ImportData godoc
// @Summary Загрузить данные
// @Description Загрузить данные в архиве
// @Tags organization
// @Security ApiKeyAuth
// @Accept mpfd
// @Produce json
// @Param data formData file true "Архив с данными"
// @Param organization_id formData string true "ID организации"
// @Success 201 {string} string "Данные успешно загружены"
// @Failure 400 {object} string "Ошибка при обработке файлов"
// @Router /organization/import [post]
func (h *Handler) ImportData(c *fiber.Ctx) error {
	mpArchive, err := c.FormFile("data")
	if err != nil {
		log.Error().Err(err).Msg("failed to get form file")
		return shared.ErrParseFormValue
	}

	organizationID, err := strconv.ParseInt(c.FormValue("organization_id"), 10, 64)
	if err != nil {
		log.Error().Msg("organizationID must be int")
		return shared.ErrParseFormValue
	}

	if err := h.ctrl.ImportData(c.UserContext(), mpArchive, organizationID); err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}
