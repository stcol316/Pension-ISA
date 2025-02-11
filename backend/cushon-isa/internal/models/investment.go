package models

import "time"

type Investment struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	FundID     string    `json:"fundId"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"` // TODO: We probably want something to confirm status of investments here
}

type CreateInvestmentRequest struct {
	CustomerID string  `json:"customerId"`
	FundID     string  `json:"fundId"`
	Amount     float64 `json:"amount"`
}

func NewInvestment(customerId, fundId string, amount float64) Investment {
	return Investment{
		CustomerID: customerId,
		FundID:     fundId,
		Amount:     amount,
	}
}
