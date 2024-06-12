package model

import (
	"github.com/yogenyslav/ldt-2024/api/internal/api/pb"
)

// ProductDao is a model for the product entity.
type ProductDao struct {
	RegistryNumberInRk  any     `json:"registry_number_in_rk" bson:"registry_number_in_rk"`
	PaidRub             any     `json:"paid_rub" bson:"paid_rub"`
	RefPrice            any     `json:"ref_price" bson:"ref_price"`
	ConclusionDate      string  `json:"conclusion_date" bson:"conclusion_date"`
	ExecutionTermUntil  string  `json:"execution_term_until" bson:"execution_term_until"`
	ExecutionTermFrom   string  `json:"execution_term_from" bson:"execution_term_from"`
	EndDateOfValidity   string  `json:"end_date_of_validity" bson:"end_date_of_validity"`
	FinalCodeKpgz       string  `json:"final_code_kpgz" bson:"final_code_kpgz"`
	NameSpgz            string  `json:"name_spgz" bson:"name_spgz"`
	ItemNameGk          string  `json:"item_name_gk" bson:"item_name_gk"`
	FinalNameKpgz       string  `json:"final_name_kpgz" bson:"final_name_kpgz"`
	Depth3CodeKpgz      string  `json:"depth3_code_kpgz" bson:"depth3_code_kpgz"`
	NameSte             string  `json:"name_ste" bson:"name_ste"`
	CharacteristicsName string  `json:"characteristics_name" bson:"characteristics_name"`
	ID                  int64   `json:"id" bson:"id"`
	IDSpgz              int64   `json:"id_spgz" bson:"id_spgz"`
	GkPriceRub          float64 `json:"gk_price_rub" bson:"gk_price_rub"`
}

// ToPb converts ProductDao to *pb.Product.
func (p *ProductDao) ToPb() *pb.Product {
	return &pb.Product{
		NameSpgz:       p.NameSpgz,
		FinalNameKpgz:  p.FinalNameKpgz,
		FinalCodeKpgz:  p.FinalCodeKpgz,
		Depth3CodeKpgz: p.Depth3CodeKpgz,
		Id:             p.ID,
	}
}
