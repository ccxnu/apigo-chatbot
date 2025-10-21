package handlers

import (
	"context"
	"fmt"
	"strings"

	"api-chatbot/domain"
	"api-chatbot/internal/whatsapp"
)

// RAGHandler handles Q&A using the RAG system
// Similar to your NestJS MessageHandler
type RAGHandler struct {
	whatsapp.BaseHandler
	chunkUseCase domain.ChunkUseCase
	filter       *whatsapp.MessageFilter
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(
	client *whatsapp.Client,
	convUseCase domain.ConversationUseCase,
	chunkUseCase domain.ChunkUseCase,
) *RAGHandler {
	return &RAGHandler{
		BaseHandler: whatsapp.BaseHandler{
			Client:      client,
			ConvUseCase: convUseCase,
		},
		chunkUseCase: chunkUseCase,
		filter:       whatsapp.NewMessageFilter(),
	}
}

// Match determines if this handler should process the message
func (h *RAGHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	// Skip commands
	if h.filter.IsCommand(msg, "commands") {
		return false
	}
	if h.filter.IsCommand(msg, "help") {
		return false
	}
	if h.filter.IsCommand(msg, "horarios") {
		return false
	}

	// Skip own messages
	if h.filter.IsFromMe(msg) {
		return false
	}

	// Only process text messages with content
	if !h.filter.IsTextMessage(msg) || !h.filter.HasBody(msg) {
		return false
	}

	return true
}

// Handle processes the message using RAG
func (h *RAGHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	// Get user's question
	query := strings.TrimSpace(msg.Body)

	// Perform similarity search to find relevant chunks
	result := h.chunkUseCase.SimilaritySearch(ctx, query, 5, 0.7)

	if !result.Success {
		// Error occurred - send error message
		return h.sendMessage(msg.ChatID, "Lo siento, ocurrió un error al procesar tu mensaje.")
	}

	// Check if we found relevant information
	if len(result.Data) == 0 {
		return h.sendMessage(msg.ChatID, "Lo siento, no encontré información relevante sobre tu consulta. ¿Podrías reformular tu pregunta?")
	}

	// Build context from chunks
	context := h.buildContext(result.Data)

	// TODO: Call LLM (Grok/OpenAI) to generate answer using the context
	// For now, send the most relevant chunk
	answer := h.generateSimpleAnswer(result.Data)

	return h.sendMessage(msg.ChatID, answer)
}

// Priority returns the handler priority (lower than command handlers)
func (h *RAGHandler) Priority() int {
	return 10 // Default priority
}

// buildContext creates context string from retrieved chunks
func (h *RAGHandler) buildContext(chunks []domain.ChunkWithSimilarity) string {
	var builder strings.Builder

	for i, chunk := range chunks {
		builder.WriteString(fmt.Sprintf("## Fuente %d: %s\n", i+1, chunk.DocTitle))
		builder.WriteString(chunk.Content)
		builder.WriteString("\n\n")
	}

	return builder.String()
}

// generateSimpleAnswer generates a simple answer from chunks
// TODO: Replace with LLM-generated response
func (h *RAGHandler) generateSimpleAnswer(chunks []domain.ChunkWithSimilarity) string {
	if len(chunks) == 0 {
		return "No encontré información relevante."
	}

	// For now, return the most relevant chunk with source
	mostRelevant := chunks[0]

	answer := fmt.Sprintf(
		"📚 Basado en: *%s*\n\n%s\n\n_Similitud: %.0f%%_",
		mostRelevant.DocTitle,
		mostRelevant.Content,
		mostRelevant.SimilarityScore*100,
	)

	// If we have multiple sources, mention them
	if len(chunks) > 1 {
		answer += fmt.Sprintf("\n\n_(También encontré información en %d documentos más)_", len(chunks)-1)
	}

	return answer
}

// sendMessage sends a text message to the chat
func (h *RAGHandler) sendMessage(chatID, message string) error {
	// Parse chat JID
	jid, err := h.parseJID(chatID)
	if err != nil {
		return fmt.Errorf("failed to parse chat ID: %w", err)
	}

	// Send message
	_, err = h.Client.SendTextMessage(jid, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// parseJID converts chat ID string to types.JID
func (h *RAGHandler) parseJID(chatID string) (interface{}, error) {
	// TODO: Implement proper JID parsing using whatsmeow types
	// For now, return placeholder
	return chatID, nil
}
