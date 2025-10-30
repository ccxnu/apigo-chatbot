package route

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
	"api-chatbot/internal/llm"
	"api-chatbot/internal/logger"

	"github.com/danielgtaylor/huma/v2"
)

// Response types - wrapped in Result[T]
type ChatCompletionsResponse struct {
	Body d.Result[d.ChatCompletionsResponse]
}

func NewExternalAPIRouter(
	chunkUseCase d.ChunkUseCase,
	embeddingService d.EmbeddingService,
	llmProvider llm.Provider,
	cache d.ParameterCache,
	apiKeyUseCase d.APIKeyUseCase,
	apiUsageRepo d.APIUsageRepository,
	mux *http.ServeMux,
	humaAPI huma.API,
) {
	// POST /v1/chat/completions
	huma.Register(humaAPI, huma.Operation{
		OperationID: "chat-completions",
		Method:      http.MethodPost,
		Path:        "/v1/chat/completions",
		Summary:     "Create chat completion with RAG",
		Description: "OpenAI-compatible chat completions endpoint with RAG support and event filtering",
		Tags:        []string{"External API"},
	}, func(ctx context.Context, input *struct {
		Body request.ChatCompletionsRequest
	}) (*ChatCompletionsResponse, error) {
		startTime := time.Now()

		logger.LogInfo(ctx, "Processing chat completion request",
			"operation", "ChatCompletions",
			"model", input.Body.Model,
			"messageCount", len(input.Body.Messages),
		)

		// Generate unique ID
		completionID := generateCompletionID()

		// Extract user message (last message)
		var userMessage string
		for i := len(input.Body.Messages) - 1; i >= 0; i-- {
			if input.Body.Messages[i].Role == "user" {
				userMessage = input.Body.Messages[i].Content
				break
			}
		}

		// Prepare RAG context if enabled
		var ragContext *d.RAGContextInfo
		var retrievedContext string

		if input.Body.RAGConfig != nil && input.Body.RAGConfig.Enabled {
			// Set defaults
			searchLimit := input.Body.RAGConfig.SearchLimit
			if searchLimit == 0 {
				searchLimit = 5
			}
			minSimilarity := input.Body.RAGConfig.MinSimilarity
			if minSimilarity == 0 {
				minSimilarity = 0.7
			}
			keywordWeight := input.Body.RAGConfig.KeywordWeight
			if keywordWeight == 0 {
				keywordWeight = 0.3
			}

			// Perform hybrid search
			logger.LogInfo(ctx, "Performing RAG search",
				"operation", "ChatCompletions",
				"query", userMessage,
				"searchLimit", searchLimit,
			)

			searchResult := chunkUseCase.HybridSearch(ctx, userMessage, searchLimit, minSimilarity, keywordWeight)
			if searchResult.Success {
				chunks := searchResult.Data

				// Filter by event if specified
				if len(input.Body.RAGConfig.EventFilter) > 0 {
					filteredChunks := filterChunksByEvent(chunks, input.Body.RAGConfig.EventFilter)
					chunks = filteredChunks
				}

				// Build RAG context
				ragContext = &d.RAGContextInfo{
					ChunksRetrieved: len(chunks),
					Sources:         make([]d.SourceInfo, 0, len(chunks)),
				}

				var contextBuilder strings.Builder
				contextBuilder.WriteString("Relevant information from knowledge base:\n\n")

				for i, chunk := range chunks {
					contextBuilder.WriteString(fmt.Sprintf("Source %d (Document: %s):\n%s\n\n",
						i+1, chunk.DocTitle, chunk.Content))

					ragContext.Sources = append(ragContext.Sources, d.SourceInfo{
						DocumentID:    chunk.DocumentID,
						DocumentTitle: chunk.DocTitle,
						ChunkID:       chunk.ID,
						Similarity:    chunk.CombinedScore,
					})
				}

				retrievedContext = contextBuilder.String()
			}
		}

		// Build LLM request
		llmRequest := llm.GenerateRequest{
			UserMessage: userMessage,
			Context:     retrievedContext,
		}

		// Add conversation history
		for _, msg := range input.Body.Messages {
			if msg.Role != "user" || msg.Content != userMessage {
				llmRequest.ConversationHistory = append(llmRequest.ConversationHistory, llm.Message{
					Role:    msg.Role,
					Content: msg.Content,
				})
			}
		}

		// Set parameters
		if input.Body.Temperature != nil {
			llmRequest.Temperature = *input.Body.Temperature
		} else {
			llmRequest.Temperature = 0.7
		}

		if input.Body.MaxTokens != nil {
			llmRequest.MaxTokens = *input.Body.MaxTokens
		} else {
			llmRequest.MaxTokens = 1000
		}

		// Get system prompt from parameters
		if param, exists := cache.Get("RAG_SYSTEM_PROMPT"); exists {
			dataMap, _ := param.GetDataAsMap()
			if systemPrompt, ok := dataMap["message"].(string); ok {
				llmRequest.SystemPrompt = systemPrompt
			}
		}

		// Call LLM
		llmResponse, err := llmProvider.GenerateResponse(ctx, llmRequest)
		if err != nil {
			logger.LogError(ctx, "LLM generation failed", err,
				"operation", "ChatCompletions",
			)
			return nil, huma.Error500InternalServerError("Failed to generate response")
		}

		// Build completion data
		completionData := d.ChatCompletionsResponse{
			ID:      completionID,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   input.Body.Model,
			Choices: []d.ChatCompletionChoice{
				{
					Index: 0,
					Message: d.ChatMessage{
						Role:    "assistant",
						Content: llmResponse.Content,
					},
					FinishReason: llmResponse.FinishReason,
				},
			},
			RAGContext: ragContext,
		}

		// Add usage info if available
		if llmResponse.TotalTokens != nil {
			completionData.Usage = &d.UsageInfo{
				PromptTokens:     safeInt(llmResponse.PromptTokens),
				CompletionTokens: safeInt(llmResponse.CompletionTokens),
				TotalTokens:      *llmResponse.TotalTokens,
			}
		}

		logger.LogInfo(ctx, "Chat completion generated successfully",
			"operation", "ChatCompletions",
			"completionID", completionID,
			"duration", time.Since(startTime).Milliseconds(),
		)

		// Wrap in Result
		return &ChatCompletionsResponse{
			Body: d.Success(completionData),
		}, nil
	})

}

// Helper functions

func filterChunksByEvent(chunks []d.ChunkWithHybridSimilarity, eventFilter []string) []d.ChunkWithHybridSimilarity {
	if len(eventFilter) == 0 {
		return chunks
	}

	filtered := make([]d.ChunkWithHybridSimilarity, 0)
	for _, chunk := range chunks {
		for _, event := range eventFilter {
			// Check if document category matches the event filter
			// Case-insensitive partial match (e.g., "INDTEC" matches "EVENT_INDTEC")
			if strings.Contains(strings.ToUpper(chunk.DocCategory), strings.ToUpper(event)) {
				filtered = append(filtered, chunk)
				break
			}
		}
	}
	return filtered
}

func generateCompletionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "chatcmpl-" + hex.EncodeToString(bytes)
}

func safeInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
