package proxy

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(
	target *url.URL,
	logger *slog.Logger,
) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(target)

			request.Out.Header.Del("X-User-ID")
			request.Out.Header.Del("X-User-Role")

			request.SetXForwarded()
		},
		ErrorHandler: func(
			writer http.ResponseWriter,
			request *http.Request,
			err error,
		) {
			if errors.Is(request.Context().Err(), context.Canceled) {
				return
			}

			logger.ErrorContext(
				request.Context(),
				"fallo comunicando con servicio de destino",
				slog.String("upstream", target.Host),
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Any("error", err),
			)

			writeProblem(
				writer,
				request,
				logger,
				problemDetails{
					Type:     "https://ricardocuellar.dev/errors/bad-gateway",
					Title:    "Error de comunicación con el servicio",
					Status:   http.StatusBadGateway,
					Detail:   "No fue posible obtener una respuesta del servicio de destino",
					Instance: request.URL.Path,
				},
			)
		},
		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}
}
