package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	mw "github.com/stcol316/cushon-isa/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)                    // Log HTTP requests
	r.Use(middleware.RealIP)                    // Extracts real client IP when behind a proxy
	r.Use(middleware.Recoverer)                 // Recovers from panics and ensure durability
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	r.Use(cors.Handler(cors.Options{
		// Note: We would want to restrict origins in production to avoid things like CSRF and DDOS,
		// e.g. "https://cushon.co.uk", "https://api.cushon.co.uk"
		// This is fine for dev
		// TODO: Probably switch to env vars
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Test route
	r.Get("/", s.HelloWorldHandler)

	// Customer routes
	//Note: API versioning
	//TODO: Split into separate services to facilitate microservice architecture
	r.Route("/v1", func(r chi.Router) {
		r.Route("/customers", func(r chi.Router) {
			r.Post("/retail", s.customerHandler.CreateRetailCustomerHandler)
			r.Get("/id/{id}", s.customerHandler.GetRetailCustomerByIdHandler)
			r.Get("/email/{email}", s.customerHandler.GetRetailCustomerByEmailHandler)
		})

		// Fund routes
		r.Route("/funds", func(r chi.Router) {
			// Note: We use pagination for our GET List calls
			r.With(mw.Paginate).Get("/", s.fundHandler.ListFundsHandler)
			r.Get("/id/{id}", s.fundHandler.GetFundByIdHandler)
		})

		// Investment routes
		r.Route("/investments", func(r chi.Router) {
			r.Post("/", s.investmentHandler.CreateInvestmentHandler)
			r.Get("/id/{id}", s.investmentHandler.GetInvestmentByIDHandler)
			r.With(mw.Paginate).Get("/customer/{customerId}", s.investmentHandler.ListCustomerInvestmentsHandler)
		})
	})

	// r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
