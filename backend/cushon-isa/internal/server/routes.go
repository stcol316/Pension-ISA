package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/stcol316/cushon-isa/internal/models"
)

// Note: We use pagination for our GET List calls
type PaginationParams struct {
	Page     int
	PageSize int
}

// Note: Best practice to define our own key to avoid collisions
type contextKey string

const paginationParamsKey contextKey = "pagination"

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
	r.Route("/customers", func(r chi.Router) {
		r.Post("/retail", s.createRetailCustomerHandler)
		r.Get("/{id}", s.getRetailCustomerHandler)
	})

	// Fund routes
	r.Route("/funds", func(r chi.Router) {
		r.With(paginate).Get("/", s.listFundsHandler)
		r.Get("/{id}", s.getFundHandler)
	})

	// Investment routes
	r.Route("/investments", func(r chi.Router) {
		r.Post("/", s.createInvestmentHandler)
		r.Get("/{id}", s.getInvestmentHandler)
		r.With(paginate).Get("/customer/{customerId}", s.ListCustomerInvestmentsHandler)
	})

	// r.Get("/health", s.healthHandler)

	return r
}

// Note: Pagination middleware
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		pageSize := 10

		if p := r.URL.Query().Get("page"); p != "" {
			if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
				page = pageInt
			}
		}

		if size := r.URL.Query().Get("page_size"); size != "" {
			if sizeInt, err := strconv.Atoi(size); err == nil && sizeInt > 0 {
				pageSize = sizeInt
			}
		}

		// Store pagination in context with previously defined key
		ctx := context.WithValue(r.Context(), paginationParamsKey, PaginationParams{
			Page:     page,
			PageSize: pageSize,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the status code
	w.WriteHeader(status)

	// Encode and write the response body
	return json.NewEncoder(w).Encode(v)
}

// Retail Customer handlers
func (s *Server) createRetailCustomerHandler(w http.ResponseWriter, r *http.Request) {
	//
	createReq := new(models.CreateRetailCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(createReq); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	rc := models.NewRetailCustomer(createReq.FirstName, createReq.LastName, createReq.Email)

	if err := s.db.CreateRetailCustomer(r.Context(), &rc); err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create customer",
		})
		return
	}

	WriteJSON(w, http.StatusOK, rc)
}

func (s *Server) getRetailCustomerHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Customer retrieved successfully"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// Fund Handlers
func (s *Server) listFundsHandler(w http.ResponseWriter, r *http.Request) {

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
