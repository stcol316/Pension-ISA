package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

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

// Note: DB healthcheck go routine
func (p *PostgresDB) StartHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			stats := p.HealthCheck()
			// We might want to capture these metrics elsewhere but we'll print for now
			log.Println(stats)
			if stats["status"] == "down" {
				log.Println("Warning: Database health degraded")
				// Could implement alert/notification system here
				continue
			}
		}
	}()
}

// Checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (p *PostgresDB) HealthCheck() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := p.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := p.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *PostgresDB) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}
