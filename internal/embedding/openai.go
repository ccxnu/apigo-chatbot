package embedding

import (
	"context"
	"fmt"

	"api-chatbot/domain"
)

type openAIRequest struct {
	Input          any    `json:"input"` // Supports string (single) or []string (batch)
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format,omitempty"` // Added based on curl request
}

type usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type data struct {
	Object    string    `json:"object"` // "embedding"
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type openAIResponse struct {
	Object string `json:"object"`
	Model  string `json:"model"`
	Data   []data `json:"data"`
	Usage  usage  `json:"usage"`
}

type OpenAIEmbeddingService struct {
	paramCache domain.ParameterCache
	httpClient domain.HTTPClient
}

func NewOpenAIEmbeddingService(paramCache domain.ParameterCache, httpClient domain.HTTPClient) *OpenAIEmbeddingService {
	return &OpenAIEmbeddingService{
		paramCache: paramCache,
		httpClient: httpClient,
	}
}

// getConfig retrieves necessary configuration from the parameter cache.
func (s *OpenAIEmbeddingService) getConfig() (apiURL, apiKey, model string, err error) {
	if param, exists := s.paramCache.Get("EMBEDDING_CONFIG"); exists {
		if data, mapErr := param.GetDataAsMap(); mapErr == nil {
			apiURL, _ = data["openaiUrl"].(string)
			apiKey, _ = data["openaiApiKey"].(string)
			model, _ = data["openaiModel"].(string)
		}
	}

	if apiURL == "" || apiKey == "" || model == "" {
		err = fmt.Errorf("OpenAI embedding configuration not found in parameters (apiURL: %s, apiKey: [hidden], model: %s)", apiURL, model)
	}

	return
}

// GenerateEmbedding generates an embedding vector from a single text string.
func (s *OpenAIEmbeddingService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	apiURL, apiKey, model, err := s.getConfig()

	if err != nil {
		return nil, err
	}

	reqBody := openAIRequest{
		Input:          text, // single string
		Model:          model,
		EncodingFormat: "float", // Explicitly set encoding format
	}

	// Create HTTP request
	httpReq := domain.HTTPRequest{
		URL:    apiURL,
		Method: "POST",
		Body:   reqBody,
		AdditionalHeaders: []domain.HTTPHeader{
			{Key: "Authorization", Value: "Bearer " + apiKey},
		},
	}

	// Execute request
	var openAIResp openAIResponse
	err = s.httpClient.Do(ctx, httpReq, &openAIResp)

	if err != nil {
		return nil, fmt.Errorf("error calling OpenAI API for single embedding: %w", err)
	}

	if len(openAIResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data in response")
	}

	// Since it's a single request, we return the first embedding.
	return openAIResp.Data[0].Embedding, nil
}

// GenerateEmbeddings generates embeddings for multiple texts using a single batch API call.
func (s *OpenAIEmbeddingService) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	const maxBatchSize = 2048 // Standard maximum for V3 models
	if len(texts) > maxBatchSize {
		return nil, fmt.Errorf("batch size of %d exceeds the OpenAI maximum of %d texts per request. Please chunk your input", len(texts), maxBatchSize)
	}

    var filteredTexts []string
    for _, text := range texts {
        if text != "" {
            filteredTexts = append(filteredTexts, text)
        }
    }

    if len(filteredTexts) == 0 {
        return nil, nil // All inputs were empty
    }

    textsToSend := filteredTexts

	apiURL, apiKey, model, err := s.getConfig()
	if err != nil {
		fmt.Printf("ERROR: Failed to generate embeddings: %v\n", err)
		return nil, err
	}

	reqBody := openAIRequest{
		Input:          textsToSend, // Use the filtered list
		Model:          model,
		EncodingFormat: "float",
	}

	httpReq := domain.HTTPRequest{
		URL:    apiURL,
		Method: "POST",
		Body:   reqBody,
		AdditionalHeaders: []domain.HTTPHeader{
			{Key: "Authorization", Value: "Bearer " + apiKey},
		},
	}

	var openAIResp openAIResponse
	err = s.httpClient.Do(ctx, httpReq, &openAIResp)

	if err != nil {
		fmt.Printf("ERROR: Failed to calling embeddings: %v\n", err)
		return nil, fmt.Errorf("error calling OpenAI API for batch embeddings: %w", err)
	}

	if len(openAIResp.Data) != len(textsToSend) {
		return nil, fmt.Errorf("OpenAI API returned %d embeddings, expected %d based on filtered input", len(openAIResp.Data), len(textsToSend))
	}

	embeddings := make([][]float32, len(textsToSend))

	for _, item := range openAIResp.Data {
		if item.Index >= 0 && item.Index < len(textsToSend) {
			embeddings[item.Index] = item.Embedding
		} else {
            return nil, fmt.Errorf("received embedding with out-of-bounds index: %d", item.Index)
        }
	}

	return embeddings, nil
}
