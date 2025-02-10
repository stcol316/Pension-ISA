package fund

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

func (r *Repository) listFunds(ctx context.Context) ([]*models.Fund, error) {
	// TODO: Implement funds listing logic
	return nil, nil
}

func (r *Repository) getFundByID(ctx context.Context, id string) (*models.Fund, error) {
	// TODO: Implement fund retrieval by ID logic
	return nil, nil
}
