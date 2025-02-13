package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/stcol316/cushon-isa/internal/config"
	"github.com/stcol316/cushon-isa/internal/customer"
	"github.com/stcol316/cushon-isa/internal/database"
	"github.com/stcol316/cushon-isa/internal/fund"
	"github.com/stcol316/cushon-isa/internal/investment"
	"github.com/stcol316/cushon-isa/internal/server"
)

func main() {

	fmt.Println("Loading config...")

	cfg, cfgerr := config.Load()
	if cfgerr != nil {
		log.Fatalf("Failed to load config: %v", cfgerr)
	}

	//Note: Easily swappable database configuration
	db_service, dberr := database.NewPostgresDB(cfg)
	if dberr != nil {
		log.Fatal(dberr)
	}
	fmt.Println("Starting Healthcheck go routine")
	db_service.StartHealthCheck(1 * time.Minute)

	// Note: Repository Pattern. Handles data access
	fmt.Println("Creating Repository Layer")
	customerRepo := customer.NewRepository(db_service.DB())
	fundRepo := fund.NewRepository(db_service.DB())
	investmentRepo := investment.NewRepository(db_service.DB())

	// Note: Service layer to handle business logic between DB and handlers
	fmt.Println("Creating Service Layer")
	customerService := customer.NewService(customerRepo)
	fundService := fund.NewService(fundRepo)
	investmentService := investment.NewService(investmentRepo)

	// Note: Presentation layer to handle APIs
	fmt.Println("Creating Presentation Layer")
	customerHandler := customer.NewHandler(customerService)
	fundHandler := fund.NewHandler(fundService)
	investmentHandler := investment.NewHandler(investmentService)

	server := server.NewServer(cfg, customerHandler, fundHandler, investmentHandler)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, db_service, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")

}

// Note: Graceful shutdown
func gracefulShutdown(apiServer *http.Server, db_service *database.PostgresDB, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	if err := db_service.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
