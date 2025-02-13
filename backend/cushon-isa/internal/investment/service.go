package investment

import (
	"context"
	"fmt"

	mw "github.com/stcol316/cushon-isa/internal/middleware"
	"github.com/stcol316/cushon-isa/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) createInvestment(ctx context.Context, req *models.CreateInvestmentRequest) (*models.Investment, error) {
	investment := models.NewInvestment(req.CustomerID, req.FundID, req.Amount)
	if err := s.repo.createInvestment(ctx, &investment); err != nil {
		return nil, fmt.Errorf("failed to make investment: %w", err)
	}

	return &investment, nil
}

func (s *Service) listInvestmentsByCustomerID(ctx context.Context, id string, page, pageSize int) (*mw.PaginatedResult, error) {
	funds, total, err := s.repo.listInvestmentsByCustomerID(ctx, id, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list funds: %w", err)
	}

	// Calculate pagination metadata
	totalPages := (total + pageSize - 1) / pageSize

	result := &mw.PaginatedResult{
		Data: funds,
	}
	result.Pagination.CurrentPage = page
	result.Pagination.PageSize = pageSize
	result.Pagination.TotalItems = total
	result.Pagination.TotalPages = totalPages
	result.Pagination.HasNext = page < totalPages
	result.Pagination.HasPrevious = page > 1

	return result, nil
}

func (s *Service) getInvestmentByID(ctx context.Context, id string) (*models.Investment, error) {
	return s.repo.getInvestmentByID(ctx, id)
}

func (s *Service) getCustomerFundTotal(ctx context.Context, customer_id, fund_id string) (*models.InvestmentSummary, error) {
	return s.repo.getCustomerFundTotal(ctx, customer_id, fund_id)
}
