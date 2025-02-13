package fund

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stcol316/cushon-isa/internal/models"
)

type FundRepository interface {
	ListFunds(ctx context.Context, page, pageSize int) ([]models.Fund, int, error)
	GetFundByID(ctx context.Context, id string) (*models.Fund, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) listFunds(ctx context.Context, page, pageSize int) ([]models.Fund, int, error) {
	offset := (page - 1) * pageSize

	// First, get total count
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM funds").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Then get paginated data
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name, description 
        FROM funds 
        ORDER BY name 
        LIMIT $1 OFFSET $2
    `, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query funds: %w", err)
	}
	defer rows.Close()

	var funds []models.Fund
	for rows.Next() {
		var fund models.Fund
		if err := rows.Scan(&fund.ID, &fund.Name, &fund.Description); err != nil {
			return nil, 0, fmt.Errorf("failed to scan fund: %w", err)
		}
		funds = append(funds, fund)
	}

	return funds, total, nil
}

func (r *Repository) getFundByID(ctx context.Context, id string) (*models.Fund, error) {
	query := `
	SELECT id, name, description, risk_level_id
	FROM funds
	WHERE id = $1
`
	var fund models.Fund
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&fund.ID,
		&fund.Name,
		&fund.Description,
		&fund.RiskLevel,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("fund not found")
		}
		return nil, fmt.Errorf("failed to get fund: %w", err)
	}

	return &fund, nil
}
