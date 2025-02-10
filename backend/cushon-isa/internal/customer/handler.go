package customer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stcol316/cushon-isa/internal/models"
	helper "github.com/stcol316/cushon-isa/pkg/helpers"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRetailCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Error handling could be improved here.
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
