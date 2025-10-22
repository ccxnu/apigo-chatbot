package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"api-chatbot/domain"
	"api-chatbot/internal/llm"
	"api-chatbot/internal/whatsapp"
	"api-chatbot/internal/whatsapp/handlers"
)

// InitializeWhatsAppService creates and initializes a WhatsApp service
func InitializeWhatsAppService(
	app Application,
	sessionUC domain.WhatsAppSessionUseCase,
	chunkUC domain.ChunkUseCase,
	userUC domain.WhatsAppUserUseCase,
	convUC domain.ConversationUseCase,
) (*whatsapp.Service, error) {
	// Check if WhatsApp is enabled in parameters
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

	// Create SQL store for device data (whatsmeow requirement)
	// This is separate from our main PostgreSQL and only stores device keys
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

	// Get or create device for this session
	ctx := context.Background()
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}

	// Create WhatsApp client first
	waClient, err := whatsapp.NewClient(whatsapp.Config{
		SessionName: sessionName,
		DeviceStore: deviceStore,
		LogLevel:    "INFO",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp client: %w", err)
	}

	// Initialize LLM provider from configuration
	llmProvider, err := createLLMProvider(app.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM provider: %w", err)
	}

	// Initialize all message handlers in priority order
	messageHandlers := []whatsapp.MessageHandler{
		// User validation handler (HIGHEST priority - runs first)
		handlers.NewUserValidationHandler(userUC, convUC, waClient, 1000),

		// Command handler (high priority)
		handlers.NewCommandHandler(waClient, 100),

		// RAG handler (lowest priority, catches all other text)
		handlers.NewRAGHandler(chunkUC, convUC, llmProvider, waClient, 50),
	}

	// Create WhatsApp service
	service, err := whatsapp.NewServiceWithClient(waClient, sessionName, sessionUC, messageHandlers)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp service: %w", err)
	}

	// Start the service
	ctx = context.Background()
	if err := service.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start WhatsApp service: %w", err)
	}

	slog.Info("WhatsApp service initialized successfully", "session", sessionName)
	return service, nil
}

// createDeviceStore creates a SQL store for whatsmeow device data
// This is required by whatsmeow and is separate from our main database
func createDeviceStore(connString string) (*sqlstore.Container, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open device database: %w", err)
	}

	container := sqlstore.NewWithDB(db, "postgres", nil)

	// Create tables if they don't exist
	ctx := context.Background()
	if err := container.Upgrade(ctx); err != nil {
		return nil, fmt.Errorf("failed to upgrade device store: %w", err)
	}

	return container, nil
}

// createLLMProvider creates an LLM provider based on configuration
func createLLMProvider(cache domain.ParameterCache) (llm.Provider, error) {
	// Get LLM configuration from cache
	param, exists := cache.Get("LLM_CONFIG")
	if !exists {
		return nil, fmt.Errorf("LLM_CONFIG parameter not found")
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM_CONFIG: %w", err)
	}

	// Extract configuration
	provider, _ := data["provider"].(string)
	apiKey, _ := data["apiKey"].(string)
	baseURL, _ := data["baseURL"].(string)
	model, _ := data["model"].(string)
	temperature, _ := data["temperature"].(float64)
	maxTokens, _ := data["maxTokens"].(float64)
	timeout, _ := data["timeout"].(float64)
	systemPrompt, _ := data["systemPrompt"].(string)

	// Validate required fields
	if apiKey == "" || baseURL == "" || model == "" {
		return nil, fmt.Errorf("LLM_CONFIG missing required fields (apiKey, baseURL, model)")
	}

	// Create LLM configuration
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

	// Create OpenAI-compatible provider (works with Groq, OpenAI, etc.)
	llmProvider := llm.NewOpenAICompatibleProvider(config)

	slog.Info("LLM provider initialized",
		"provider", provider,
		"model", model,
		"baseURL", baseURL,
	)

	return llmProvider, nil
}
