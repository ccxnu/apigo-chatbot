package middleware

import (
	"context"
	"net/http"
	"strings"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

// APIKeyContextKey is the context key for storing the validated API key
type APIKeyContextKey string

const (
	APIKeyKey APIKeyContextKey = "api_key"
)

// APIKeyAuth middleware validates API keys from Authorization header
func APIKeyAuth(apiKeyUseCase d.APIKeyUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Extract API key from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.LogWarn(ctx, "Missing Authorization header",
					"middleware", "APIKeyAuth",
					"path", r.URL.Path,
				)
				writeJSONError(w, http.StatusUnauthorized, "ERR_MISSING_API_KEY", "Authorization header is required")
				return
			}

			// Parse "Bearer sk_xxx"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.LogWarn(ctx, "Invalid Authorization header format",
					"middleware", "APIKeyAuth",
					"path", r.URL.Path,
				)
				writeJSONError(w, http.StatusUnauthorized, "ERR_INVALID_AUTH_HEADER", "Authorization header must be in format: Bearer <api_key>")
				return
			}

			apiKey := parts[1]
			ipAddress := getClientIP(r)
			endpoint := r.URL.Path

			logger.LogInfo(ctx, "Validating API key",
				"middleware", "APIKeyAuth",
				"endpoint", endpoint,
				"ipAddress", ipAddress,
			)

			// Validate API key
			result := apiKeyUseCase.ValidateAPIKey(ctx, apiKey, ipAddress, endpoint)
			if !result.Success {
				logger.LogWarn(ctx, "API key validation failed",
					"middleware", "APIKeyAuth",
					"code", result.Code,
					"message", result.Info,
					"endpoint", endpoint,
					"ipAddress", ipAddress,
				)
				writeJSONError(w, http.StatusUnauthorized, result.Code, result.Info)
				return
			}

			logger.LogInfo(ctx, "API key validated successfully",
				"middleware", "APIKeyAuth",
				"keyID", result.Data.ID,
				"keyName", result.Data.Name,
				"endpoint", endpoint,
			)

			// Store validated API key in context for downstream use
			ctx = context.WithValue(ctx, APIKeyKey, result.Data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAPIKeyFromContext retrieves the validated API key from the request context
func GetAPIKeyFromContext(ctx context.Context) (*d.APIKey, bool) {
	apiKey, ok := ctx.Value(APIKeyKey).(*d.APIKey)
	return apiKey, ok
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// RemoteAddr includes port, strip it
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// writeJSONError writes a JSON error response
func writeJSONError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// Write simple JSON error response
	response := `{"error":{"code":"` + code + `","message":"` + message + `"}}`
	w.Write([]byte(response))
}
