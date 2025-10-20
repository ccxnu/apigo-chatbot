# API Endpoints Documentation

This document lists all available API endpoints in the ISTS Chatbot API.

## Access Documentation

- **Interactive UI**: http://localhost:8080/docs
- **OpenAPI JSON**: http://localhost:8080/openapi.json
- **OpenAPI YAML**: http://localhost:8080/openapi.yaml

## Endpoint Overview

Total Endpoints: **31** (26 Huma + 5 system)

### By Module:
- **System**: 1 endpoint
- **Parameters**: 6 endpoints
- **Documents**: 7 endpoints
- **Chunks**: 7 endpoints
- **Chunk Statistics**: 5 endpoints

---

## System Endpoints

### Health Check
**GET** `/health`
- **Description**: Returns the health status of the API
- **OpenAPI**: ✅ Documented
- **Response**: Always returns `Success: true, Code: OK`

---

## Parameters Module (6 endpoints)

All parameter endpoints use Huma with OpenAPI 3.1 documentation.

### 1. Get All Parameters
**POST** `/huma/parameters/get-all`
- **Description**: Retrieves all active system parameters from cache or database
- **Input**: None
- **Response**: Array of parameters
- **Success Code**: `OK`

### 2. Get Parameter by Code
**POST** `/huma/parameters/get-by-code`
- **Description**: Retrieves a specific parameter by its unique code
- **Input**: `{ "code": "PARAMETER_CODE" }`
- **Response**: Single parameter object
- **Success Code**: `OK`
- **Error Codes**: `ERR_PARAM_NOT_FOUND`

### 3. Add Parameter
**POST** `/huma/parameters/add`
- **Description**: Creates a new system parameter with validation
- **Input**:
  ```json
  {
    "name": "Parameter Name",
    "code": "PARAM_CODE",
    "data": { "key": "value" },
    "description": "Parameter description"
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_PARAM_CODE_EXISTS`

### 4. Update Parameter
**POST** `/huma/parameters/update`
- **Description**: Updates an existing parameter
- **Input**:
  ```json
  {
    "code": "PARAM_CODE",
    "name": "Updated Name",
    "data": { "key": "new_value" },
    "description": "Updated description"
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_PARAM_NOT_FOUND`

