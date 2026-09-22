package middlewares

import (
	"net/http"
)

// Middleware es la firma estándar de cualquier middleware
type Middleware func(http.Handler) http.Handler

// Chain aplica múltiples middlewares a un handler en el orden que se especifican
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
