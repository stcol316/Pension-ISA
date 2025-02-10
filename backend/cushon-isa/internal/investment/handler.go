package investment

import (
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateInvestmentHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetInvestmentByIDHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListCustomerInvestmentsHandler(w http.ResponseWriter, r *http.Request) {

}
