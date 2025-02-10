package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

// Note: Best practice to define our own key to avoid collisions
type contextKey string

const paginationParamsKey contextKey = "pagination"

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
		ctx := context.WithValue(r.Context(), paginationParamsKey, PaginationParams{
			Page:     page,
			PageSize: pageSize,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
