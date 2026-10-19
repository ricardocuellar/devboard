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
		},
		ErrorLog: slog.NewLogLogger(
			logger.Handler(),
			slog.LevelError,
		),
	}
}
