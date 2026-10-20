package proxy

import (
	"log/slog"
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
		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}
}
