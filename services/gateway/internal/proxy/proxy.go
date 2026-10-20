package proxy

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func New(
	target *url.URL,
	responseHeaderTimeout time.Duration,
	logger *slog.Logger,
) *httputil.ReverseProxy {

	transport := http.DefaultTransport.(*http.Transport).Clone()

	transport.DialContext = (&net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	transport.TLSHandshakeTimeout = 3 * time.Second
	transport.ResponseHeaderTimeout = responseHeaderTimeout

	return &httputil.ReverseProxy{
		Transport: transport,
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

			problem := problemDetails{
				Type:     "https://ricardocuellar.dev/errors/bad-gateway",
				Title:    "Error de comunicación con el servicio",
				Status:   http.StatusBadGateway,
				Detail:   "No fue posible obtener una respuesta del servicio de destino",
				Instance: request.URL.Path,
			}

			var networkError net.Error

			if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &networkError)) && networkError.Timeout() {
				problem.Type = "https://ricardocuellar.dev/errors/gateway-timeout"
				problem.Title = "Tiempo de espera agotado"
				problem.Status = http.StatusGatewayTimeout
				problem.Detail = "El servicio destino no respondió dentro del tiempo permitido"
			}

			logger.ErrorContext(
				request.Context(),
				"fallo comunicando con servicio de destino",
				slog.String("upstream", target.Host),
				slog.String("request_id", request.Header.Get("X-Request-ID")),
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", problem.Status),
				slog.Any("error", err),
			)

			writeProblem(
				writer,
				request,
				logger,
				problem,
			)
		},
		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}
}
