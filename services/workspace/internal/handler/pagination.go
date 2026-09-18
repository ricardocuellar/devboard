package handler

import (
	"net/http"
	"strconv"
)

// PaginationParams reperesenta parámetros de paginación
type PaginationParams struct {
	Cursor string
	Limit  int
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

// ParsePagination lee cursor y limit de la query string
func ParsePagination(request *http.Request) PaginationParams {
	cursor := request.URL.Query().Get("cursor")
	limit := defaultLimit

	if limitQuery := request.URL.Query().Get("limit"); limitQuery != "" {
		if parsed, err := strconv.Atoi(limitQuery); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	return PaginationParams{Cursor: cursor, Limit: limit}
}

// PaginatedResponse es el formato estándar para cualquier listado paginado
type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
