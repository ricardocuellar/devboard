package middlewares

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"time"
)

const RequestIDHeader = "X-Request-ID"

func RequestID(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			requestID := rand.Text()
			start := time.Now()

			request.Header.Set(RequestIDHeader, requestID)
			writer.Header().Set(RequestIDHeader, requestID)

			defer func() {
				logger.InfoContext(
					request.Context(),
					"solicitud procesada por Gateway",
					slog.String("request_id", requestID),
					slog.String("method", request.Method),
					slog.String("path", request.URL.Path),
					slog.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
