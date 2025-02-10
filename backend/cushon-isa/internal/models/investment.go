package models

import "time"

type Investment struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId"`
	FundID     string    `json:"fundId"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"` // We probably want something to confirm status of investments here
}
