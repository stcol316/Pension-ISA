package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/stcol316/cushon-isa/internal/models"
)

type Service interface {
	CreateRetailCustomer(ctx context.Context, customer *models.RetailCustomer) error
	GetRetailCustomerByEmail(ctx context.Context, email string) (*models.RetailCustomer, error)
	GetRetailCustomerByID(ctx context.Context, id string) (*models.RetailCustomer, error)
	ListFunds(ctx context.Context) ([]*models.Fund, error)
	GetFundByID(ctx context.Context, id string) (*models.Fund, error)
	CreateInvestment(ctx context.Context, investment *models.Investment) error
	ListInvestmentsByCustomerId(ctx context.Context, id string) ([]*models.Investment, error)
	GetInvestmentByID(ctx context.Context, id string) (*models.Investment, error)
}

type postgresDB struct {
	db *sql.DB
}

var (
	db_name     = "dev_db"
	db_user     = "dev_user"
	db_password = "bEeBwv2JWoFp16Pq1+se3qNGaVoAAgAvKFgBkn5eGeQ="
	port        = "5433"
	host        = "localhost"
)

func NewPostgresDB() (*postgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, host, port, db_name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresDB{db: db}, nil
}

func (s *postgresDB) CreateRetailCustomer(ctx context.Context, customer *models.RetailCustomer) error {
	query := `
	INSERT INTO retail_customers (first_name, last_name, email)
	VALUES ($1, $2, $3)
`

	_, err := s.db.ExecContext(ctx, query,
		customer.FirstName,
		customer.LastName,
		customer.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to create retail customer: %w", err)
	}

	return nil
}

func (s *postgresDB) GetRetailCustomerByEmail(ctx context.Context, email string) (*models.RetailCustomer, error) {
	query := `
	SELECT id, first_name, last_name, email
	FROM retail_customers
	WHERE email = $1
`
	var customer models.RetailCustomer
	err := s.db.QueryRowContext(ctx, query, email).Scan(
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

func (s *postgresDB) GetRetailCustomerByID(ctx context.Context, id string) (*models.RetailCustomer, error) {
	query := `
	SELECT id, first_name, last_name, email
	FROM retail_customers
	WHERE id = $1
`
	var customer models.RetailCustomer
	err := s.db.QueryRowContext(ctx, query, id).Scan(
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

func (s *postgresDB) ListFunds(ctx context.Context) ([]*models.Fund, error) {
	// TODO: Implement funds listing logic
	return nil, nil
}

func (s *postgresDB) GetFundByID(ctx context.Context, id string) (*models.Fund, error) {
	// TODO: Implement fund retrieval by ID logic
	return nil, nil
}

func (s *postgresDB) CreateInvestment(ctx context.Context, investment *models.Investment) error {
	// TODO: Implement investment creation logic
	return nil
}

func (s *postgresDB) ListInvestmentsByCustomerId(ctx context.Context, id string) ([]*models.Investment, error) {
	// TODO: Implement investments listing by customer ID logic
	return nil, nil
}

func (s *postgresDB) GetInvestmentByID(ctx context.Context, id string) (*models.Investment, error) {
	// TODO: Implement investment retrieval by ID logic
	return nil, nil
}
