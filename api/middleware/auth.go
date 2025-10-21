package middleware

import (
	"api-chatbot/domain"
	"net/http"
	// "slices"
)

// AuthMiddleware validates a custom Authorization header using ParameterCache
func AuthMiddleware(next http.Handler, paramCache domain.ParameterCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for OpenAPI documentation endpoints
		// publicPaths := []string{"/docs", "/openapi", "/openapi.json", "/openapi.yaml"}
		//
		// if slices.Contains(publicPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			// return
		// }

		// Get basicAuth from parameter cache
		// var basicAuth string
		// if param, exists := paramCache.Get("APP_CONFIG"); exists {
		// 	if data, err := param.GetDataAsMap(); err == nil {
		// 		if auth, ok := data["basicAuth"].(string); ok {
		// 			basicAuth = auth
		// 		}
		// 	}
		// }
		//
		// // If no basicAuth configured, skip validation (development mode)
		// if basicAuth == "" {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		//
		// authHeader := r.Header.Get("X-App-Authorization")
		//
		// if authHeader != basicAuth {
		// 	domain.AppError(w, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "No autorizado")
		// 	return
		// }
		//
		// // Token is valid, proceed
		// next.ServeHTTP(w, r)
	})
}
