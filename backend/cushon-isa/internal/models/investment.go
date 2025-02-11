package models

import (
	"time"
)

type Investment struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	FundID     string    `json:"fundId"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"` // TODO: We probably want something to confirm status of investments here
}

type InvestmentSummary struct {
	CustomerID      string  `json:"customer_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	FundID          string  `json:"fund_id"`
	FundName        string  `json:"fund_name"`
	TotalInvestment float64 `json:"total_investment"`
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
