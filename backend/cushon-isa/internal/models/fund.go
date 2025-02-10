package models

type Fund struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	MinInvestment float64 `json:"minInvestment"`
	IsActive      bool    `json:"isActive"`
}
