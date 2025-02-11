package middleware

import (
	"context"
	"net/http"
	"strconv"
)

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

type PaginationParams struct {
	Page     int
	PageSize int
}

// Note: Best practice to define our own key to avoid collisions
type contextKey string

const PaginationParamsKey contextKey = "pagination"

// Note: Pagination middleware
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		pageSize := 10

		if p := r.URL.Query().Get("page"); p != "" {
			if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
				page = pageInt
			}
		}

		if size := r.URL.Query().Get("page_size"); size != "" {
			if sizeInt, err := strconv.Atoi(size); err == nil && sizeInt > 0 {
				pageSize = sizeInt
			}
		}

		// Store pagination in context with previously defined key
		ctx := context.WithValue(r.Context(), PaginationParamsKey, PaginationParams{
			Page:     page,
			PageSize: pageSize,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPaginationParams(ctx context.Context) (PaginationParams, bool) {
	params, ok := ctx.Value(PaginationParamsKey).(PaginationParams)
	return params, ok
}
