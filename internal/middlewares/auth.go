package middlewares

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ricardocuellar/devboard/internal/auth"
	"github.com/ricardocuellar/devboard/internal/domain"
)

type contextKey string

const (
	claimsKey contextKey = "auth_claims"
)

// RequireAuth es un middleware que verifica el JWT en el header Authorization
func RequireAuth(jwtManager *auth.JWTManager, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractBearerToken(r)

			if err != nil {
				respondAuthError(w, r, logger, http.StatusUnauthorized, "No autorizado", "Se requiere un token de acceso válido")
				return
			}

			// Verificar token

			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				respondAuthError(w, r, logger, http.StatusUnauthorized, "No autorizado", "El token de acceso es inválido o expiro")

				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", domain.ErrUnauthorized
	}

	// Bearer <token>
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", domain.ErrUnauthorized
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", domain.ErrUnauthorized
	}

	return token, nil
}

// GetClaims extrae los JWT claims del context del request
func GetClaims(r *http.Request) *auth.Claims {
	claims, _ := r.Context().Value(claimsKey).(*auth.Claims)
	return claims
}

// GetUserID es un helper que extrae el UserID del context
func GetUserID(r *http.Request) string {
	claims := GetClaims(r)
	if claims == nil {
		return ""
	}
	return claims.UserID
}

// RequireRole es un middleware que verifica que el usuario tiene el rol requerido
func RequireRole(role string, logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r)
			if claims == nil {
				logger.ErrorContext(r.Context(), "RequireRole usado sin RequireAuth",
					slog.String("path", r.URL.Path))

				respondAuthError(w, r, logger, http.StatusUnauthorized, "No autorizado", "No se encontró información de autenticación")

				return
			}

			if claims.Role != role {
				respondAuthError(w, r, logger, http.StatusForbidden, "Acceso prohibido", "No tienes permisos para acceder a este recurso")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// respondAuthError
func respondAuthError(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	status int,
	title string,
	detail string,
) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)

	response := struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
	}{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(
			r.Context(),
			"error al escribir respuesta de autenticación",
			"error", err,
			"path", r.URL.Path,
		)
	}
}
