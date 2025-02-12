package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/stcol316/cushon-isa/internal/config"
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

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	var (
		db_name     = cfg.DBName
		db_user     = cfg.DBUser
		db_password = cfg.DBPassword
		port        = cfg.DBPort
		host        = cfg.DBHost
	)

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
