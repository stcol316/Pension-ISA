package investment

import (
	"context"
	"database/sql"

	"github.com/stcol316/cushon-isa/internal/models"
)

type FundRepository interface {
	listFunds(ctx context.Context) ([]*models.Fund, error)
	getFundByID(ctx context.Context, id string) (*models.Fund, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) createInvestment(ctx context.Context, investment *models.Investment) error {
	// TODO: Implement investment creation logic
	return nil
}

func (r *Repository) listInvestmentsByCustomerID(ctx context.Context, id string) ([]*models.Investment, error) {
	// TODO: Implement investments listing by customer ID logic
	return nil, nil
}

func (r *Repository) getInvestmentByID(ctx context.Context, id string) (*models.Investment, error) {
	// TODO: Implement investment retrieval by ID logic
	return nil, nil
}
