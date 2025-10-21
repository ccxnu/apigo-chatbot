package llm

import (
	"context"
)

// Provider represents an LLM provider interface
// Allows easy switching between Grok, OpenAI, Anthropic, etc.
type Provider interface {
	// GenerateResponse generates a response based on the prompt and context
	GenerateResponse(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

	// GetProviderName returns the name of the LLM provider
	GetProviderName() string

	// IsAvailable checks if the provider is properly configured and available
	IsAvailable() bool
}

// GenerateRequest contains the input for generating a response
type GenerateRequest struct {
	// SystemPrompt sets the behavior/personality of the AI
	SystemPrompt string

	// UserMessage is the user's question or input
	UserMessage string

	// Context provides relevant information to help answer the question
	Context string

	// Temperature controls randomness (0.0 = deterministic, 1.0 = creative)
	Temperature float64

	// MaxTokens limits the response length
	MaxTokens int

	// ConversationHistory for multi-turn conversations (optional)
	ConversationHistory []Message
}

// GenerateResponse contains the LLM's response
type GenerateResponse struct {
	// Content is the generated text response
	Content string

	// TokensUsed indicates how many tokens were consumed
	TokensUsed int

	// Model is the specific model that generated the response
	Model string

	// FinishReason indicates why generation stopped (e.g., "stop", "length")
	FinishReason string
}

// Message represents a conversation message
type Message struct {
	Role    string // "system", "user", or "assistant"
	Content string
}

// Config holds configuration for LLM providers
type Config struct {
	// Provider name: "grok", "openai", "anthropic"
	Provider string

	// API key for authentication
	APIKey string

	// Model name (e.g., "grok-beta", "gpt-4", "claude-3-opus")
	Model string

	// Default temperature
	Temperature float64

	// Default max tokens
	MaxTokens int

	// Request timeout in seconds
	Timeout int

	// System prompt template
	SystemPrompt string
}

// Error types
type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeInvalidConfig   = "LLM_INVALID_CONFIG"
	ErrCodeAPIError        = "LLM_API_ERROR"
	ErrCodeTimeout         = "LLM_TIMEOUT"
	ErrCodeRateLimit       = "LLM_RATE_LIMIT"
	ErrCodeInvalidResponse = "LLM_INVALID_RESPONSE"
	ErrCodeUnavailable     = "LLM_UNAVAILABLE"
)
