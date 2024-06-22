package model

// UniqueCodesReq модель запроса на получение уникальный кодов товаров.
type UniqueCodesReq struct {
	OrganizationID int64 `json:"organization_id" validate:"required,ge=1"`
}

// UniqueCodesResp модель ответа на запрос получения регулярных товаров.
type UniqueCodesResp struct {
	Codes []UniqueCodeDto `json:"codes"`
}

// UniqueCodeDto представление уникального товара с типом регулярности для передачи внутри сервиса.
type UniqueCodeDto struct {
	Segment string `json:"segment"`
	Name    string `json:"name"`
	Regular bool   `json:"regular"`
}
