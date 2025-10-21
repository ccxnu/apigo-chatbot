package handlers

import (
	"context"
	"fmt"
	"strings"

	"api-chatbot/domain"
)

// RAGHandler handles Q&A using the RAG system
// Similar to your NestJS MessageHandler
type RAGHandler struct {
	chunkUseCase domain.ChunkUseCase
	priority     int
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(
	chunkUseCase domain.ChunkUseCase,
	priority int,
) *RAGHandler {
	return &RAGHandler{
		chunkUseCase: chunkUseCase,
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

	// Only process text messages with content
	msgType := msg.MessageType
	if msgType == "" {
		msgType = "text"
	}

	return msgType == "text" && len(msg.Body) > 0
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
	contextStr := h.buildContext(result.Data)
	_ = contextStr // Will be used when LLM is integrated

	// TODO: Call LLM (Grok/OpenAI) to generate answer using the context
	// For now, send the most relevant chunk
	answer := h.generateSimpleAnswer(result.Data)

	return h.sendMessage(msg.ChatID, answer)
}

// Priority returns the handler priority (lower than command handlers)
func (h *RAGHandler) Priority() int {
	return h.priority
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

// sendMessage is a placeholder - actual sending will be done by the service
func (h *RAGHandler) sendMessage(chatID, message string) error {
	// TODO: Actual message sending will be handled by the WhatsApp service
	// For now, just log that we would send this message
	return nil
}
