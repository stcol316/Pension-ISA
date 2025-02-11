package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stcol316/cushon-isa/internal/models"
	"github.com/stcol316/cushon-isa/pkg/helpers"
	helper "github.com/stcol316/cushon-isa/pkg/helpers"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRetailCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Error handling and validation could be improved here.
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		helpers.RespondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	// We could probably split this out into a common decode and validate helper function
	req := new(models.CreateRetailCustomerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.handleCustomerError(w, err)
		return
	}

	customer, err := h.service.createRetailCustomer(r.Context(), req)
	if err != nil {
		h.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *Handler) GetRetailCustomerByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid customer ID format")
		return
	}

	customer, err := h.service.getRetailCustomerByID(r.Context(), id)
	if err != nil {
		h.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *Handler) GetRetailCustomerByEmailHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "email")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer email is required")
		return
	}

	customer, err := h.service.getRetailCustomerByEmail(r.Context(), id)
	if err != nil {
		h.handleCustomerError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *Handler) handleCustomerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		helper.RespondWithError(w, http.StatusNotFound, "customer not found")
	default:
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
