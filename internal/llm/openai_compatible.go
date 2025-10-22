package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-chatbot/internal/logger"
)

// OpenAICompatibleProvider implements Provider interface for OpenAI-compatible APIs
// This works with OpenAI, Groq, and other providers that follow the OpenAI API format
type OpenAICompatibleProvider struct {
	config  Config
	client  *http.Client
	baseURL string
}

// NewOpenAICompatibleProvider creates a new OpenAI-compatible provider
func NewOpenAICompatibleProvider(config Config) *OpenAICompatibleProvider {
	timeout := 30 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	return &OpenAICompatibleProvider{
		config: config,
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: config.BaseURL,
	}
}

// GenerateResponse generates a response using the OpenAI-compatible API
func (p *OpenAICompatibleProvider) GenerateResponse(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	// Build messages array
	messages := []map[string]string{}

	// Add system prompt
	if req.SystemPrompt != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": req.SystemPrompt,
		})
	}

	// Add conversation history if provided
	for _, msg := range req.ConversationHistory {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	// Add context if provided
	if req.Context != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": fmt.Sprintf("Contexto relevante:\n\n%s", req.Context),
		})
	}

	// Add user message
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.UserMessage,
	})

	// Build request body
	requestBody := map[string]interface{}{
		"model":    p.config.Model,
		"messages": messages,
	}

	// Add optional parameters
	if req.Temperature > 0 {
		requestBody["temperature"] = req.Temperature
	} else if p.config.Temperature > 0 {
		requestBody["temperature"] = p.config.Temperature
	}

	if req.MaxTokens > 0 {
		requestBody["max_tokens"] = req.MaxTokens
	} else if p.config.MaxTokens > 0 {
		requestBody["max_tokens"] = p.config.MaxTokens
	}

	// Marshal request
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, &Error{
			Code:    ErrCodeInvalidConfig,
			Message: "failed to marshal request",
			Err:     err,
		}
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", p.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &Error{
			Code:    ErrCodeAPIError,
			Message: "failed to create HTTP request",
			Err:     err,
		}
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	// Log request
	logger.LogInfo(ctx, "Sending LLM request",
		"provider", p.config.Provider,
		"model", p.config.Model,
		"baseURL", p.baseURL,
		"messagesCount", len(messages),
	)

	// Send request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		logger.LogError(ctx, "LLM API request failed", err,
			"provider", p.config.Provider,
			"model", p.config.Model,
		)
		return nil, &Error{
			Code:    ErrCodeAPIError,
			Message: "HTTP request failed",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{
			Code:    ErrCodeAPIError,
			Message: "failed to read response body",
			Err:     err,
		}
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		logger.LogError(ctx, "LLM API returned non-OK status", nil,
			"provider", p.config.Provider,
			"statusCode", resp.StatusCode,
			"response", string(body),
		)
		return nil, &Error{
			Code:    ErrCodeAPIError,
			Message: fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)),
		}
	}

	// Parse response
	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			QueueTime         float64 `json:"queue_time"`
			PromptTokens      int     `json:"prompt_tokens"`
			PromptTime        float64 `json:"prompt_time"`
			CompletionTokens  int     `json:"completion_tokens"`
			CompletionTime    float64 `json:"completion_time"`
			TotalTokens       int     `json:"total_tokens"`
			TotalTime         float64 `json:"total_time"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.LogError(ctx, "Failed to parse LLM response", err,
			"provider", p.config.Provider,
			"response", string(body),
		)
		return nil, &Error{
			Code:    ErrCodeInvalidResponse,
			Message: "failed to parse API response",
			Err:     err,
		}
	}

	// Validate response
	if len(apiResponse.Choices) == 0 {
		return nil, &Error{
			Code:    ErrCodeInvalidResponse,
			Message: "API returned no choices",
		}
	}

	queueTimeMs := int(apiResponse.Usage.QueueTime * 1000)
	promptTimeMs := int(apiResponse.Usage.PromptTime * 1000)
	completionTimeMs := int(apiResponse.Usage.CompletionTime * 1000)
	totalTimeMs := int(apiResponse.Usage.TotalTime * 1000)

	logger.LogInfo(ctx, "LLM response received",
		"provider", p.config.Provider,
		"model", apiResponse.Model,
		"totalTokens", apiResponse.Usage.TotalTokens,
		"totalTimeMs", totalTimeMs,
	)

	return &GenerateResponse{
		Content:          apiResponse.Choices[0].Message.Content,
		Model:            apiResponse.Model,
		FinishReason:     apiResponse.Choices[0].FinishReason,
		QueueTimeMs:      &queueTimeMs,
		PromptTokens:     &apiResponse.Usage.PromptTokens,
		PromptTimeMs:     &promptTimeMs,
		CompletionTokens: &apiResponse.Usage.CompletionTokens,
		CompletionTimeMs: &completionTimeMs,
		TotalTokens:      &apiResponse.Usage.TotalTokens,
		TotalTimeMs:      &totalTimeMs,
	}, nil
}

// GetProviderName returns the provider name
func (p *OpenAICompatibleProvider) GetProviderName() string {
	return p.config.Provider
}

// IsAvailable checks if the provider is configured
func (p *OpenAICompatibleProvider) IsAvailable() bool {
	return p.config.APIKey != "" && p.config.BaseURL != ""
}
