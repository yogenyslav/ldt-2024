package handler

import (
	"net/http"

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
// @Success 201 {string} string "Данные успешно загружены"
// @Failure 400 {object} string "Ошибка при обработке файлов"
// @Router /organization/import [post]
func (h *Handler) ImportData(c *fiber.Ctx) error {
	mpArchive, err := c.FormFile("data")
	if err != nil {
		log.Error().Err(err).Msg("failed to get form file")
		return shared.ErrParseFormFile
	}

	org, ok := c.Locals(shared.OrganizationKey).(string)
	if !ok {
		log.Error().Err(err).Msg("failed to get organization")
		return shared.ErrCtxConvertType
	}

	if err := h.ctrl.ImportData(c.UserContext(), mpArchive, org); err != nil {
		return err
	}

	return c.SendStatus(http.StatusCreated)
}
