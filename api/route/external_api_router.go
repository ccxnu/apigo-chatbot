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
	conversationUseCase d.ConversationUseCase,
	mux *http.ServeMux,
	humaAPI huma.API,
) {
	// POST /v1/chat/completions
	huma.Register(humaAPI, huma.Operation{
		OperationID: "chat-completions",
		Method:      http.MethodPost,
		Path:        "/api/v1/chat/completions",
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
			"deviceID", input.Body.IdDevice,
		)

		// Generate unique ID
		completionID := generateCompletionID()

		// Get or create conversation using device ID as chat ID
		chatID := input.Body.IdDevice
		convResult := conversationUseCase.GetOrCreateConversation(
			ctx,
			chatID,
			input.Body.DeviceAddress, // Using device address as phone number
			nil,                      // No contact name
			false,                    // Not a group
			nil,                      // No group name
		)

		if !convResult.Success {
			logger.LogError(ctx, "Failed to get or create conversation", nil,
				"operation", "ChatCompletions",
				"deviceID", chatID,
			)
			// Continue anyway, but won't save history
		}

		var conversationID int
		if convResult.Success && convResult.Data != nil {
			conversationID = convResult.Data.ID
		}

		// Retrieve conversation history from database
		var conversationHistory []llm.Message
		if conversationID > 0 {
			logger.LogInfo(ctx, "Attempting to retrieve conversation history",
				"operation", "ChatCompletions",
				"deviceID", chatID,
				"conversationID", conversationID,
			)
			historyResult := conversationUseCase.GetConversationHistory(ctx, chatID, 50)
			logger.LogInfo(ctx, "History result",
				"operation", "ChatCompletions",
				"success", historyResult.Success,
				"dataLength", len(historyResult.Data),
			)
			if historyResult.Success && len(historyResult.Data) > 0 {
				// Convert database messages to LLM messages
				for _, msg := range historyResult.Data {
					if msg.Body != nil && *msg.Body != "" {
						role := "user"
						if msg.FromMe || msg.SenderType == "bot" {
							role = "assistant"
						}
						conversationHistory = append(conversationHistory, llm.Message{
							Role:    role,
							Content: *msg.Body,
						})
					}
				}
				logger.LogInfo(ctx, "Retrieved conversation history",
					"operation", "ChatCompletions",
					"deviceID", chatID,
					"historyCount", len(conversationHistory),
				)
			}
		}

		// Extract user message (last message)
		var userMessage string
		for i := len(input.Body.Messages) - 1; i >= 0; i-- {
			if input.Body.Messages[i].Role == "user" {
				userMessage = input.Body.Messages[i].Content
				break
			}
		}

		// Query expansion: If event_filter is provided, expand the query with event keywords
		// This improves semantic similarity for generic queries like "De qué es este evento"
		expandedQuery := userMessage
		if input.Body.RAGConfig != nil && len(input.Body.RAGConfig.EventFilter) > 0 {
			eventCategory := input.Body.RAGConfig.EventFilter[0]
			// Add event name to query for better semantic matching
			if eventCategory == "DOC_INDTEC" {
				expandedQuery = userMessage + " INDTEC congreso tecnología"
				logger.LogInfo(ctx, "Query expanded for better semantic matching",
					"operation", "ChatCompletions",
					"originalQuery", userMessage,
					"expandedQuery", expandedQuery,
				)
			}
		}

		// Save user message to database
		if conversationID > 0 && userMessage != "" {
			userMessageID := fmt.Sprintf("msg-%s-user", completionID)
			userParams := d.CreateConversationMessageParams{
				ConversationID: conversationID,
				MessageID:      userMessageID,
				FromMe:         false,
				SenderType:     "user",
				MessageType:    "text",
				Body:           &userMessage,
				Timestamp:      time.Now().Unix(),
				IsForwarded:    false,
			}
			saveResult := conversationUseCase.StoreMessageWithStats(ctx, userParams)
			logger.LogInfo(ctx, "Saved user message",
				"operation", "ChatCompletions",
				"success", saveResult.Success,
				"conversationID", conversationID,
			)
		}

		// Prepare RAG context if enabled
		var ragContext *d.RAGContextInfo
		var retrievedContext string
		var selectedCategory *string

		if input.Body.RAGConfig != nil && input.Body.RAGConfig.Enabled {
			// Set defaults
			searchLimit := input.Body.RAGConfig.SearchLimit
			if searchLimit == 0 {
				searchLimit = 7
			}
			minSimilarity := input.Body.RAGConfig.MinSimilarity
			if minSimilarity == 0 {
				minSimilarity = 0.7
			}
			keywordWeight := input.Body.RAGConfig.KeywordWeight
			if keywordWeight == 0 {
				keywordWeight = 0.3
			}

			// Determine category filter from EventFilter (use first one if provided)
			if len(input.Body.RAGConfig.EventFilter) > 0 {
				selectedCategory = &input.Body.RAGConfig.EventFilter[0]
				logger.LogInfo(ctx, "Using category filter for RAG search",
					"operation", "ChatCompletions",
					"category", *selectedCategory,
				)
			}

			// Context Injection: Always inject base context for specific event categories
			// This ensures the LLM has essential information even for generic queries
			var contextBuilder strings.Builder
			if selectedCategory != nil && *selectedCategory != "" {
				baseContextCode := "BASE_CONTEXT_" + *selectedCategory
				if param, exists := cache.Get(baseContextCode); exists {
					dataMap, _ := param.GetDataAsMap()
					if baseContext, ok := dataMap["context"].(string); ok && baseContext != "" {
						contextBuilder.WriteString("Essential Event Information:\n")
						contextBuilder.WriteString(baseContext)
						contextBuilder.WriteString("\n\n")
						logger.LogInfo(ctx, "Base context injected for event category",
							"operation", "ChatCompletions",
							"category", *selectedCategory,
							"baseContextCode", baseContextCode,
						)
					}
				}
			}

			// Perform hybrid search with category filter using expanded query
			logger.LogInfo(ctx, "Performing RAG search",
				"operation", "ChatCompletions",
				"query", expandedQuery,
				"searchLimit", searchLimit,
				"category", func() string {
					if selectedCategory != nil {
						return *selectedCategory
					}
					return "none"
				}(),
			)

			var searchResult d.Result[[]d.ChunkWithHybridSimilarity]
			if selectedCategory != nil {
				// Use category-filtered search with expanded query
				searchResult = chunkUseCase.HybridSearchWithCategory(ctx, expandedQuery, searchLimit, minSimilarity, keywordWeight, selectedCategory)
			} else {
				// Use regular search with expanded query
				searchResult = chunkUseCase.HybridSearch(ctx, expandedQuery, searchLimit, minSimilarity, keywordWeight)
			}

			if searchResult.Success && len(searchResult.Data) > 0 {
				chunks := searchResult.Data

				// Build RAG context with retrieved chunks
				ragContext = &d.RAGContextInfo{
					ChunksRetrieved: len(chunks),
					Sources:         make([]d.SourceInfo, 0, len(chunks)),
				}

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
			} else {
				// No chunks retrieved, but we may still have base context
				logger.LogWarn(ctx, "No chunks retrieved from RAG search",
					"operation", "ChatCompletions",
					"query", expandedQuery,
					"hasBaseContext", contextBuilder.Len() > 0,
				)
			}

			retrievedContext = contextBuilder.String()
		}

		// Build LLM request
		llmRequest := llm.GenerateRequest{
			UserMessage:         userMessage,
			Context:             retrievedContext,
			ConversationHistory: conversationHistory, // Use history from database
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

		// Set system prompt based on category
		// If category is specified, try to get category-specific prompt first (e.g., "RAG_SYSTEM_PROMPT_DOC_INDECT")
		// Otherwise, fallback to general "RAG_SYSTEM_PROMPT"
		if selectedCategory != nil && *selectedCategory != "" {
			categoryPromptCode := "RAG_SYSTEM_PROMPT_" + *selectedCategory
			if param, exists := cache.Get(categoryPromptCode); exists {
				dataMap, _ := param.GetDataAsMap()
				if systemPrompt, ok := dataMap["message"].(string); ok {
					llmRequest.SystemPrompt = systemPrompt
					logger.LogInfo(ctx, "Using category-specific system prompt",
						"operation", "ChatCompletions",
						"promptCode", categoryPromptCode,
					)
				}
			}
		}

		// Fallback to general system prompt if no category-specific prompt was found
		if llmRequest.SystemPrompt == "" {
			if param, exists := cache.Get("RAG_SYSTEM_PROMPT"); exists {
				dataMap, _ := param.GetDataAsMap()
				if systemPrompt, ok := dataMap["message"].(string); ok {
					llmRequest.SystemPrompt = systemPrompt
					logger.LogInfo(ctx, "Using general system prompt",
						"operation", "ChatCompletions",
					)
				}
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

		// Save assistant response to database
		if conversationID > 0 && llmResponse.Content != "" {
			assistantMessageID := fmt.Sprintf("msg-%s-assistant", completionID)
			assistantParams := d.CreateConversationMessageParams{
				ConversationID:   conversationID,
				MessageID:        assistantMessageID,
				FromMe:           true,
				SenderType:       "bot",
				MessageType:      "text",
				Body:             &llmResponse.Content,
				Timestamp:        time.Now().Unix(),
				IsForwarded:      false,
				PromptTokens:     llmResponse.PromptTokens,
				CompletionTokens: llmResponse.CompletionTokens,
				TotalTokens:      llmResponse.TotalTokens,
			}
			conversationUseCase.StoreMessageWithStats(ctx, assistantParams)
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
