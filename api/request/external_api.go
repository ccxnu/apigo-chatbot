package request

import "api-chatbot/domain"

// ChatCompletionsRequest represents the OpenAI-compatible chat completions request
type ChatCompletionsRequest struct {
	domain.Base
	Model       string             `json:"model" validate:"required" doc:"Model to use for completion"`
	Messages    []ChatMessageInput `json:"messages" validate:"required,min=1,dive" doc:"Array of messages in the conversation"`
	Temperature *float64           `json:"temperature,omitempty" validate:"omitempty,gte=0,lte=2" doc:"Sampling temperature (0-2)"`
	MaxTokens   *int               `json:"max_tokens,omitempty" validate:"omitempty,gt=0" doc:"Maximum tokens to generate"`
	Stream      *bool              `json:"stream,omitempty" doc:"Whether to stream the response (default: false)"`
	RAGConfig   *RAGConfig         `json:"rag_config,omitempty" doc:"RAG-specific configuration"`
}

// ChatMessageInput represents an input message in the conversation
type ChatMessageInput struct {
	Role    string `json:"role" validate:"required,oneof=system user assistant"`
	Content string `json:"content" validate:"required"`
}

// RAGConfig contains RAG-specific configuration
type RAGConfig struct {
	Enabled        bool     `json:"enabled"`
	SearchLimit    int      `json:"search_limit" validate:"omitempty,gt=0,lte=50"`
	MinSimilarity  float64  `json:"min_similarity" validate:"omitempty,gte=0,lte=1"`
	KeywordWeight  float64  `json:"keyword_weight" validate:"omitempty,gte=0,lte=1"`
	EventFilter    []string `json:"event_filter,omitempty"` // Filter by event categories (e.g., ["EVENT_INDTEC"])
}

// EmbeddingsRequest represents the OpenAI-compatible embeddings request
type EmbeddingsRequest struct {
	domain.Base
	Input          interface{} `json:"input" validate:"required" doc:"Text or array of texts to embed"`
	Model          string      `json:"model" validate:"required" doc:"Model to use for embeddings"`
	EncodingFormat string      `json:"encoding_format,omitempty" validate:"omitempty,oneof=float base64" doc:"Encoding format for embeddings"`
}

// SearchRequest represents the knowledge base search request
type SearchRequest struct {
	domain.Base
	Query          string   `json:"query" validate:"required" doc:"Search query text"`
	Limit          int      `json:"limit" validate:"omitempty,gt=0,lte=100" doc:"Maximum number of results to return"`
	MinSimilarity  float64  `json:"min_similarity" validate:"omitempty,gte=0,lte=1" doc:"Minimum similarity score (0-1)"`
	SearchType     string   `json:"search_type" validate:"omitempty,oneof=vector hybrid keyword" doc:"Type of search to perform"`
	KeywordWeight  float64  `json:"keyword_weight" validate:"omitempty,gte=0,lte=1" doc:"Weight for keyword search in hybrid mode"`
	EventFilter    []string `json:"event_filter,omitempty" doc:"Filter by event categories (e.g., ['INDTEC', 'CONGRESO'])"`
	IncludeContent bool     `json:"include_content" doc:"Whether to include chunk content in results"`
}
