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

// RAGHandler handles Q&A using the RAG system with conversation history
type RAGHandler struct {
	chunkUseCase domain.ChunkUseCase
	convUseCase  domain.ConversationUseCase
	llmProvider  llm.Provider
	client       WhatsAppClient
	priority     int
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(
	chunkUseCase domain.ChunkUseCase,
	convUseCase domain.ConversationUseCase,
	llmProvider llm.Provider,
	client WhatsAppClient,
	priority int,
) *RAGHandler {
	return &RAGHandler{
		chunkUseCase: chunkUseCase,
		convUseCase:  convUseCase,
		llmProvider:  llmProvider,
		client:       client,
		priority:     priority,
	}
}

// Match determines if this handler should process the message
func (h *RAGHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	// Skip commands (anything starting with /)
	if len(msg.Body) > 0 && msg.Body[0] == '/' {
		return false
	}

	// Skip own messages
	if msg.FromMe {
		return false
	}

	// ONLY respond to personal/direct messages - skip groups
	if msg.IsGroup {
		return false
	}

	// Skip WhatsApp status broadcasts
	if msg.ChatID == "status@broadcast" {
		return false
	}

	// Only process text messages with content
	msgType := msg.MessageType
	if msgType == "" {
		msgType = "text"
	}

	return msgType == "text" && len(msg.Body) > 0
}

// Handle processes the message using RAG with conversation history
func (h *RAGHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	// Get user's question
	query := strings.TrimSpace(msg.Body)

	// 1. Get or create conversation
	convResult := h.convUseCase.GetOrCreateConversation(
		ctx,
		msg.ChatID,
		msg.ChatID, // Using chatID as phone number for now
		nil,        // contactName
		msg.IsGroup,
		nil, // groupName
	)
	if !convResult.Success {
		logger.LogError(ctx, "Failed to get/create conversation", nil,
			"chatID", msg.ChatID,
			"error", convResult.Code,
		)
		return h.sendMessage(msg.ChatID, "Lo siento, ocurrió un error al procesar tu mensaje.")
	}

	conversation := convResult.Data
	timestamp := time.Now().Unix()

	// 2. Store incoming user message
	storeResult := h.convUseCase.StoreMessage(
		ctx,
		conversation.ID,
		msg.MessageID,
		false, // from user
		query,
		timestamp,
	)
	if !storeResult.Success {
		logger.LogWarn(ctx, "Failed to store user message", "error", storeResult.Code)
		// Continue anyway - storing message is not critical
	}

	// 3. Get conversation history for context
	historyResult := h.convUseCase.GetConversationHistory(ctx, msg.ChatID, 10)
	var conversationHistory []llm.Message
	if historyResult.Success && len(historyResult.Data) > 0 {
		// Convert to LLM message format (reverse order - oldest first)
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

	// 4. Perform hybrid similarity search to find relevant chunks
	searchResult := h.chunkUseCase.HybridSearch(ctx, query, 5, 0.2, 0.15)

	if !searchResult.Success {
		logger.LogError(ctx, "Hybrid search failed", nil, "error", searchResult.Code)
		return h.sendMessage(msg.ChatID, "Lo siento, ocurrió un error al buscar información relevante.")
	}

	// 5. Check if we found relevant information
	if len(searchResult.Data) == 0 {
		// No relevant chunks found - use LLM with conversation history only
		answer, err := h.generateLLMResponse(ctx, query, "", conversationHistory)
		if err != nil {
			return h.sendMessage(msg.ChatID, "Lo siento, no encontré información relevante sobre tu consulta. ¿Podrías reformular tu pregunta?")
		}
		// Store assistant message
		h.storeAssistantMessage(ctx, conversation.ID, answer, timestamp+1)
		return h.sendMessage(msg.ChatID, answer)
	}

	// 6. Build context from retrieved chunks
	contextStr := h.buildHybridContext(searchResult.Data)

	// 7. Generate response using LLM with RAG context and conversation history
	answer, err := h.generateLLMResponse(ctx, query, contextStr, conversationHistory)
	if err != nil {
		logger.LogError(ctx, "LLM generation failed", err)
		// Fallback to simple answer if LLM fails
		answer = h.generateSimpleAnswer(searchResult.Data)
	}

	// 8. Store assistant response
	h.storeAssistantMessage(ctx, conversation.ID, answer, timestamp+2)

	// 9. Send response
	return h.sendMessage(msg.ChatID, answer)
}

// Priority returns the handler priority (lower than command handlers)
func (h *RAGHandler) Priority() int {
	return h.priority
}

// buildHybridContext creates context string from hybrid search results
func (h *RAGHandler) buildHybridContext(chunks []domain.ChunkWithHybridSimilarity) string {
	var builder strings.Builder

	for i, chunk := range chunks {
		builder.WriteString(fmt.Sprintf("## Fuente %d: %s\n", i+1, chunk.DocTitle))
		builder.WriteString(chunk.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

// generateLLMResponse generates a response using the LLM provider
func (h *RAGHandler) generateLLMResponse(ctx context.Context, query, ragContext string, conversationHistory []llm.Message) (string, error) {
	// Check if LLM provider is available
	if h.llmProvider == nil || !h.llmProvider.IsAvailable() {
		return "", fmt.Errorf("LLM provider not available")
	}

	// Build request
	request := llm.GenerateRequest{
		SystemPrompt: "Eres un asistente virtual del instituto educativo. Tu objetivo es ayudar a estudiantes y profesores con información académica de manera clara, precisa y amigable. Siempre basa tus respuestas en el contexto proporcionado.",
		UserMessage:  query,
		Context:      ragContext,
		ConversationHistory: conversationHistory,
		Temperature:  0.7,
		MaxTokens:    1000,
	}

	// Generate response
	response, err := h.llmProvider.GenerateResponse(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to generate LLM response: %w", err)
	}

	return response.Content, nil
}

// storeAssistantMessage stores the assistant's response in conversation history
func (h *RAGHandler) storeAssistantMessage(ctx context.Context, conversationID int, message string, timestamp int64) {
	result := h.convUseCase.StoreMessage(
		ctx,
		conversationID,
		fmt.Sprintf("assistant_%d", timestamp),
		true, // from me (assistant)
		message,
		timestamp,
	)
	if !result.Success {
		logger.LogWarn(ctx, "Failed to store assistant message", "error", result.Code)
	}
}

// generateSimpleAnswer generates a simple fallback answer from chunks
func (h *RAGHandler) generateSimpleAnswer(chunks []domain.ChunkWithHybridSimilarity) string {
	if len(chunks) == 0 {
		return "No encontré información relevante."
	}

	// Return the most relevant chunk with source
	mostRelevant := chunks[0]

	answer := fmt.Sprintf(
		"📚 Basado en: *%s*\n\n%s\n\n_Relevancia: %.0f%%_",
		mostRelevant.DocTitle,
		mostRelevant.Content,
		mostRelevant.CombinedScore*100,
	)

	// If we have multiple sources, mention them
	if len(chunks) > 1 {
		answer += fmt.Sprintf("\n\n_(También encontré información en %d documentos más)_", len(chunks)-1)
	}

	return answer
}

// sendMessage sends a text message via WhatsApp
func (h *RAGHandler) sendMessage(chatID, message string) error {
	return h.client.SendText(chatID, message)
}
