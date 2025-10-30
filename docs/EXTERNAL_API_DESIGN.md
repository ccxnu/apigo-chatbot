# External API Design

## Overview

This document outlines the design for the External API (Phase 4), providing Claude-style API endpoints for external clients to access the chatbot's RAG capabilities programmatically.

## Architecture

```
External Client
    ↓ (API Key in Authorization header)
API Key Auth Middleware
    ↓ (validates key, checks rate limits, IP whitelist)
Rate Limiting Middleware
    ↓
External API Endpoints (/v1/*)
    ↓
Use Cases (reuse existing ChunkUseCase, LLM, etc.)
    ↓
Repository Layer
    ↓
PostgreSQL Database
```

## Authentication Flow

### API Key Format
```
Authorization: Bearer sk_live_abc123...xyz
```

### Validation Steps
1. Extract API key from `Authorization` header
2. Query `cht_api_keys` table for the key (hashed lookup)
3. Validate:
   - Key exists and is active
   - Key hasn't expired
   - Request IP matches allowed IPs (if configured)
   - Endpoint is in allowed permissions
4. Check rate limiting (requests per hour from key config)
5. Update `key_last_used_at`
6. Pass request to endpoint handler

## Database Layer

### New Stored Procedures

```sql
-- sp_create_api_key: Create new API key
-- sp_get_api_key_by_value: Retrieve API key by value (for auth)
-- sp_update_api_key_last_used: Update last used timestamp
-- sp_delete_api_key: Delete/deactivate API key
-- sp_get_all_api_keys: List all API keys (admin endpoint)
-- sp_track_api_usage: Track usage statistics
```

### New Table: cht_api_usage

Track API usage for billing/analytics:

```sql
CREATE TABLE cht_api_usage (
    usg_id SERIAL PRIMARY KEY,
    usg_api_key_id INT REFERENCES cht_api_keys(key_id),
    usg_endpoint VARCHAR(100),
    usg_method VARCHAR(10),
    usg_status_code INT,
    usg_tokens_used INT,
    usg_request_time_ms INT,
    usg_ip_address VARCHAR(45),
    usg_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Domain Layer

### New Interfaces

```go
// domain/api_key.go
type APIKeyRepository interface {
    Create(ctx context.Context, params CreateAPIKeyParams) (*CreateAPIKeyResult, error)
    GetByValue(ctx context.Context, keyValue string) (*APIKey, error)
    UpdateLastUsed(ctx context.Context, keyID int) error
    Delete(ctx context.Context, keyID int) (*DeleteAPIKeyResult, error)
    GetAll(ctx context.Context) ([]APIKey, error)
}

type APIKeyUseCase interface {
    CreateAPIKey(ctx context.Context, name, keyType string, permissions []string, rateLimit int) Result[*APIKey]
    ValidateAPIKey(ctx context.Context, keyValue, ipAddress, endpoint string) Result[*APIKey]
    RevokeAPIKey(ctx context.Context, keyID int) Result[Data]
    ListAPIKeys(ctx context.Context) Result[[]APIKey]
}

type APIUsageRepository interface {
    Track(ctx context.Context, params TrackAPIUsageParams) error
    GetUsageStats(ctx context.Context, keyID int, from, to time.Time) ([]APIUsage, error)
}
```

## API Endpoints

### 1. POST `/v1/chat/completions`

Claude/OpenAI-compatible chat completion with RAG.

**Request:**
```json
{
  "model": "rag-default",
  "messages": [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "What are the enrollment requirements?"}
  ],
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false,
  "rag_config": {
    "enabled": true,
    "search_limit": 5,
    "min_similarity": 0.7,
    "keyword_weight": 0.3
  }
}
```

**Response (non-streaming):**
```json
{
  "id": "chatcmpl-abc123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "rag-default",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Based on the institute's documentation..."
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 56,
    "completion_tokens": 31,
    "total_tokens": 87
  },
  "rag_context": {
    "chunks_retrieved": 3,
    "sources": [
      {
        "document_id": 1,
        "document_title": "Enrollment Guide 2025",
        "chunk_id": 42,
        "similarity": 0.89
      }
    ]
  }
}
```

**Response (streaming):**
```
data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"rag-default","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}

