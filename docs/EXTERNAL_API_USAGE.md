# External API Usage Guide

This guide explains how to use the External API endpoints with event-based document filtering for the InDTec React frontend.

## Overview

The External API provides three main endpoints:
1. **POST /v1/chat/completions** - Chat with RAG (event-filtered)
2. **POST /v1/embeddings** - Generate embeddings
3. **POST /v1/search** - Search knowledge base (event-filtered)

All requests must include the base request structure from the Next.js frontend (`withBody` helper).

## Event Filtering

Documents are organized by event categories (stored in `doc_category` field):
- `EVENT_INDTEC` - InDTec conference documents
- `EVENT_CONGRESO` - Congreso event documents
- `EVENT_TARIFA` - Tarifa event documents
- `EVENT_TRABAJOS` - Trabajos event documents

The API supports filtering by event keywords. You can pass partial matches (e.g., "INDTEC" matches "EVENT_INDTEC").

## Authentication

All requests require the base structure with session/device tracking:

```typescript
import { withBody } from '@/lib/body';

const requestBody = withBody(
  {
    // Your API-specific payload
  },
  'chatbot-message',      // process
  'device-id-12345',      // idDevice
  '192.168.1.1'          // deviceAddress (IP)
);
```

## 1. Chat Completions (RAG-enabled)

**Endpoint:** `POST /v1/chat/completions`

**Request:**
```typescript
const payload = withBody(
  {
    model: "rag-default",
    messages: [
      {
        role: "user",
        content: "¿Cuál es el objetivo de InDTec?"
      }
    ],
    temperature: 0.7,
    max_tokens: 1000,
    rag_config: {
      enabled: true,
      search_limit: 5,
      min_similarity: 0.7,
      keyword_weight: 0.3,
      event_filter: ["INDTEC"]  // Filter only InDTec documents
    }
  },
  'chatbot-indtec',
  deviceId,
  ipAddress
);

const response = await fetch(`${API_ENDPOINT}/v1/chat/completions`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-App-Authorization': `X-Auth ${APP_AUTHORIZATION}`,
    'Authorization': `Bearer ${API_TOKEN}`,
  },
  body: JSON.stringify(payload)
});
```

**Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "",
  "data": {
    "id": "chatcmpl-abc123",
    "object": "chat.completion",
    "created": 1677652288,
    "model": "rag-default",
    "choices": [
      {
        "index": 0,
        "message": {
          "role": "assistant",
          "content": "El objetivo de InDTec es..."
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
          "document_title": "InDTec 2025 - Objetivos",
          "chunk_id": 42,
          "similarity": 0.89
        }
      ]
    }
  }
}
```

## 2. Embeddings

**Endpoint:** `POST /v1/embeddings`

**Request:**
```typescript
const payload = withBody(
  {
    input: "¿Cuáles son las modalidades de InDTec?",
    model: "text-embedding-3-small"
  },
  'embedding-generation',
  deviceId,
  ipAddress
);
```

**Response:**
```json
{
  "success": true,
  "code": "OK",
  "data": {
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
}
```

## 3. Knowledge Base Search

**Endpoint:** `POST /v1/search`

**Request:**
```typescript
const payload = withBody(
  {
    query: "objetivos específicos",
    limit: 10,
    min_similarity: 0.7,
    search_type: "hybrid",
    keyword_weight: 0.3,
    event_filter: ["INDTEC", "CONGRESO"],  // Search across multiple events
    include_content: true
  },
  'knowledge-search',
  deviceId,
  ipAddress
);
```

**Response:**
```json
{
  "success": true,
  "code": "OK",
  "data": {
    "results": [
      {
        "chunk_id": 42,
        "document_id": 1,
        "document_title": "InDTec - Objetivos Específicos",
        "content": "Los objetivos específicos de InDTec incluyen...",
        "similarity_score": 0.89,
        "keyword_score": 0.65,
        "combined_score": 0.82
      }
    ],
    "total": 1,
    "search_type": "hybrid"
  }
}
```

## Event-Specific Examples

### InDTec Frontend (`/indtec/page.tsx`)

```typescript
const chatWithInDTec = async (message: string, history: any[]) => {
  const payload = withBody(
    {
      model: "rag-default",
      messages: [
        ...history,
        { role: "user", content: message }
      ],
      rag_config: {
        enabled: true,
        search_limit: 5,
        event_filter: ["INDTEC"]  // Only search InDTec documents
      }
    },
    'chatbot-indtec',
    deviceId,
    ipAddress
  );

  const response = await fetch(`${API_ENDPOINT}/v1/chat/completions`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-App-Authorization': `X-Auth ${APP_AUTHORIZATION}`,
      'Authorization': `Bearer ${API_TOKEN}`,
    },
    body: JSON.stringify(payload)
  });

  return await response.json();
};
```

### Congreso Frontend (`/congreso/page.tsx`)

```typescript
const chatWithCongreso = async (message: string) => {
  const payload = withBody(
    {
      model: "rag-default",
      messages: [{ role: "user", content: message }],
      rag_config: {
        enabled: true,
        event_filter: ["CONGRESO"]  // Only search Congreso documents
      }
    },
    'chatbot-congreso',
    deviceId,
    ipAddress
  );

  // ... same fetch logic
};
```

### Multi-Event Search

```typescript
const searchAcrossEvents = async (query: string) => {
  const payload = withBody(
    {
      query,
      limit: 20,
      event_filter: ["INDTEC", "CONGRESO", "TARIFA"],  // Search multiple events
      include_content: true
    },
    'multi-event-search',
    deviceId,
    ipAddress
  );

  const response = await fetch(`${API_ENDPOINT}/v1/search`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-App-Authorization': `X-Auth ${APP_AUTHORIZATION}`,
      'Authorization': `Bearer ${API_TOKEN}`,
    },
    body: JSON.stringify(payload)
  });

  return await response.json();
};
```

## Response Format

All responses follow the standard Result structure:

```typescript
interface Result<T> {
  success: boolean;    // true if operation succeeded
  code: string;        // "OK" or error code (e.g., "ERR_INTERNAL_DB")
  info: string;        // Human-readable message
  data: T;            // Response payload
}
```

## Error Handling

```typescript
const response = await fetch(endpoint, options);
const result = await response.json();

if (!result.success) {
  console.error(`Error ${result.code}: ${result.info}`);
  // Handle specific error codes
  switch (result.code) {
    case 'ERR_EMBEDDING_GENERATION':
      // Retry or fallback
      break;
    case 'ERR_INTERNAL_DB':
      // Show user-friendly error
      break;
  }
}
```

## Best Practices

1. **Always filter by event** - Use `event_filter` to ensure users only see relevant content for their section
2. **Cache device ID** - Store the device ID from `uid.ts` for consistent user tracking
3. **Handle conversation history** - Pass previous messages for context-aware responses
4. **Set reasonable limits** - Use `search_limit: 5` for chat, higher for dedicated search pages
5. **Include content selectively** - Only set `include_content: true` when you need the full text

## Environment Variables

```env
# Backend API endpoint
CHATBOT_API_ENDPOINT=https://your-backend-api.com

# API token for authentication
CHATBOT_API_TOKEN=sk_live_xxxxxxxxxxxxx

# Authorization header value
CHATBOT_API_AUTHORIZATION=wiaAchcHks3rBxIhJQem1nLoMDwdoQ==
```

## Testing

Test the endpoints with curl:

```bash
# Chat completion
curl -X POST https://your-api.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "X-App-Authorization: X-Auth wiaAchcHks3rBxIhJQem1nLoMDwdoQ==" \
  -H "Authorization: Bearer $API_TOKEN" \
  -d '{
    "idSession": "test-session",
    "idRequest": "550e8400-e29b-41d4-a716-446655440000",
    "process": "test-chat",
    "idDevice": "test-device",
    "deviceAddress": "127.0.0.1",
    "dateProcess": "2025-10-30T10:00:00Z",
    "model": "rag-default",
    "messages": [{"role": "user", "content": "Hola"}],
    "rag_config": {
      "enabled": true,
      "event_filter": ["INDTEC"]
    }
  }'
```
