package fund

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	mw "github.com/stcol316/cushon-isa/internal/middleware"
	helper "github.com/stcol316/cushon-isa/pkg/helpers"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListFundsHandler(w http.ResponseWriter, r *http.Request) {
	params, ok := mw.GetPaginationParams(r.Context())
	if !ok {
		log.Printf("Failed to get pagination params from context")
		helper.RespondWithError(w, http.StatusInternalServerError, "pagination error")
		return
	}

	fmt.Printf("Pagination Params: Page=%d, PageSize=%d\n", params.Page, params.PageSize)

	result, err := h.service.listFunds(r.Context(), params.Page, params.PageSize)
	if err != nil {
		h.handleFundError(w, err)
		return
	}
	helper.RespondWithJSON(w, http.StatusOK, result)

}

func (h *Handler) GetFundByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "customer ID is required")
		return
	}

	customer, err := h.service.getFundByID(r.Context(), id)
	if err != nil {
		h.handleFundError(w, err)
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, customer)
}

func (h *Handler) handleFundError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		helper.RespondWithError(w, http.StatusNotFound, "no funds found")
	default:
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
