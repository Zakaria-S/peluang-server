package dto

type PriceRangeRequest struct {
	Min int `json:"min" validate:"required,number"`
	Max int `json:"max" validate:"required,number,gtfield=Min"`
}

type PriceRangeResponse struct {
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Slug string `json:"slug"`
}