data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"rag-default","choices":[{"index":0,"delta":{"content":"Based"},"finish_reason":null}]}

data: {"id":"chatcmpl-abc123","object":"chat.completion.chunk","created":1677652288,"model":"rag-default","choices":[{"index":0,"delta":{"content":" on"},"finish_reason":null}]}

...

data: [DONE]
```

### 2. POST `/v1/embeddings`

Generate embeddings for text input.

**Request:**
```json
{
  "input": "What are the enrollment requirements?",
  "model": "text-embedding-3-small"
}
```

**Response:**
```json
{
  "object": "list",
  "data": [
    {
      "object": "embedding",
      "embedding": [0.0023064255, -0.009327292, ...],
      "index": 0
    }
  ],
  "model": "text-embedding-3-small",
  "usage": {
    "prompt_tokens": 8,
    "total_tokens": 8
  }
}
```

### 3. POST `/v1/search`

Direct knowledge base search with filters.

**Request:**
```json
{
  "query": "enrollment requirements",
  "limit": 10,
  "min_similarity": 0.7,
  "search_type": "hybrid",
  "keyword_weight": 0.3,
  "filters": {
    "document_ids": [1, 2, 3],
    "date_from": "2025-01-01",
    "date_to": "2025-12-31"
  },
  "include_content": true
}
```

**Response:**
```json
{
  "results": [
    {
      "chunk_id": 42,
      "document_id": 1,
      "document_title": "Enrollment Guide 2025",
      "content": "To enroll, students must...",
      "similarity_score": 0.89,
      "keyword_score": 0.65,
      "combined_score": 0.82,
      "metadata": {
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    }
  ],
  "total": 1,
  "search_type": "hybrid"
}
```

## Middleware Implementation

### API Key Auth Middleware

```go
// api/middleware/api_key_auth.go
func APIKeyAuth(apiKeyUseCase domain.APIKeyUseCase) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract API key from Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                writeError(w, 401, "Missing API key")
                return
            }

            // Parse "Bearer sk_xxx"
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                writeError(w, 401, "Invalid Authorization header format")
                return
            }

            apiKey := parts[1]
            ipAddress := getClientIP(r)
            endpoint := r.URL.Path

            // Validate API key
            result := apiKeyUseCase.ValidateAPIKey(r.Context(), apiKey, ipAddress, endpoint)
            if !result.Success {
                writeError(w, 401, result.Message)
                return
            }

            // Store API key in context for downstream use
            ctx := context.WithValue(r.Context(), "api_key", result.Data)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Rate Limiting Middleware

