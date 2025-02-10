package fund

import (
	"database/sql"
	"errors"
	"net/http"

	helper "github.com/stcol316/cushon-isa/pkg/helpers"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListFundsHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetFundByIdHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleFundError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		helper.RespondWithError(w, http.StatusNotFound, "no funds found")
	default:
		helper.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
