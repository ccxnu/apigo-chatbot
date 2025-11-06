package domain

import "time"

// ChatCompletionsResponse represents the OpenAI-compatible chat completions response
type ChatCompletionsResponse struct {
	ID         string                 `json:"id"`
	Object     string                 `json:"object"` // "chat.completion"
	Created    int64                  `json:"created"`
	Model      string                 `json:"model"`
	Choices    []ChatCompletionChoice `json:"choices"`
	Usage      *UsageInfo             `json:"usage,omitempty"`
	RAGContext *RAGContextInfo        `json:"rag_context,omitempty"`
}

// ChatCompletionChoice represents a choice in the chat completion
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"` // "stop", "length", "content_filter"
}

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// UsageInfo represents token usage information
type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// RAGContextInfo provides information about the RAG context used
type RAGContextInfo struct {
	ChunksRetrieved int          `json:"chunks_retrieved"`
	Sources         []SourceInfo `json:"sources"`
}

// SourceInfo represents information about a source document
type SourceInfo struct {
	DocumentID    int     `json:"document_id"`
	DocumentTitle string  `json:"document_title"`
	ChunkID       int     `json:"chunk_id"`
	Similarity    float64 `json:"similarity"`
}

// EmbeddingsResponse represents the OpenAI-compatible embeddings response
type EmbeddingsResponse struct {
	Object string          `json:"object"` // "list"
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  *EmbeddingUsage `json:"usage,omitempty"`
}

// EmbeddingData represents a single embedding
type EmbeddingData struct {
	Object    string    `json:"object"` // "embedding"
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage represents token usage for embeddings
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// SearchResponse represents the knowledge base search response
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int            `json:"total"`
	SearchType string         `json:"search_type"`
}

// SearchResult represents a single search result
type SearchResult struct {
	ChunkID         int             `json:"chunk_id"`
	DocumentID      int             `json:"document_id"`
	DocumentTitle   string          `json:"document_title"`
	Content         *string         `json:"content,omitempty"`
	SimilarityScore *float64        `json:"similarity_score,omitempty"`
	KeywordScore    *float64        `json:"keyword_score,omitempty"`
	CombinedScore   *float64        `json:"combined_score,omitempty"`
	Metadata        *SearchMetadata `json:"metadata,omitempty"`
}

// SearchMetadata contains additional metadata about the result
type SearchMetadata struct {
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StreamChunk represents a chunk in the streaming response
type StreamChunk struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"` // "chat.completion.chunk"
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []StreamChoiceChunk `json:"choices"`
}

// StreamChoiceChunk represents a choice in a stream chunk
type StreamChoiceChunk struct {
	Index        int              `json:"index"`
	Delta        ChatMessageDelta `json:"delta"`
	FinishReason *string          `json:"finish_reason,omitempty"`
}

// ChatMessageDelta represents a delta update in streaming
type ChatMessageDelta struct {
	Role    *string `json:"role,omitempty"`
	Content *string `json:"content,omitempty"`
}
