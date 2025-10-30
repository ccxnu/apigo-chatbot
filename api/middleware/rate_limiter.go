package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
	"golang.org/x/time/rate"
)

// RateLimiterStore holds rate limiters for each API key
type RateLimiterStore struct {
	limiters map[int]*rate.Limiter
	mu       sync.RWMutex
}

// NewRateLimiterStore creates a new rate limiter store
func NewRateLimiterStore() *RateLimiterStore {
	store := &RateLimiterStore{
		limiters: make(map[int]*rate.Limiter),
	}

	// Start cleanup goroutine to remove old limiters (every hour)
	go store.cleanup()

	return store
}

// GetLimiter returns a rate limiter for the given API key
func (s *RateLimiterStore) GetLimiter(apiKey *d.APIKey) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, exists := s.limiters[apiKey.ID]
	if !exists {
		// Calculate rate limit per second and burst
		// If rate_limit is 1000/hour, that's ~0.28 requests/second
		// We'll use a burst of rate_limit/60 (per minute) to allow short bursts
		ratePerSecond := float64(apiKey.RateLimit) / 3600.0
		burst := apiKey.RateLimit / 60
		if burst < 1 {
			burst = 1
		}

		limiter = rate.NewLimiter(rate.Limit(ratePerSecond), burst)
		s.limiters[apiKey.ID] = limiter
	}

	return limiter
}

// cleanup removes old limiters periodically to prevent memory leaks
func (s *RateLimiterStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		// In a production system, you'd want to track last access time
		// and remove limiters that haven't been used in X hours
		// For now, we'll keep all limiters
		s.mu.Unlock()
	}
}

// RateLimiter middleware enforces rate limits per API key
func RateLimiter(store *RateLimiterStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get API key from context (set by APIKeyAuth middleware)
			apiKey, ok := GetAPIKeyFromContext(ctx)
			if !ok {
				// Should never happen if APIKeyAuth middleware is applied first
				logger.LogError(ctx, "API key not found in context", nil,
					"middleware", "RateLimiter",
					"path", r.URL.Path,
				)
				writeJSONError(w, http.StatusInternalServerError, "ERR_INTERNAL", "API key not found in context")
				return
			}

			// Get rate limiter for this API key
			limiter := store.GetLimiter(apiKey)

			// Check if request is allowed
			if !limiter.Allow() {
				logger.LogWarn(ctx, "Rate limit exceeded",
					"middleware", "RateLimiter",
					"keyID", apiKey.ID,
					"keyName", apiKey.Name,
					"rateLimit", apiKey.RateLimit,
					"path", r.URL.Path,
				)

				// Add rate limit headers
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", apiKey.RateLimit))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()))

				writeJSONError(w, http.StatusTooManyRequests, "ERR_RATE_LIMIT_EXCEEDED",
					fmt.Sprintf("Rate limit of %d requests per hour exceeded", apiKey.RateLimit))
				return
			}

			// Add rate limit headers to successful responses
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", apiKey.RateLimit))
			// Note: Getting precise remaining count from rate.Limiter is not straightforward
			// In production, you might want to use a different rate limiting implementation
			// that tracks exact counts (e.g., Redis-based)

			next.ServeHTTP(w, r)
		})
	}
}
