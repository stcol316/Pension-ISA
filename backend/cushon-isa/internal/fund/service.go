package fund

import (
	"context"
	"fmt"

	"github.com/stcol316/cushon-isa/internal/models"
)

type Service struct {
	repo *Repository
}

type PaginatedResult struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int  `json:"current_page"`
		PageSize    int  `json:"page_size"`
		TotalItems  int  `json:"total_items"`
		TotalPages  int  `json:"total_pages"`
		HasNext     bool `json:"has_next"`
		HasPrevious bool `json:"has_previous"`
	} `json:"pagination"`
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) listFunds(ctx context.Context, page int, pageSize int) (*PaginatedResult, error) {

	funds, total, err := s.repo.listFunds(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list funds: %w", err)
	}

	// Calculate pagination metadata
	totalPages := (total + pageSize - 1) / pageSize

	result := &PaginatedResult{
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
