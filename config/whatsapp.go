package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"api-chatbot/domain"
	"api-chatbot/internal/llm"
	"api-chatbot/internal/mailer"
	"api-chatbot/internal/whatsapp"
	"api-chatbot/internal/whatsapp/handlers"
	"api-chatbot/usecase"
)

func InitializeWhatsAppService(
	app Application,
	sessionUC domain.WhatsAppSessionUseCase,
	chunkUC domain.ChunkUseCase,
	userUC domain.WhatsAppUserUseCase,
	regUC domain.RegistrationUseCase,
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
		handlers.NewCommandHandler(waClient, app.Cache, regUC, userUC, convUC, 100),
		handlers.NewRegistrationHandler(regUC, userUC, convUC, waClient, app.Cache, 1000),
		handlers.NewRAGHandler(chunkUC, convUC, userUC, llmProvider, waClient, app.Cache, 50),
	}

	service, err := whatsapp.NewServiceWithClient(waClient, sessionName, sessionUC, messageHandlers, app.Cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp service: %w", err)
	}

	ctx = context.Background()
	if err := service.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start WhatsApp service: %w", err)
	}

	// Register service with global manager for QR code access
	whatsapp.GetManager().SetService(service)

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

// InitializeRegistrationUseCase creates and initializes the registration use case with OTP mailer
func InitializeRegistrationUseCase(
	regRepo domain.RegistrationRepository,
	userRepo domain.WhatsAppUserRepository,
	userUC domain.WhatsAppUserUseCase,
	httpClient domain.HTTPClient,
	cache domain.ParameterCache,
	timeout time.Duration,
) domain.RegistrationUseCase {
	// Get email configuration from parameters
	var tikeeURL, senderEmail string

	if param, exists := cache.Get("EMAIL_CONFIG"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if url, ok := data["tikeeURL"].(string); ok {
				tikeeURL = url
			}
			if sender, ok := data["senderEmail"].(string); ok {
				senderEmail = sender
			}
		}
	}

	// Default values if not configured
	if tikeeURL == "" {
		tikeeURL = "http://20.84.48.225:5056/api/emails/enviarDirecto"
		slog.Warn("EMAIL_CONFIG.tikeeURL not found, using default", "url", tikeeURL)
	}
	if senderEmail == "" {
		senderEmail = "automatizaciones@tikee.tech"
		slog.Warn("EMAIL_CONFIG.senderEmail not found, using default (AWS SES verified)", "email", senderEmail)
	}

	// Create OTP mailer
	otpMailer := mailer.NewOTPMailer(httpClient, tikeeURL, senderEmail, cache, timeout)

	// Create registration use case
	regUC := usecase.NewRegistrationUseCase(
		regRepo,
		userRepo,
		userUC,
		otpMailer,
		cache,
		timeout,
	)

	slog.Info("Registration use case initialized",
		"tikeeURL", tikeeURL,
		"senderEmail", senderEmail,
	)

	return regUC
}
