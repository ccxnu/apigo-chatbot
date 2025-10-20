package middleware

import (
	"api-chatbot/domain"
	"api-chatbot/internal/jwttoken"
	"context"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			t := strings.Split(authHeader, " ")
			if len(t) != 2 {
				domain.AppError(w, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "")
				return
			}

			authToken := t[1]
			authorized, err := jwttoken.IsAuthorized(authToken, secret)
			if !authorized {
				domain.AppError(w, http.StatusUnauthorized, "ERR_UNAUTHORIZED", err.Error())
				return
			}

			userID, err := jwttoken.ExtractIDFromToken(authToken, secret)
			if err != nil {
				domain.AppError(w, http.StatusUnauthorized, "ERR_UNAUTHORIZED", err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), "x-user-id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
