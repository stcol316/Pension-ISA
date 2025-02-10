package fund

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

func (s *Service) listFunds(ctx context.Context) ([]*models.Fund, error) {

	return nil, nil
}

func (s *Service) getFundByID(ctx context.Context, id string) (*models.Fund, error) {
	return s.repo.getFundByID(ctx, id)
}
