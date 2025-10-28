package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"api-chatbot/domain"
	"api-chatbot/internal/llm"
	"api-chatbot/internal/whatsapp"
	"api-chatbot/internal/whatsapp/handlers"
)

func InitializeWhatsAppService(
	app Application,
	sessionUC domain.WhatsAppSessionUseCase,
	chunkUC domain.ChunkUseCase,
	userUC domain.WhatsAppUserUseCase,
	convUC domain.ConversationUseCase,
) (*whatsapp.Service, error) {
	param, exists := app.Cache.Get("WHATSAPP_CONFIG")
	if !exists {
		slog.Info("WhatsApp service disabled - WHATSAPP_CONFIG parameter not found")
		return nil, nil
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return nil, fmt.Errorf("failed to parse WHATSAPP_CONFIG: %w", err)
	}

	enabled, ok := data["enabled"].(bool)
	if !ok || !enabled {
		slog.Info("WhatsApp service disabled in configuration")
		return nil, nil
	}

	sessionName, _ := data["sessionName"].(string)
	if sessionName == "" {
		sessionName = "chatbot-session"
	}

	dbConnString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		app.Env.Database.Host,
		app.Env.Database.Port,
		app.Env.Database.User,
		app.Env.Database.Password,
		app.Env.Database.Name,
	)

	container, err := createDeviceStore(dbConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to create device store: %w", err)
	}

	ctx := context.Background()
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}

	waClient, err := whatsapp.NewClient(whatsapp.Config{
		SessionName: sessionName,
		DeviceStore: deviceStore,
		LogLevel:    "INFO",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp client: %w", err)
	}

	llmProvider, err := createLLMProvider(app.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM provider: %w", err)
	}

	messageHandlers := []whatsapp.MessageHandler{
		handlers.NewUserValidationHandler(userUC, convUC, waClient, app.Cache, 1000),
		handlers.NewCommandHandler(waClient, app.Cache, 100),
		handlers.NewRAGHandler(chunkUC, convUC, llmProvider, waClient, app.Cache, 50),
	}

	service, err := whatsapp.NewServiceWithClient(waClient, sessionName, sessionUC, messageHandlers, app.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp service: %w", err)
	}

	ctx = context.Background()
	if err := service.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start WhatsApp service: %w", err)
	}

	slog.Info("WhatsApp service initialized successfully", "session", sessionName)
	return service, nil
}

func createDeviceStore(connString string) (*sqlstore.Container, error) {
	// Use pgx driver instead of lib/pq
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open device database: %w", err)
	}

	container := sqlstore.NewWithDB(db, "pgx", nil)

	ctx := context.Background()
	if err := container.Upgrade(ctx); err != nil {
		return nil, fmt.Errorf("failed to upgrade device store: %w", err)
	}

	return container, nil
}

func createLLMProvider(cache domain.ParameterCache) (llm.Provider, error) {
	param, exists := cache.Get("LLM_CONFIG")
	if !exists {
		return nil, fmt.Errorf("LLM_CONFIG parameter not found")
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM_CONFIG: %w", err)
	}

	provider, _ := data["provider"].(string)
	apiKey, _ := data["apiKey"].(string)
	baseURL, _ := data["baseURL"].(string)
	model, _ := data["model"].(string)
	temperature, _ := data["temperature"].(float64)
	maxTokens, _ := data["maxTokens"].(float64)
	timeout, _ := data["timeout"].(float64)
	systemPrompt, _ := data["systemPrompt"].(string)

	if apiKey == "" || baseURL == "" || model == "" {
		return nil, fmt.Errorf("LLM_CONFIG missing required fields (apiKey, baseURL, model)")
	}

	config := llm.Config{
		Provider:     provider,
		APIKey:       apiKey,
		BaseURL:      baseURL,
		Model:        model,
		Temperature:  temperature,
		MaxTokens:    int(maxTokens),
		Timeout:      int(timeout),
		SystemPrompt: systemPrompt,
	}

	llmProvider := llm.NewOpenAICompatibleProvider(config)

	slog.Info("LLM provider initialized",
		"provider", provider,
		"model", model,
		"baseURL", baseURL,
	)

	return llmProvider, nil
}
