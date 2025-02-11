package fund

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

func (s *Service) listFunds(ctx context.Context, page, pageSize int) (*mw.PaginatedResult, error) {

	funds, total, err := s.repo.listFunds(ctx, page, pageSize)
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

func (s *Service) getFundByID(ctx context.Context, id string) (*models.Fund, error) {
	return s.repo.getFundByID(ctx, id)
}
