package investment

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	isaerrors "github.com/stcol316/cushon-isa/internal/errors"
	"github.com/stcol316/cushon-isa/internal/models"
)

// TODO: Use interfaces at service level instead of "repo *Repository"
type InvestmentRepository interface {
	CreateInvestment(ctx context.Context, investment *models.Investment) error
	ListInvestmentsByCustomerID(ctx context.Context, id string, page, pageSize int) ([]models.Investment, int, error)
	GetInvestmentByID(ctx context.Context, id string) (*models.Investment, error)
	GetCustomerFundTotal(ctx context.Context, customerID, fundID string) (*models.InvestmentSummary, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) createInvestment(ctx context.Context, investment *models.Investment) error {
	// Note: Limit customers to one fund. Remove to allow multiple
	var existingFundID *string
	exerr := r.db.QueryRowContext(ctx, `
        SELECT DISTINCT fund_id 
        FROM investments 
        WHERE customer_id = $1 
        LIMIT 1
    `, investment.CustomerID).Scan(&existingFundID)

	if exerr != nil && exerr != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing investments: %w", exerr)
	}

	// If customer has existing investments, ensure it's the same fund
	if exerr != sql.ErrNoRows && *existingFundID != investment.FundID {
		return isaerrors.ErrDifferentFundNotAllowed
	}

	// Note: We begin a transaction that will rollback our actions if either insert or materialized view refresh fails
	// This ensures data consistency between the investments table and view
	tx, txerr := r.db.BeginTx(ctx, nil)
	if txerr != nil {
		return fmt.Errorf("failed to begin transaction: %w", txerr)
	}
	defer tx.Rollback()

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

	log.Printf("Attempting to refresh materialized view")
	// Note: Refresh materialized view
	_, err = tx.ExecContext(ctx, "REFRESH MATERIALIZED VIEW customer_fund_totals")
	if err != nil {
		return fmt.Errorf("failed to refresh materialized view: %w", err)
	}

	log.Printf("Successfully refreshed materialized view")
	// Everything worked as expected so we commit
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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

	// Scan rows into our investments slice
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

// Note: This is fetching data from the materialized view
func (r *Repository) getCustomerFundTotal(ctx context.Context, customerID, fundID string) (*models.InvestmentSummary, error) {
	query := `
        SELECT 
            customer_id,
            first_name,
            last_name,
            email,
            fund_id,
            fund_name,
            total_investment
        FROM customer_fund_totals
        WHERE customer_id = $1 AND fund_id = $2`

	var summary models.InvestmentSummary
	err := r.db.QueryRowContext(ctx, query, customerID, fundID).Scan(
		&summary.CustomerID,
		&summary.FirstName,
		&summary.LastName,
		&summary.Email,
		&summary.FundID,
		&summary.FundName,
		&summary.TotalInvestment,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying investment summary: %w", err)
	}

	return &summary, nil
}
