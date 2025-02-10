package customer

import (
	"context"
	"fmt"

	"github.com/stcol316/cushon-isa/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createRetailCustomer(ctx context.Context, req *models.CreateRetailCustomerRequest) (*models.RetailCustomer, error) {
	customer := models.NewRetailCustomer(req.FirstName, req.LastName, req.Email)
	if err := s.repo.createRetailCustomer(ctx, &customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &customer, nil
}

func (s *Service) getRetailCustomerByID(ctx context.Context, id string) (*models.RetailCustomer, error) {
	return s.repo.getRetailCustomerByID(ctx, id)
}

func (s *Service) getRetailCustomerByEmail(ctx context.Context, email string) (*models.RetailCustomer, error) {
	return s.repo.getRetailCustomerByEmail(ctx, email)
}
