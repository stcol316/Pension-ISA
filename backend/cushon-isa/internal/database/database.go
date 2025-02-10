package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/stcol316/cushon-isa/internal/models"
)

type Service interface {
	CreateInvestment(ctx context.Context, investment *models.Investment) error
	ListInvestmentsByCustomerId(ctx context.Context, id string) ([]*models.Investment, error)
	GetInvestmentByID(ctx context.Context, id string) (*models.Investment, error)
}

type PostgresDB struct {
	db *sql.DB
}

var (
	db_name     = "dev_db"
	db_user     = "dev_user"
	db_password = "bEeBwv2JWoFp16Pq1+se3qNGaVoAAgAvKFgBkn5eGeQ="
	port        = "5433"
	host        = "localhost"
)

func NewPostgresDB() (*PostgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, host, port, db_name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) DB() *sql.DB {
	return p.db
}