```go
// api/middleware/rate_limiter.go
func RateLimiter() func(http.Handler) http.Handler {
    // Use in-memory rate limiter (e.g., golang.org/x/time/rate)
    // Key: API key ID
    // Limit: From api_key.rate_limit (requests per hour)

    limiters := sync.Map{}

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            apiKey := r.Context().Value("api_key").(*domain.APIKey)

            // Get or create limiter for this API key
            limiterKey := fmt.Sprintf("api_key_%d", apiKey.ID)

            limiter, _ := limiters.LoadOrStore(limiterKey, rate.NewLimiter(
                rate.Limit(apiKey.RateLimit/3600.0), // per second
                apiKey.RateLimit/60, // burst (per minute)
            ))

            if !limiter.(*rate.Limiter).Allow() {
                writeError(w, 429, "Rate limit exceeded")
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

## Implementation Plan

### Phase 1: Database & Repository Layer
1. Create migration file for API key procedures and usage tracking table
2. Implement APIKeyRepository
3. Implement APIUsageRepository
4. Write unit tests for repositories

### Phase 2: Use Case Layer
1. Implement APIKeyUseCase
2. Add key generation with secure random values
3. Add key validation logic
4. Write unit tests for use cases

### Phase 3: Middleware
1. Implement API key authentication middleware
2. Implement rate limiting middleware
3. Implement usage tracking middleware
4. Test middleware integration

### Phase 4: External API Endpoints
1. Create external API router (`api/route/external_api_router.go`)
2. Implement `/v1/chat/completions` (non-streaming)
3. Implement `/v1/embeddings`
4. Implement `/v1/search`
5. Add OpenAPI documentation for external endpoints

### Phase 5: Streaming Support
1. Implement SSE (Server-Sent Events) for streaming
2. Update `/v1/chat/completions` to support `stream: true`
3. Test streaming functionality

### Phase 6: Testing & Documentation
1. End-to-end testing with real API keys
2. Load testing for rate limiting
3. Update API documentation
4. Create client SDK examples (Python, JavaScript, curl)

## Security Considerations

1. **API Key Storage**: Hash API keys in database (like passwords)
2. **HTTPS Only**: Enforce HTTPS in production
3. **IP Whitelisting**: Optional but recommended for sensitive data
4. **Rate Limiting**: Prevent abuse and DoS attacks
5. **Audit Logging**: Track all API usage in `cht_api_usage`
6. **Key Rotation**: Support key expiration and rotation
7. **Permissions**: Fine-grained permissions per endpoint

## Usage Tracking & Analytics

Track the following metrics:
- Total requests per API key
- Tokens used (for billing)
- Average response time
- Error rate
- Most used endpoints
- Peak usage times

## API Key Management (Admin Endpoints)

```
POST   /admin/api-keys          - Create new API key
GET    /admin/api-keys          - List all API keys
GET    /admin/api-keys/:id      - Get API key details
PUT    /admin/api-keys/:id      - Update API key (permissions, rate limit)
DELETE /admin/api-keys/:id      - Revoke API key
GET    /admin/api-keys/:id/usage - Get usage statistics
```

## Client SDK Example (Python)

```python
import requests

class ChatbotClient:
    def __init__(self, api_key: str, base_url: str = "https://api.example.com"):
        self.api_key = api_key
        self.base_url = base_url
        self.headers = {
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json"
        }

    def chat(self, messages: list, rag_enabled: bool = True, stream: bool = False):
        response = requests.post(
            f"{self.base_url}/v1/chat/completions",
            json={
                "model": "rag-default",
                "messages": messages,
                "stream": stream,
                "rag_config": {"enabled": rag_enabled}
            },
            headers=self.headers,
            stream=stream
        )

        if stream:
            for line in response.iter_lines():
                if line:
                    yield line.decode('utf-8')
        else:
            return response.json()

    def search(self, query: str, limit: int = 10):
        response = requests.post(
            f"{self.base_url}/v1/search",
            json={"query": query, "limit": limit},
            headers=self.headers
        )
        return response.json()

# Usage
client = ChatbotClient("sk_live_abc123")
response = client.chat([
    {"role": "user", "content": "What are the enrollment requirements?"}
])
print(response["choices"][0]["message"]["content"])
```

## Performance Considerations

1. **Caching**: Cache API keys in memory (invalidate on update)
2. **Connection Pooling**: Reuse database connections
3. **Async Processing**: Use goroutines for usage tracking
4. **Response Compression**: Enable gzip compression
5. **CDN**: Serve static documentation via CDN

## Monitoring & Alerts

Monitor:
- API key authentication failures (potential security issues)
- Rate limit hits (may need to adjust limits)
- Error rates per endpoint
- Response time percentiles (p50, p95, p99)
- Token usage trends

Set alerts for:
- High error rates (> 5%)
- Slow responses (> 2s p95)
- Unusual traffic patterns
- Failed authentication attempts (> 10/min from same IP)
