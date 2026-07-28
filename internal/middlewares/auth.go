package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ricardocuellar/devboard/internal/auth"
	"github.com/ricardocuellar/devboard/internal/domain"
	"github.com/ricardocuellar/devboard/internal/handler"
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
				handler.RespondError(w, r, logger, domain.ErrUnauthorized)
				return
			}

			// Verificar token

			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				handler.RespondError(w, r, logger, domain.ErrUnauthorized)
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

				handler.RespondError(w, r, logger, domain.ErrUnauthorized)
				return
			}

			if claims.Role != role {
				handler.RespondError(w, r, logger, domain.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
