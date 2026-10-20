package proxy

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type problemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func writeProblem(
	writer http.ResponseWriter,
	request *http.Request,
	logger *slog.Logger,
	problem problemDetails,
) {
	writer.Header().Set(
		"Content-Type",
		"application/problem+json",
	)

	writer.WriteHeader(problem.Status)

	if err := json.NewEncoder(writer).Encode(problem); err != nil {
		logger.ErrorContext(
			request.Context(),
			"error escribiendo problem details",
			slog.Any("error", err),
			slog.String("path", request.URL.Path),
		)
	}
}
