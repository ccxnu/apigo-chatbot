package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"api-chatbot/domain"
	"api-chatbot/internal/llm"
	"api-chatbot/internal/logger"
)

type RAGHandler struct {
	chunkUseCase domain.ChunkUseCase
	convUseCase  domain.ConversationUseCase
	userUseCase  domain.WhatsAppUserUseCase
	llmProvider  llm.Provider
	client       WhatsAppClient
	paramCache   domain.ParameterCache
	priority     int
}

func NewRAGHandler(
	chunkUseCase domain.ChunkUseCase,
	convUseCase domain.ConversationUseCase,
	userUseCase domain.WhatsAppUserUseCase,
	llmProvider llm.Provider,
	client WhatsAppClient,
	paramCache domain.ParameterCache,
	priority int,
) *RAGHandler {
	return &RAGHandler{
		chunkUseCase: chunkUseCase,
		convUseCase:  convUseCase,
		userUseCase:  userUseCase,
		llmProvider:  llmProvider,
		client:       client,
		paramCache:   paramCache,
		priority:     priority,
	}
}

func (h *RAGHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	if len(msg.Body) > 0 && msg.Body[0] == '/' {
		return false
	}
	if msg.FromMe || msg.IsGroup || msg.ChatID == "status@broadcast" {
		return false
	}

	msgType := msg.MessageType
	if msgType == "" {
		msgType = "text"
	}

	return msgType == "text" && len(msg.Body) > 0
}

func (h *RAGHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	query := strings.TrimSpace(msg.Body)

	// Check if user is registered
	userResult := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
	isRegistered := userResult.Success && userResult.Data != nil

	// For unregistered users, check chat limit
	if !isRegistered {
		guestLimit := h.getParamInt("GUEST_CHAT_LIMIT", 20)

		// Get conversation to count messages
		convResult := h.convUseCase.GetOrCreateConversation(ctx, msg.ChatID, msg.SenderName, nil, msg.IsGroup, nil)
		if convResult.Success && convResult.Data != nil {
			// Count user messages from today (last 24 hours)
			historyResult := h.convUseCase.GetConversationHistory(ctx, msg.ChatID, 100)
			todayStart := time.Now().Add(-24 * time.Hour).Unix()
			userMessageCount := 0

			if historyResult.Success {
				for _, histMsg := range historyResult.Data {
					if !histMsg.FromMe && histMsg.Timestamp >= todayStart {
						userMessageCount++
					}
				}
			}

			logger.LogInfo(ctx, "Guest user message count check",
				"whatsapp", msg.From,
				"messageCount", userMessageCount,
				"limit", guestLimit,
			)

			// Check if limit is reached
			if userMessageCount >= guestLimit {
				limitMsg := h.getParam("MESSAGE_GUEST_LIMIT_REACHED",
					"📊 Has alcanzado el límite de mensajes para usuarios no registrados.\n\n✅ Para continuar chateando sin límites, regístrate usando:\n\n/register")
				return h.sendMessage(msg.ChatID, limitMsg)
			}

			// Warn when approaching limit (1 message left)
			if userMessageCount == guestLimit-1 {
				warningTemplate := h.getParam("MESSAGE_GUEST_LIMIT_WARNING",
					"⚠️ Te queda %d mensaje disponible hoy.\n\n💡 Regístrate con /register para chat ilimitado.")
				warningMsg := fmt.Sprintf(warningTemplate, guestLimit-userMessageCount)
				// Send warning but continue processing
				h.sendMessage(msg.ChatID, warningMsg)
			}
		}
	}

	convResult := h.convUseCase.GetOrCreateConversation(
		ctx,
		msg.ChatID,
		msg.SenderName,
		nil,
		msg.IsGroup,
		nil,
	)
	if !convResult.Success {
		logger.LogError(ctx, "Failed to get/create conversation", nil,
			"chatID", msg.ChatID,
			"error", convResult.Code,
		)
		return h.sendMessage(msg.ChatID, h.getParam("RAG_ERROR_MESSAGE", "Lo siento, ocurrió un error al procesar tu mensaje."))
	}

	conversation := convResult.Data
	timestamp := time.Now().Unix()

	storeResult := h.convUseCase.StoreMessage(ctx, conversation.ID, msg.MessageID, false, query, timestamp)
	if !storeResult.Success {
		logger.LogWarn(ctx, "Failed to store user message", "error", storeResult.Code)
	}

	historyLimit := h.getParamInt("RAG_CONVERSATION_HISTORY_LIMIT", 10)
	historyResult := h.convUseCase.GetConversationHistory(ctx, msg.ChatID, historyLimit)
	var conversationHistory []llm.Message
	if historyResult.Success && len(historyResult.Data) > 0 {
		for i := len(historyResult.Data) - 1; i >= 0; i-- {
			msgHistory := historyResult.Data[i]
			if msgHistory.Body != nil {
				role := "user"
				if msgHistory.FromMe {
					role = "assistant"
				}
				conversationHistory = append(conversationHistory, llm.Message{
					Role:    role,
					Content: *msgHistory.Body,
				})
			}
		}
	}

	searchLimit := h.getParamInt("RAG_SEARCH_LIMIT", 5)
	minSimilarity := h.getParamFloat("RAG_MIN_SIMILARITY", 0.2)
	keywordWeight := h.getParamFloat("RAG_KEYWORD_WEIGHT", 0.15)

	searchResult := h.chunkUseCase.HybridSearch(ctx, query, searchLimit, minSimilarity, keywordWeight)

	if !searchResult.Success {
		logger.LogError(ctx, "Hybrid search failed", nil, "error", searchResult.Code)
		return h.sendMessage(msg.ChatID, h.getParam("RAG_ERROR_MESSAGE", "Lo siento, ocurrió un error al buscar información relevante."))
	}

	var answer string
	var err error

	var llmResponse *llm.GenerateResponse

	if len(searchResult.Data) == 0 {
		// No results found - include contact information in context
		contactInfo := h.getContactInformation()
		llmResponse, err = h.generateLLMResponse(ctx, query, contactInfo, conversationHistory)
		if err != nil {
			noResultsMsg := h.getParam("RAG_NO_RESULTS_MESSAGE", "Lo siento, no encontré información relevante sobre tu consulta.")
			contactMsg := h.formatContactsForMessage()
			if contactMsg != "" {
				noResultsMsg += "\n\n" + contactMsg
			}
			return h.sendMessage(msg.ChatID, noResultsMsg)
		}
		answer = llmResponse.Content
	} else {
		contextStr := h.buildHybridContext(searchResult.Data)
		llmResponse, err = h.generateLLMResponse(ctx, query, contextStr, conversationHistory)
		if err != nil {
			logger.LogError(ctx, "LLM generation failed", err)
			answer = h.generateSimpleAnswer(searchResult.Data)
		} else {
			answer = llmResponse.Content
		}
	}

	h.storeAssistantMessage(ctx, conversation.ID, answer, timestamp+2, llmResponse)

	return h.sendMessage(msg.ChatID, answer)
}

