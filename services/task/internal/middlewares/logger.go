// Package middlewares logger
package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ricardocuellar/devboard/services/task/internal/requestid"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseRecorder(writer http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (rr *responseRecorder) WriteHeader(code int) {
	if !rr.written {
		rr.statusCode = code
		rr.written = true
		rr.ResponseWriter.WriteHeader(code)
	}
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if !rr.written {
		rr.WriteHeader(http.StatusOK)
	}
	return rr.ResponseWriter.Write(b)
}

// Logger middleware
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			ctx := requestid.WithContext(request.Context(), request.Header.Get("X-Request-ID"))

			request = request.WithContext(ctx)

			id := requestid.FromContext(ctx)
			writer.Header().Set("X-Request-ID", id)

			start := time.Now()
			rec := newResponseRecorder(writer)

			// llamar al siguiente handler
			next.ServeHTTP(rec, request)

			//log estructurado
			logger.InfoContext(request.Context(), "request completado",
				slog.String("request_id", id),
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_addr", request.RemoteAddr),
			)

		})
	}
}
