package main

import (
	"log/slog"
	"net/http"
	"time"

	"api-chatbot/api/dal"
	"api-chatbot/api/middleware"
	"api-chatbot/api/route"
	"api-chatbot/config"
	"api-chatbot/internal/embedding"
	"api-chatbot/internal/httpclient"
	"api-chatbot/internal/whatsapp"
	"api-chatbot/repository"
	"api-chatbot/usecase"
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

	// Initialize WhatsApp service
	whatsappService := initializeWhatsAppService(app, timeout)
	if whatsappService != nil {
		defer whatsappService.Stop()
		slog.Info("WhatsApp service running")
	}

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

// initializeWhatsAppService creates and starts the WhatsApp service
func initializeWhatsAppService(app config.Application, timeout time.Duration) *whatsapp.Service {
	// Initialize use cases needed for WhatsApp
	dataAccess := dal.NewDAL(app.Db)

	// Session use case
	sessionRepo := repository.NewWhatsAppSessionRepository(dataAccess)
	sessionUC := usecase.NewWhatsAppSessionUseCase(sessionRepo, app.Cache, timeout)

	// Chunk use case for RAG
	chunkRepo := repository.NewChunkRepository(dataAccess)
	statsRepo := repository.NewChunkStatisticsRepository(dataAccess)
	httpClient := httpclient.NewHTTPClient(app.Cache)
	embeddingService := embedding.NewOpenAIEmbeddingService(app.Cache, httpClient)
	chunkUC := usecase.NewChunkUseCase(chunkRepo, statsRepo, app.Cache, embeddingService, timeout)

	// Initialize WhatsApp service (returns nil if disabled in config)
	service, err := config.InitializeWhatsAppService(app, sessionUC, chunkUC)
	if err != nil {
		slog.Error("Failed to initialize WhatsApp service", "error", err)
		return nil
	}

	return service
}

