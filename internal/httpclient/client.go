package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-chatbot/domain"
)

type httpClient struct {
	client     *http.Client
	paramCache domain.ParameterCache
}

func NewHTTPClient(paramCache domain.ParameterCache) domain.HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		paramCache: paramCache,
	}
}

// addHeaders prepares the http.Header map based on the HTTPRequest
func (c *httpClient) addHeaders(req domain.HTTPRequest) http.Header {
	headers := make(http.Header)

	if req.Body != nil {
		headers.Set("Content-Type", "application/json")
	}

	if req.AuthKey != "" && req.AuthValue != "" {
		headers.Set(req.AuthKey, req.AuthValue)
	}

	for _, header := range req.AdditionalHeaders {
		if header.Value != "" {
			headers.Set(header.Key, header.Value)
		}
	}

	return headers
}

// Do executes an HTTP request and returns a Result with the response
func (c *httpClient) Do(ctx context.Context, req domain.HTTPRequest, result any) error {
	var body io.Reader
	var jsonData []byte

	// Prepare request body (JSON marshaling)
	if req.Body != nil {
		var err error
		jsonData, err = json.Marshal(req.Body)
		if err != nil {
			err = fmt.Errorf("Failed to marshal request data: %w", err)
			return err
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Create the http.Request with context
	request, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
	if err != nil {
		err = fmt.Errorf("Failed to create request: %w", err)
		return err
	}

	request.Header = c.addHeaders(req)

	// Debug log disabled - uncomment to debug HTTP requests
	/*
	if len(jsonData) > 0 && req.URL != "" {
		bodyPreview := string(jsonData)
		if len(bodyPreview) > 1000 {
			bodyPreview = bodyPreview[:1000] + "..."
		}
		fmt.Printf("\n=== HTTP Request to %s ===\n%s\n===\n\n", req.URL, bodyPreview)
	}
	*/

	response, err := c.client.Do(request)

	if err != nil {
		err = fmt.Errorf("External service request failed: %w", err)
		return err
	}

	defer response.Body.Close()

	// Check for HTTP errors
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		errBody, _ := io.ReadAll(response.Body)
		err = fmt.Errorf("Error in external service: status %d. Message: %s", response.StatusCode, string(errBody))
		return err
	}

	// Decode the response body
	if result != nil {
		if err := json.NewDecoder(response.Body).Decode(result); err != nil {
			err = fmt.Errorf("Failed to decode response JSON into result structure: %w", err)
			return err
		}
	}

	return nil
}
