package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/api/pb"
	"github.com/yogenyslav/ldt-2024/chat/internal/stock/model"
)

// UniqueCodes godoc
// @Summary Регулярные товары
// @Description Получить набор уникальных записей с разделением на регулярные и нерегулярные товары
// @Tags stock
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.UniqueCodesResp "Список с товарами"
// @Router /stock/unique_codes [get]
func (h *Handler) UniqueCodes(c *fiber.Ctx) error {
	uniqueCodes, err := h.predictor.UniqueCodes(c.UserContext(), &pb.UniqueCodesReq{})
	if err != nil {
		return err
	}
	codes := uniqueCodes.GetCodes()
	resp := make([]model.UniqueCodeDto, len(codes))
	for i := 0; i < len(codes); i++ {
		resp[i] = model.UniqueCodeDto{
			Segment: codes[i].GetSegment(),
			Name:    codes[i].GetName(),
			Regular: codes[i].GetRegular(),
		}
	}
	return c.Status(http.StatusOK).JSON(model.UniqueCodesResp{
		Codes: resp,
	})
}
