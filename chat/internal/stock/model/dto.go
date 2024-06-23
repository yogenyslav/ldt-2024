package model

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
