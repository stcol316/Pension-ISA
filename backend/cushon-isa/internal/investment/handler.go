package investment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	mw "github.com/stcol316/cushon-isa/internal/middleware"
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

func (h *Handler) CreateInvestmentHandler(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		helpers.RespondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	req := new(models.CreateInvestmentRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.handleInvestmentError(w, err)
		return
	}

	investment, err := h.service.createInvestment(r.Context(), req)
	if err != nil {
		h.handleInvestmentError(w, err)
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, investment)
}

func (h *Handler) GetInvestmentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "investment ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid investment ID format")
		return
	}

	investment, err := h.service.getInvestmentByID(r.Context(), id)
	if err != nil {
		h.handleInvestmentError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, investment)
}

func (h *Handler) ListCustomerInvestmentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListCustomerInvestmentsHandler")
	id := chi.URLParam(r, "customerId")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer ID is required")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid customer ID format")
		return
	}

	params, ok := mw.GetPaginationParams(r.Context())
	if !ok {
		log.Printf("Failed to get pagination params from context")
		helper.RespondWithError(w, http.StatusInternalServerError, "pagination error")
		return
	}

	fmt.Printf("Pagination Params: Page=%d, PageSize=%d\n", params.Page, params.PageSize)

	result, err := h.service.listInvestmentsByCustomerID(r.Context(), id, params.Page, params.PageSize)
	if err != nil {
		h.handleInvestmentError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, result)
}

func (h *Handler) GetCustomerFundTotalHandler(w http.ResponseWriter, r *http.Request) {
	customer_id := chi.URLParam(r, "customerId")
	fund_id := chi.URLParam(r, "fundId")
	if customer_id == "" || fund_id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer ID and fund ID required")
		return
	}

	if _, err := uuid.Parse(customer_id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid customer ID format")
		return
	}

	if _, err := uuid.Parse(fund_id); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "invalid fund ID format")
		return
	}

	investment, err := h.service.getCustomerFundTotal(r.Context(), customer_id, fund_id)
	if err != nil {
		h.handleInvestmentError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, investment)
}

func (h *Handler) handleInvestmentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		helper.RespondWithError(w, http.StatusNotFound, "no investments found")
	default:
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
