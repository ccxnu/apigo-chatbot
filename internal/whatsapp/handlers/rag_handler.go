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
	llmProvider  llm.Provider
	client       WhatsAppClient
	paramCache   domain.ParameterCache
	priority     int
}

func NewRAGHandler(
	chunkUseCase domain.ChunkUseCase,
	convUseCase domain.ConversationUseCase,
	llmProvider llm.Provider,
	client WhatsAppClient,
	paramCache domain.ParameterCache,
	priority int,
) *RAGHandler {
	return &RAGHandler{
		chunkUseCase: chunkUseCase,
		convUseCase:  convUseCase,
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
		llmResponse, err = h.generateLLMResponse(ctx, query, "", conversationHistory)
		if err != nil {
			return h.sendMessage(msg.ChatID, h.getParam("RAG_NO_RESULTS_MESSAGE", "Lo siento, no encontré información relevante sobre tu consulta."))
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