func (h *RAGHandler) Priority() int {
	return h.priority
}

func (h *RAGHandler) buildHybridContext(chunks []domain.ChunkWithHybridSimilarity) string {
	var builder strings.Builder

	for i, chunk := range chunks {
		builder.WriteString(fmt.Sprintf("## Fuente %d: %s\n", i+1, chunk.DocTitle))
		builder.WriteString(chunk.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func (h *RAGHandler) generateLLMResponse(ctx context.Context, query, ragContext string, conversationHistory []llm.Message) (*llm.GenerateResponse, error) {
	if h.llmProvider == nil || !h.llmProvider.IsAvailable() {
		return nil, fmt.Errorf("LLM provider not available")
	}

	systemPrompt := h.getParam("RAG_SYSTEM_PROMPT", "Eres un asistente virtual del instituto educativo.")
	temperature := h.getParamFloat("RAG_LLM_TEMPERATURE", 0.7)
	maxTokens := h.getParamInt("RAG_LLM_MAX_TOKENS", 1000)

	request := llm.GenerateRequest{
		SystemPrompt:        systemPrompt,
		UserMessage:         query,
		Context:             ragContext,
		ConversationHistory: conversationHistory,
		Temperature:         temperature,
		MaxTokens:           maxTokens,
	}

	response, err := h.llmProvider.GenerateResponse(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to generate LLM response: %w", err)
	}

	return response, nil
}

func (h *RAGHandler) storeAssistantMessage(ctx context.Context, conversationID int, message string, timestamp int64, llmResponse *llm.GenerateResponse) {
	if llmResponse == nil {
		result := h.convUseCase.StoreMessage(ctx, conversationID, fmt.Sprintf("assistant_%d", timestamp), true, message, timestamp)
		if !result.Success {
			logger.LogWarn(ctx, "Failed to store assistant message", "error", result.Code)
		}
		return
	}

	params := domain.CreateConversationMessageParams{
		ConversationID:   conversationID,
		MessageID:        fmt.Sprintf("assistant_%d", timestamp),
		FromMe:           true,
		SenderType:       "bot",
		MessageType:      "text",
		Body:             &message,
		Timestamp:        timestamp,
		IsForwarded:      false,
		QueueTimeMs:      llmResponse.QueueTimeMs,
		PromptTokens:     llmResponse.PromptTokens,
		PromptTimeMs:     llmResponse.PromptTimeMs,
		CompletionTokens: llmResponse.CompletionTokens,
		CompletionTimeMs: llmResponse.CompletionTimeMs,
		TotalTokens:      llmResponse.TotalTokens,
		TotalTimeMs:      llmResponse.TotalTimeMs,
	}

	result := h.convUseCase.StoreMessageWithStats(ctx, params)
	if !result.Success {
		logger.LogWarn(ctx, "Failed to store assistant message with stats", "error", result.Code)
	}
}

func (h *RAGHandler) generateSimpleAnswer(chunks []domain.ChunkWithHybridSimilarity) string {
	if len(chunks) == 0 {
		return h.getParam("RAG_NO_RELEVANT_INFO", "No encontré información relevante.")
	}

	mostRelevant := chunks[0]

	answer := fmt.Sprintf(
		"📚 Basado en: *%s*\n\n%s\n\n_Relevancia: %.0f%%_",
		mostRelevant.DocTitle,
		mostRelevant.Content,
		mostRelevant.CombinedScore*100,
	)

	if len(chunks) > 1 {
		answer += fmt.Sprintf("\n\n_(También encontré información en %d documentos más)_", len(chunks)-1)
	}

	return answer
}

func (h *RAGHandler) sendMessage(chatID, message string) error {
	return h.client.SendText(chatID, message)
}

func (h *RAGHandler) getParam(code, defaultValue string) string {
	param, exists := h.paramCache.Get(code)
	if !exists {
		return defaultValue
	}
	data, err := param.GetDataAsMap()
	if err != nil {
		return defaultValue
	}
	if msg, ok := data["message"].(string); ok {
		return msg
	}
	if val, ok := data["value"].(string); ok {
		return val
	}
	return defaultValue
}

func (h *RAGHandler) getParamInt(code string, defaultValue int) int {
	param, exists := h.paramCache.Get(code)
	if !exists {
		return defaultValue
	}
	data, err := param.GetDataAsMap()
	if err != nil {
		return defaultValue
	}
	if val, ok := data["value"].(float64); ok {
		return int(val)
	}
	return defaultValue
}

func (h *RAGHandler) getParamFloat(code string, defaultValue float64) float64 {
	param, exists := h.paramCache.Get(code)
	if !exists {
		return defaultValue
	}
	data, err := param.GetDataAsMap()
	if err != nil {
		return defaultValue
	}
	if val, ok := data["value"].(float64); ok {
		return val
	}
	return defaultValue
}

// getContactInformation retrieves contact information to be used as context for LLM
func (h *RAGHandler) getContactInformation() string {
	param, exists := h.paramCache.Get("RAG_INFORMATION_CONTACT")
	if !exists {
		return ""
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return ""
	}

	message, _ := data["message"].(string)
	contacts, ok := data["contacts"].([]interface{})
	if !ok || len(contacts) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("## Información de contacto institucional\n\n")
	if message != "" {
		builder.WriteString(message)
		builder.WriteString("\n\n")
	}

	for _, contactData := range contacts {
		contact, ok := contactData.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := contact["name"].(string)
		position, _ := contact["position"].(string)
		email, _ := contact["email"].(string)

		if name != "" {
			builder.WriteString(fmt.Sprintf("- **%s**", name))
			if position != "" {
				builder.WriteString(fmt.Sprintf(" (%s)", position))
			}
			if email != "" {
				builder.WriteString(fmt.Sprintf("\n  Email: %s", email))
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// formatContactsForMessage formats contacts for direct WhatsApp message
func (h *RAGHandler) formatContactsForMessage() string {
	param, exists := h.paramCache.Get("RAG_INFORMATION_CONTACT")
	if !exists {
		return ""
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return ""
	}

	message, _ := data["message"].(string)
	contacts, ok := data["contacts"].([]interface{})
	if !ok || len(contacts) == 0 {
		return ""
	}

	var builder strings.Builder
	if message != "" {
		builder.WriteString(message)
		builder.WriteString("\n\n")
	}

	for _, contactData := range contacts {
		contact, ok := contactData.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := contact["name"].(string)
		position, _ := contact["position"].(string)
		email, _ := contact["email"].(string)

		if name != "" {
			builder.WriteString(fmt.Sprintf("👤 *%s*", name))
			if position != "" {
				builder.WriteString(fmt.Sprintf("\n   _%s_", position))
			}
			if email != "" {
				builder.WriteString(fmt.Sprintf("\n   📧 %s", email))
			}
			builder.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(builder.String())
}
