package investment

import (
	"context"
	"database/sql"
	"fmt"

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
	query := `
	INSERT INTO investments (customer_id, fund_id, amount)
	VALUES ($1, $2, $3)
`

	_, err := r.db.ExecContext(ctx, query,
		investment.CustomerID,
		investment.FundID,
		investment.Amount,
	)
	if err != nil {
		return fmt.Errorf("failed to make investment: %w", err)
	}

	return nil
}

func (r *Repository) listInvestmentsByCustomerID(ctx context.Context, id string, page, pageSize int) ([]models.Investment, int, error) {
	offset := (page - 1) * pageSize

	// First, get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM investments").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Then get paginated data
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, customer_id, fund_id, amount, created_at 
        FROM investments
		WHERE customer_id = $1
        ORDER BY created_at
        LIMIT $2 OFFSET $3
    `, id, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query investments: %w", err)
	}
	defer rows.Close()

	var investments []models.Investment
	for rows.Next() {
		var investment models.Investment
		if err := rows.Scan(&investment.ID,
			&investment.CustomerID,
			&investment.FundID,
			&investment.Amount,
			&investment.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan investment: %w", err)
		}
		investments = append(investments, investment)
	}

	return investments, total, nil

}

func (r *Repository) getInvestmentByID(ctx context.Context, id string) (*models.Investment, error) {
	query := `
	SELECT id, customer_id, fund_id, amount
	FROM investments
	WHERE id = $1
`
	var investment models.Investment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&investment.ID,
		&investment.CustomerID,
		&investment.FundID,
		&investment.Amount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("investment not found")
		}
		return nil, fmt.Errorf("failed to get investment: %w", err)
	}

	return &investment, nil
}
