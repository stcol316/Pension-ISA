package investment

import (
	"context"

	"github.com/stcol316/cushon-isa/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createInvestment(ctx context.Context) (*models.Investment, error) {

	return nil, nil
}

func (s *Service) listInvestmentsByCustomerID(ctx context.Context, id string) ([]*models.Investment, error) {
	return s.repo.listInvestmentsByCustomerID(ctx, id)
}

func (s *Service) getInvestmentByID(ctx context.Context, id string) (*models.Investment, error) {
	return s.repo.getInvestmentByID(ctx, id)
}
