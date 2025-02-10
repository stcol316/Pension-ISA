package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	mw "github.com/stcol316/cushon-isa/internal/api/middleware"
	"github.com/stcol316/cushon-isa/internal/models"
	helper "github.com/stcol316/cushon-isa/pkg/helpers"
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
			r.Post("/retail", s.createRetailCustomerHandler)
			r.Get("/id/{id}", s.getRetailCustomerByIdHandler)
			r.Get("/email/{email}", s.getRetailCustomerByEmailHandler)
		})

		// Fund routes
		r.Route("/funds", func(r chi.Router) {
			r.With(mw.Paginate).Get("/", s.listFundsHandler)
			r.Get("/{id}", s.getFundHandler)
		})

		// Investment routes
		r.Route("/investments", func(r chi.Router) {
			r.Post("/", s.createInvestmentHandler)
			r.Get("/{id}", s.getInvestmentHandler)
			r.With(mw.Paginate).Get("/customer/{customerId}", s.ListCustomerInvestmentsHandler)
		})
	})

	// r.Get("/health", s.healthHandler)

	return r
}

// Retail Customer handlers
func (s *Server) createRetailCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Error handling could be improved here.
	// We could probably split this out into a common decode and validate helper function
	req := new(models.CreateRetailCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.handleCustomerError(w, err)
		return
	}

	customer := models.NewRetailCustomer(req.FirstName, req.LastName, req.Email)
	if err := s.db.CreateRetailCustomer(r.Context(), &customer); err != nil {
		s.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (s *Server) getRetailCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer ID is required")
		return
	}

	customer, err := s.db.GetRetailCustomerByID(r.Context(), id)
	if err != nil {
		s.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)

}

func (s *Server) getRetailCustomerByEmailHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "email")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer email is required")
		return
	}

	customer, err := s.db.GetRetailCustomerByEmail(r.Context(), id)
	if err != nil {
		s.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (s *Server) handleCustomerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		helper.RespondWithError(w, http.StatusNotFound, "customer not found")
	default:
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

// Fund Handlers
func (s *Server) listFundsHandler(w http.ResponseWriter, r *http.Request) {
	// // Get pagination from context
	// params := r.Context().Value("pagination").(PaginationParams)

	// // Calculate offset
	// offset := (params.Page - 1) * params.PageSize

	// // Use in your database query
	// query := `
	// 		SELECT id, name, description
	// 		FROM funds
	// 		LIMIT $1 OFFSET $2
	// 	`

	// funds, err := s.db.QueryContext(r.Context(), query, params.PageSize, offset)
	// if err != nil {
	// 	WriteJSON(w, http.StatusInternalServerError, map[string]string{
	// 		"error": "Failed to fetch funds",
	// 	})
	// 	return
	// }

	// // Return paginated results
	// WriteJSON(w, http.StatusOK, map[string]interface{}{
	// 	"funds": funds,
	// 	"pagination": map[string]int{
	// 		"page":      params.Page,
	// 		"page_size": params.PageSize,
	// 	},
	// })
}

func (s *Server) getFundHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Fund retrieved successfully"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// Investment Handlers
func (s *Server) createInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Investment created successfully"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) getInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Investment retrieved successfully"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) ListCustomerInvestmentsHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Customer investments retrieved successfully"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
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
