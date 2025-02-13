package customer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stcol316/cushon-isa/internal/models"
)

type CustomerRepository interface {
	CreateRetailCustomer(ctx context.Context, customer *models.RetailCustomer) error
	GetRetailCustomerByEmail(ctx context.Context, email string) (*models.RetailCustomer, error)
	GetRetailCustomerByID(ctx context.Context, id string) (*models.RetailCustomer, error)
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) createRetailCustomer(ctx context.Context, customer *models.RetailCustomer) error {
	query := `
	INSERT INTO retail_customers (first_name, last_name, email)
	VALUES ($1, $2, $3)
`

	_, err := r.db.ExecContext(ctx, query,
		customer.FirstName,
		customer.LastName,
		customer.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to create retail customer: %w", err)
	}

	return nil
}

func (r *Repository) getRetailCustomerByEmail(ctx context.Context, email string) (*models.RetailCustomer, error) {
	query := `
	SELECT id, first_name, last_name, email
	FROM retail_customers
	WHERE email = $1
`
	var customer models.RetailCustomer
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

func (r *Repository) getRetailCustomerByID(ctx context.Context, id string) (*models.RetailCustomer, error) {
	query := `
	SELECT id, first_name, last_name, email
	FROM retail_customers
	WHERE id = $1
`
	var customer models.RetailCustomer
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}
