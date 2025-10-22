package middleware

import (
	"context"
	"net/http"
	"strings"

	"api-chatbot/domain"
	"api-chatbot/internal/jwttoken"
)

// ContextKey type for context keys
type ContextKey string

const (
	// AdminClaimsKey is the context key for admin claims
	AdminClaimsKey ContextKey = "admin_claims"
)

// JWTMiddleware validates JWT access tokens and injects claims into context
// This middleware protects admin routes
func JWTMiddleware(tokenService *jwttoken.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				domain.AppError(w, http.StatusUnauthorized, "ERR_NO_TOKEN", "Authorization header required")
				return
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				domain.AppError(w, http.StatusUnauthorized, "ERR_INVALID_TOKEN_FORMAT", "Authorization header must be 'Bearer <token>'")
				return
			}

			token := parts[1]

			// Validate token
			claims, err := tokenService.ValidateAccessToken(token)
			if err != nil {
				domain.AppError(w, http.StatusUnauthorized, "ERR_INVALID_TOKEN", "Invalid or expired token")
				return
			}

			// Inject claims into context
			ctx := context.WithValue(r.Context(), AdminClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission creates a middleware that checks if the user has a specific permission
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims, ok := r.Context().Value(AdminClaimsKey).(*jwttoken.CustomClaims)
			if !ok {
				domain.AppError(w, http.StatusUnauthorized, "ERR_NO_CLAIMS", "No authentication claims found")
				return
			}

			// Check if user has the required permission
			hasPermission := false
			for _, perm := range claims.Permissions {
				if perm == permission || perm == "*" {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				domain.AppError(w, http.StatusForbidden, "ERR_INSUFFICIENT_PERMISSIONS", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims, ok := r.Context().Value(AdminClaimsKey).(*jwttoken.CustomClaims)
			if !ok {
				domain.AppError(w, http.StatusUnauthorized, "ERR_NO_CLAIMS", "No authentication claims found")
				return
			}

			// Check if user has the required role
			if claims.Role != role && claims.Role != "super_admin" {
				domain.AppError(w, http.StatusForbidden, "ERR_INSUFFICIENT_ROLE", "Insufficient role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetAdminClaims extracts admin claims from request context
func GetAdminClaims(ctx context.Context) (*jwttoken.CustomClaims, bool) {
	claims, ok := ctx.Value(AdminClaimsKey).(*jwttoken.CustomClaims)
	return claims, ok
}
