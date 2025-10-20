# API Response Format Standard

All API endpoints in this project use a consistent response format defined in `api/common/result.go`.

## Standard Response Structure

```go
type Result[T any] struct {
    Success bool   `json:"success"` // Operation status
    Code    string `json:"code"`    // Response code
    Info    string `json:"info"`    // Human-readable message
    Data    T      `json:"data"`    // Actual response data
}
```

## Response Codes

### Success Response
- **Success**: `true`
- **Code**: `"OK"`
- **Info**: Message from parameter cache (e.g., "Operation successful")
- **Data**: Actual response data

### Error Response
- **Success**: `false`
- **Code**: Error code with `"ERR_"` prefix (e.g., `"ERR_PARAM_NOT_FOUND"`)
- **Info**: Error message from parameter cache
- **Data**: Usually `null` or empty object

## Error Code Naming Convention

All error codes follow the pattern: `ERR_[MODULE]_[DESCRIPTION]`

### Common Error Codes

**System Errors:**
- `ERR_INTERNAL_DB` - Database error
- `ERR_INTERNAL_SERVER` - Server error
- `ERR_VALIDATION_BODY` - Request validation failed

**Parameter Errors:**
- `ERR_PARAM_NOT_FOUND` - Parameter not found
- `ERR_PARAM_CODE_EXISTS` - Parameter code already exists

**Document Errors:**
- `ERR_DOCUMENT_NOT_FOUND` - Document not found
- `ERR_REQUIRED_FIELDS` - Required fields missing

**Chunk Errors:**
- `ERR_CHUNK_NOT_FOUND` - Chunk not found
- `ERR_CHUNK_STATS_NOT_FOUND` - Chunk statistics not found

## Examples

### Success Response Example

**Request:**
```bash
POST /parameters/get-by-code
{
  "code": "SYSTEM_NAME"
}
```

**Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Operation successful",
  "data": {
    "id": 1,
    "name": "System Name",
    "code": "SYSTEM_NAME",
    "data": {
      "value": "ISTS Chatbot"
    },
    "description": "The name of the system",
    "active": true,
    "createdAt": "2025-01-15T10:00:00Z",
    "updatedAt": "2025-01-15T10:00:00Z"
  }
}
```

### Error Response Example

**Request:**
```bash
POST /parameters/get-by-code
{
  "code": "INVALID_CODE"
}
```

**Response:**
```json
{
  "success": false,
  "code": "ERR_PARAM_NOT_FOUND",
  "info": "Parameter not found in the system",
  "data": null
}
```

### List Response Example

**Request:**
```bash
POST /documents/get-all
{
  "limit": 10,
  "offset": 0
}
```

**Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Documents retrieved successfully",
  "data": [
    {
      "id": 1,
      "category": "ACADEMIC",
      "title": "Student Handbook",
      "summary": "Complete guide for students",
      "source": "https://example.com/handbook.pdf",
      "publishedAt": "2025-01-01T00:00:00Z",
      "active": true,
      "createdAt": "2025-01-15T10:00:00Z",
      "updatedAt": "2025-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "category": "ACADEMIC",
      "title": "Course Catalog",
      "summary": "Available courses for 2025",
      "source": null,
      "publishedAt": "2025-01-02T00:00:00Z",
      "active": true,
      "createdAt": "2025-01-16T10:00:00Z",
      "updatedAt": "2025-01-16T10:00:00Z"
    }
  ]
}
```

### Create Response Example

**Request:**
```bash
POST /documents/create
{
  "category": "ACADEMIC",
  "title": "New Document",
  "summary": "Document summary",
  "source": "https://example.com/doc.pdf",
  "publishedAt": "2025-01-19T10:00:00Z"
}
```

**Success Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Document created successfully",
  "data": {
    "docId": 123
  }
}
```

**Error Response (Validation Failed):**
```json
{
  "success": false,
  "code": "ERR_REQUIRED_FIELDS",
  "info": "Category and title are required fields",
  "data": null
}
```

## Error Message Source

All error messages come from the **parameter cache**:

1. Error codes are stored as parameters in `cht_parameters` table
2. Each error code has a `data.message` field with the human-readable message
3. Use cases retrieve messages from parameter cache using `getErrorMessage(errorCode)`
4. This allows changing error messages without code deployment

**Example Parameter:**
```sql
INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
VALUES (
  'Error: Parameter Not Found',
  'ERR_PARAM_NOT_FOUND',
  '{"message": "Parameter not found in the system"}',
  'Error message for parameter not found'
);
```

## Huma Integration

All Huma endpoints use the same `common.Result[T]` type:

```go
type GetParameterResponse struct {
    Body common.Result[*domain.Parameter]
}

func handler(ctx context.Context, input *GetParameterInput) (*GetParameterResponse, error) {
    result := paramUseCase.GetByCode(ctx, input.Body.Code)

    resp := &GetParameterResponse{}
    resp.Body = result // Already formatted as common.Result

    return resp, nil
}
```

## Benefits

✅ **Consistency** - All endpoints use the same format
✅ **Type Safety** - Generic `Result[T]` ensures correct data types
✅ **Flexibility** - Error messages can be updated without code changes
✅ **Client-Friendly** - Always know structure of response
✅ **Documentation** - OpenAPI spec reflects actual response format
✅ **Error Handling** - Easy to check `success` field and handle errors

## Client Implementation Example

### JavaScript/TypeScript

```typescript
interface ApiResponse<T> {
  success: boolean;
  code: string;
  info: string;
  data: T;
}

async function getParameter(code: string): Promise<Parameter | null> {
  const response = await fetch('/parameters/get-by-code', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code })
  });

  const result: ApiResponse<Parameter> = await response.json();

  if (result.success) {
    return result.data;
  } else {
    console.error(`Error ${result.code}: ${result.info}`);
    return null;
  }
}
```

### Python

```python
from typing import Generic, TypeVar, Optional
from dataclasses import dataclass

T = TypeVar('T')

@dataclass
class ApiResponse(Generic[T]):
    success: bool
    code: str
    info: str
    data: Optional[T]

def get_parameter(code: str) -> Optional[dict]:
    response = requests.post('/parameters/get-by-code', json={'code': code})
    result = response.json()

    if result['success']:
        return result['data']
    else:
        print(f"Error {result['code']}: {result['info']}")
        return None
```

## Testing Responses

```bash
# Success case
curl -X POST http://localhost:8080/health | jq
{
  "success": true,
  "code": "OK",
  "info": "Service is healthy",
  "data": {
    "status": "ok",
    "timestamp": "2025-01-19T10:30:00Z",
    "version": "1.0.0"
  }
}

# Error case
curl -X POST http://localhost:8080/parameters/get-by-code \
  -H "Content-Type: application/json" \
  -d '{"code": "INVALID"}' | jq
{
  "success": false,
  "code": "ERR_PARAM_NOT_FOUND",
  "info": "Parameter not found in the system",
  "data": null
}
```

## Summary

- ✅ Use `common.Result[T]` for ALL responses
- ✅ Success always has `Code: "OK"`
- ✅ Errors always have `Code: "ERR_..."`
- ✅ Messages come from parameter cache
- ✅ Consistent across original and Huma endpoints
