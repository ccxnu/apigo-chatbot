package main

import (
	"log/slog"
	"net/http"
	"time"

	"api-chatbot/api/middleware"
	"api-chatbot/api/route"
	"api-chatbot/config"
)

func main() {

	app := config.App()

	defer app.Shutdown() // Gracefully close logger and DB

	// Get context timeout from parameter cache
	timeout := 2 * time.Second // default
	if param, exists := app.Cache.Get("APP_CONFIG"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if ctxTimeout, ok := data["contextTimeout"].(float64); ok {
				timeout = time.Duration(ctxTimeout) * time.Second
			}
		}
	}

	mux := http.NewServeMux()

	humaAPI := config.NewHumaAPI(mux, app.Cache)

	route.Setup(app.Cache, timeout, app.Db, mux, humaAPI)

	// Global middlewares (order matters: Logging -> CORS -> Auth -> Handler)
	handler := middleware.LoggingMiddleware(
		middleware.CORSMiddleware(
			middleware.AuthMiddleware(mux, app.Cache),
		),
	)

	serverAddress := ":8080"
	slog.Info("Server starting", "address", serverAddress, "port", 8080)
	slog.Info("OpenAPI documentation available", "url", "http://localhost:8080/docs")
	if err := http.ListenAndServe(serverAddress, handler); err != nil {
		slog.Error("Could not start server", "error", err)
	}
}