### 5. Delete Parameter
**POST** `/huma/parameters/delete`
- **Description**: Soft deletes a parameter (sets active = false)
- **Input**: `{ "code": "PARAM_CODE" }`
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_PARAM_NOT_FOUND`

### 6. Reload Cache (Legacy)
**POST** `/parameters/reload-cache`
- **Description**: Reloads parameter cache from database
- **OpenAPI**: ❌ Not documented (legacy endpoint)

---

## Documents Module (7 endpoints)

All document endpoints use Huma with OpenAPI 3.1 documentation.

### 1. Get All Documents
**POST** `/huma/documents/get-all`
- **Description**: Retrieves all active documents with pagination
- **Input**:
  ```json
  {
    "limit": 100,
    "offset": 0
  }
  ```
- **Response**: Array of documents
- **Success Code**: `OK`

### 2. Get Document by ID
**POST** `/huma/documents/get-by-id`
- **Description**: Retrieves a specific document by its ID
- **Input**: `{ "docId": 1 }`
- **Response**: Single document object
- **Success Code**: `OK`
- **Error Codes**: `ERR_DOCUMENT_NOT_FOUND`

### 3. Get Documents by Category
**POST** `/huma/documents/get-by-category`
- **Description**: Retrieves documents filtered by category with pagination
- **Input**:
  ```json
  {
    "category": "ACADEMIC",
    "limit": 100,
    "offset": 0
  }
  ```
- **Response**: Array of documents
- **Success Code**: `OK`

### 4. Search Documents by Title
**POST** `/huma/documents/search-by-title`
- **Description**: Searches documents by title pattern (case-insensitive)
- **Input**:
  ```json
  {
    "titlePattern": "handbook",
    "limit": 100
  }
  ```
- **Response**: Array of matching documents
- **Success Code**: `OK`

### 5. Create Document
**POST** `/huma/documents/create`
- **Description**: Creates a new document in the knowledge base
- **Input**:
  ```json
  {
    "category": "ACADEMIC",
    "title": "Student Handbook",
    "summary": "Complete guide for students",
    "source": "https://example.com/handbook.pdf",
    "publishedAt": "2025-01-19T10:00:00Z"
  }
  ```
- **Response**: `{ "docId": 123 }`
- **Success Code**: `OK`
- **Error Codes**: `ERR_REQUIRED_FIELDS`

### 6. Update Document
**POST** `/huma/documents/update`
- **Description**: Updates an existing document
- **Input**:
  ```json
  {
    "docId": 1,
    "category": "ACADEMIC",
    "title": "Updated Title",
    "summary": "Updated summary",
    "source": "https://example.com/new.pdf",
    "publishedAt": "2025-01-20T10:00:00Z"
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_DOCUMENT_NOT_FOUND`

### 7. Delete Document
**POST** `/huma/documents/delete`
- **Description**: Soft deletes a document (sets active = false)
- **Input**: `{ "docId": 1 }`
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_DOCUMENT_NOT_FOUND`

---

## Chunks Module (7 endpoints)

All chunk endpoints use Huma with OpenAPI 3.1 documentation.

### 1. Get Chunks by Document
**POST** `/huma/chunks/get-by-document`
- **Description**: Retrieves all chunks for a specific document
- **Input**: `{ "docId": 1 }`
- **Response**: Array of chunks
- **Success Code**: `OK`

### 2. Get Chunk by ID
**POST** `/huma/chunks/get-by-id`
- **Description**: Retrieves a specific chunk by its ID including embedding
- **Input**: `{ "chunkId": 1 }`
- **Response**: Single chunk object with embedding
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_NOT_FOUND`

### 3. Vector Similarity Search 🔍
**POST** `/huma/chunks/similarity-search`
- **Description**: Performs semantic search using vector embeddings for RAG
- **Tags**: `RAG`
- **Input**:
  ```json
  {
    "queryEmbedding": [0.1, 0.2, ...],  // 1536 dimensions
    "limit": 5,
    "minSimilarity": 0.7
  }
  ```
- **Response**: Array of chunks with similarity scores, ordered by relevance
- **Success Code**: `OK`
- **Use Case**: Main RAG retrieval endpoint

### 4. Create Chunk
**POST** `/huma/chunks/create`
- **Description**: Creates a new chunk with optional embedding
- **Input**:
  ```json
  {
    "documentId": 1,
    "content": "Chunk text content",
    "embedding": [0.1, 0.2, ...]  // Optional, 1536 dimensions
  }
  ```
- **Response**: `{ "chunkId": 456 }`
- **Success Code**: `OK`
- **Error Codes**: `ERR_DOCUMENT_NOT_FOUND`
- **Note**: Automatically initializes statistics record

### 5. Update Chunk Embedding
**POST** `/huma/chunks/update-embedding`
- **Description**: Updates the embedding vector for a chunk
- **Input**:
  ```json
  {
    "chunkId": 1,
    "embedding": [0.1, 0.2, ...]  // 1536 dimensions
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_NOT_FOUND`

### 6. Delete Chunk
**POST** `/huma/chunks/delete`
- **Description**: Hard deletes a chunk (cascades to statistics)
- **Input**: `{ "chunkId": 1 }`
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_NOT_FOUND`

### 7. Bulk Create Chunks
**POST** `/huma/chunks/bulk-create`
- **Description**: Creates multiple chunks for a document in a single operation
- **Input**:
  ```json
  {
    "documentId": 1,
    "contents": ["Chunk 1", "Chunk 2", "Chunk 3"],
    "embeddings": [[...], [...], [...]]  // Optional
  }
  ```
- **Response**: `{ "chunksCreated": 3 }`
- **Success Code**: `OK`
- **Use Case**: Efficient for document ingestion

---

## Chunk Statistics Module (5 endpoints)

All statistics endpoints use Huma with OpenAPI 3.1 documentation.

### 1. Get Chunk Statistics
**POST** `/huma/chunk-statistics/get-by-chunk`
- **Description**: Retrieves all statistics and quality metrics for a specific chunk
- **Input**: `{ "chunkId": 1 }`
- **Response**: Full statistics object with all metrics
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_STATS_NOT_FOUND`

### 2. Get Top Chunks by Usage 📊
**POST** `/huma/chunk-statistics/get-top-by-usage`
- **Description**: Retrieves the most frequently used chunks for analytics
- **Tags**: `Analytics`
- **Input**: `{ "limit": 10 }`
- **Response**: Array of top chunks with usage counts
- **Success Code**: `OK`
- **Use Case**: Identify popular knowledge

### 3. Increment Chunk Usage
**POST** `/huma/chunk-statistics/increment-usage`
- **Description**: Increments usage count and updates last used timestamp
- **Input**: `{ "chunkId": 1 }`
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_STATS_NOT_FOUND`
- **Use Case**: Call this when a chunk is used in RAG responses

### 4. Update RAG Quality Metrics 📈
**POST** `/huma/chunk-statistics/update-quality-metrics`
- **Description**: Updates quality and relevance metrics for a chunk
- **Tags**: `RAG`
- **Input**:
  ```json
  {
    "chunkId": 1,
    "precisionAtK": 0.85,
    "recallAtK": 0.90,
    "f1AtK": 0.87,
    "mrr": 0.95,
    "map": 0.88,
    "ndcg": 0.92
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_STATS_NOT_FOUND`
- **Note**: Only provided metrics are updated

**Metric Definitions:**
- **Precision@K**: Accuracy of top K results
- **Recall@K**: Coverage of relevant items in top K
- **F1@K**: Harmonic mean of Precision and Recall
- **MRR**: Mean Reciprocal Rank
- **MAP**: Mean Average Precision
- **NDCG**: Normalized Discounted Cumulative Gain

### 5. Update Chunk Staleness
**POST** `/huma/chunk-statistics/update-staleness`
- **Description**: Updates staleness metric to track content freshness
- **Input**:
  ```json
  {
    "chunkId": 1,
    "stalenessDays": 30
  }
  ```
- **Response**: Success message
- **Success Code**: `OK`
- **Error Codes**: `ERR_CHUNK_STATS_NOT_FOUND`

---

## Legacy Endpoints (Original Chi Routes)

These endpoints continue to work but are **not documented** in OpenAPI:

- `POST /parameters/get-all`
- `POST /parameters/get-by-code`
- `POST /parameters/add`
- `POST /parameters/update`
- `POST /parameters/delete`
- `POST /parameters/reload-cache`
- `POST /documents/*` (all 7 endpoints)
- `POST /chunks/*` (all 7 endpoints)
- `POST /chunk-statistics/*` (all 5 endpoints)

**Note**: These are functionally identical to Huma endpoints but use the `/path` instead of `/huma/path`.

---

## Response Format

All endpoints return responses in this format:

```json
{
  "success": true,           // or false for errors
  "code": "OK",             // or "ERR_..." for errors
  "info": "message",        // from parameter cache
  "data": { ... }           // actual response data
}
```

See `API_RESPONSE_FORMAT.md` for detailed examples.

---

## Common Error Codes

- `OK` - Success
- `ERR_INTERNAL_DB` - Database error
- `ERR_INTERNAL_SERVER` - Server error
- `ERR_VALIDATION_BODY` - Request validation failed
- `ERR_PARAM_NOT_FOUND` - Parameter not found
- `ERR_PARAM_CODE_EXISTS` - Parameter code already exists
- `ERR_DOCUMENT_NOT_FOUND` - Document not found
- `ERR_REQUIRED_FIELDS` - Required fields missing
- `ERR_CHUNK_NOT_FOUND` - Chunk not found
- `ERR_CHUNK_STATS_NOT_FOUND` - Chunk statistics not found

---

## Testing Examples

### Health Check
```bash
curl http://localhost:8080/health
```

### Get All Parameters
```bash
curl -X POST http://localhost:8080/huma/parameters/get-all \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Create Document
```bash
curl -X POST http://localhost:8080/huma/documents/create \
  -H "Content-Type: application/json" \
  -d '{
    "category": "ACADEMIC",
    "title": "Test Document",
    "summary": "Test summary"
  }'
```

### Similarity Search (RAG)
```bash
curl -X POST http://localhost:8080/huma/chunks/similarity-search \
  -H "Content-Type: application/json" \
  -d '{
    "queryEmbedding": [0.1, 0.2, ...],
    "limit": 5,
    "minSimilarity": 0.7
  }'
```

---

## Interactive Testing

Visit **http://localhost:8080/docs** for interactive API documentation where you can:
- ✅ Browse all endpoints
- ✅ View request/response schemas
- ✅ Try endpoints directly from browser
- ✅ See examples and validation rules
- ✅ Download OpenAPI spec

---

## Summary

- **Total Endpoints**: 31
- **OpenAPI Documented**: 26 (all Huma routes)
- **Legacy Endpoints**: 19 (still functional)
- **Response Format**: Consistent `common.Result[T]` across all endpoints
- **Error Handling**: Codes from parameter cache
- **Documentation**: Auto-generated from Go code
