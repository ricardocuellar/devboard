package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr              string
	IdentityURL           *url.URL
	WorkspaceURL          *url.URL
	TaskURL               *url.URL
	ResponseHeaderTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: envOrDefault("GATEWAY_HTTP_ADDR", ":8080"),
	}

	_, port, err := net.SplitHostPort(cfg.HTTPAddr)
	if err != nil {
		return Config{}, fmt.Errorf(
			"GATEWAY_HTTP_ADDR debe tener formato host:puerto: %w",
			err,
		)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return Config{}, fmt.Errorf(
			"GATEWAY_HTTP_ADDR debe tener un puerto entre 1 y 65535",
		)
	}

	cfg.IdentityURL, err = loadServiceURL(
		"IDENTITY_HTTP_URL",
		"http://localhost:8081",
	)
	if err != nil {
		return Config{}, err
	}

	cfg.WorkspaceURL, err = loadServiceURL(
		"WORKSPACE_HTTP_URL",
		"http://localhost:8082",
	)
	if err != nil {
		return Config{}, err
	}

	cfg.TaskURL, err = loadServiceURL(
		"TASK_HTTP_URL",
		"http://localhost:8083",
	)
	if err != nil {
		return Config{}, err
	}

	timeoutValue := envOrDefault(
		"GATEWAY_RESPONSE_HEADER_TIMEOUT",
		"5s",
	)

	cfg.ResponseHeaderTimeout, err = time.ParseDuration(timeoutValue)
	if err != nil {
		return Config{}, fmt.Errorf(
			"GATEWAY_RESPONSE_HEADER_TIMEOUT inválido: %w",
			err,
		)
	}

	if cfg.ResponseHeaderTimeout <= 0 {
		return Config{}, fmt.Errorf(
			"GATEWAY_RESPONSE_HEADER_TIMEOUT debe ser mayor que cero",
		)
	}

	return cfg, nil
}

func loadServiceURL(key, defaultValue string) (*url.URL, error) {
	value := envOrDefault(key, defaultValue)

	target, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("%s inválida: %w", key, err)
	}

	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf(
			"%s debe utilizar http o https",
			key,
		)
	}

	if target.Hostname() == "" {
		return nil, fmt.Errorf(
			"%s debe incluir un host",
			key,
		)
	}

	if target.User != nil ||
		(target.Path != "" && target.Path != "/") ||
		target.RawQuery != "" ||
		target.ForceQuery ||
		target.Fragment != "" {
		return nil, fmt.Errorf(
			"%s debe contener únicamente esquema y host, con puerto opcional",
			key,
		)
	}

	return target, nil
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
