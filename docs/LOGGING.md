# Logging Guide

This application uses Go's `slog` package for structured logging with automatic request tracking.

## Log Format

All logs use JSON format with these standard fields:
- `time`: ISO 8601 timestamp
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `msg`: Short description of the event
- `request_id`: UUID tracking the request across all layers (when available)

## Environment-Based Output

Logging output is automatically configured based on the `APP_CONFIG.appEnv` parameter:

- **Development** (`appEnv: "development"`): Logs to `stdout` only
- **Production** (`appEnv: "production"`): Logs to `logs/app.log` with automatic rotation

### Log Rotation (Production Only)

- Max file size: 100MB
- Max backups: 5 files
- Max age: 30 days
- Compression: Enabled for rotated files

## HTTP Request Logging

All HTTP requests are automatically logged by the `LoggingMiddleware`:

```json
{
  "time": "2025-10-20T05:19:47.567282326-05:00",
  "level": "WARN",
  "msg": "HTTP",
  "request_id": "d7d17806-40c1-4dd0-be64-4baf644c4ac3",
  "method": "POST",
  "path": "/api/v1/parameters/get-by-code",
  "status": 422,
  "duration_ms": 0
}
```

### Log Levels for HTTP Requests

- `INFO`: 2xx responses (success)
- `WARN`: 4xx responses (client errors)
- `ERROR`: 5xx responses (server errors)

### Request ID Tracking

Every request automatically gets a unique `request_id`:
- Client can provide via `X-Request-ID` header
- Auto-generated if not provided
- Available in response via `X-Request-ID` header
- Propagated through all layers (middleware → router → use case)

## Logging in Use Cases

Import the logger helper:

```go
import "api-chatbot/internal/logger"
```

### Log an Error

```go
func (u *parameterUseCase) GetByCode(ctx context.Context, code string) d.Result[*d.Parameter] {
    param, err := u.repo.GetByCode(ctx, code)
    if err != nil {
        logger.LogError(ctx, "Failed to fetch parameter by code", err,
            "operation", "GetByCode",
            "code", code,
        )
        return d.Error[*d.Parameter](u.cache, "ERR_INTERNAL_DB")
    }
    // ...
}
```

Output:
```json
{
  "time": "2025-10-20T05:19:47.567282326-05:00",
  "level": "ERROR",
  "msg": "Failed to fetch parameter by code",
  "request_id": "d7d17806-40c1-4dd0-be64-4baf644c4ac3",
  "error": "connection timeout",
  "operation": "GetByCode",
  "code": "APP_CONFIG"
}
```

### Log a Warning

```go
if param == nil {
    logger.LogWarn(ctx, "Parameter not found",
        "operation", "GetByCode",
        "code", code,
    )
    return d.Error[*d.Parameter](u.cache, "ERR_PARAM_NOT_FOUND")
}
```

Output:
```json
{
  "time": "2025-10-20T05:19:47.567282326-05:00",
  "level": "WARN",
  "msg": "Parameter not found",
  "request_id": "d7d17806-40c1-4dd0-be64-4baf644c4ac3",
  "operation": "GetByCode",
  "code": "NONEXISTENT"
}
```

### Log Info

```go
logger.LogInfo(ctx, "Cache reloaded successfully",
    "operation", "ReloadCache",
    "param_count", len(params),
)
```

### Log Debug

```go
logger.LogDebug(ctx, "Cache hit",
    "operation", "GetByCode",
    "code", code,
)
```

## Available Logger Functions

All functions in `internal/logger/logger.go`:

```go
// Log error with automatic request_id tracking
logger.LogError(ctx context.Context, msg string, err error, args ...any)

// Log info with automatic request_id tracking
logger.LogInfo(ctx context.Context, msg string, args ...any)

// Log warning with automatic request_id tracking
logger.LogWarn(ctx context.Context, msg string, args ...any)

// Log debug with automatic request_id tracking
logger.LogDebug(ctx context.Context, msg string, args ...any)
```

## Best Practices

### 1. Always Pass Context

```go
// ✅ Good: Pass context to get request_id
logger.LogError(ctx, "Database error", err, "table", "users")

// ❌ Bad: Using slog directly without context
slog.Error("Database error", "error", err)
```

### 2. Add Relevant Context

Include operation-specific details:

```go
logger.LogError(ctx, "Failed to process payment", err,
    "operation", "ProcessPayment",
    "user_id", userID,
    "amount", amount,
    "currency", currency,
)
```

### 3. Use Appropriate Log Levels

- `ERROR`: Actual errors that need investigation
- `WARN`: Unexpected but recoverable situations (e.g., not found, validation)
- `INFO`: Important business events (e.g., cache reload, user login)
- `DEBUG`: Detailed debugging information (disabled in production)

### 4. Don't Log Sensitive Data

```go
// ❌ Bad: Logging sensitive information
logger.LogInfo(ctx, "User login", "password", password)

// ✅ Good: Log only necessary information
logger.LogInfo(ctx, "User login", "user_id", userID)
```

## Configuration

Configuration is stored in the `LOG_CONFIG` parameter in the database:

```json
{
  "level": "info",
  "format": "json",
  "output": "both",
  "filePath": "logs/app.log",
  "maxSizeMB": 100,
  "maxBackups": 5,
  "maxAgeDays": 30
}
```

**Note**: The `output` field is overridden by `APP_CONFIG.appEnv`:
- Development: Forces `stdout`
- Production: Forces `file`

To change log level, update the database:

```sql
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{level}', '"debug"')
WHERE prm_code = 'LOG_CONFIG';
```

Then restart the application.

## Example: Complete Use Case with Logging

```go
package usecase

import (
    "context"
    "time"

    d "api-chatbot/domain"
    "api-chatbot/internal/logger"
)

type documentUseCase struct {
    repo           d.DocumentRepository
    contextTimeout time.Duration
}

func (u *documentUseCase) GetByID(ctx context.Context, id int) d.Result[*d.Document] {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    logger.LogDebug(ctx, "Fetching document", "document_id", id)

    doc, err := u.repo.GetByID(ctx, id)
    if err != nil {
        logger.LogError(ctx, "Failed to fetch document", err,
            "operation", "GetByID",
            "document_id", id,
        )
        return d.Error[*d.Document](u.cache, "ERR_INTERNAL_DB")
    }

    if doc == nil {
        logger.LogWarn(ctx, "Document not found",
            "operation", "GetByID",
            "document_id", id,
        )
        return d.Error[*d.Document](u.cache, "ERR_DOCUMENT_NOT_FOUND")
    }

    logger.LogInfo(ctx, "Document fetched successfully",
        "operation", "GetByID",
        "document_id", id,
        "title", doc.Title,
    )

    return d.Success(doc)
}
```

## Viewing Logs

### Development (stdout)
Logs appear in your terminal:

```bash
./main
```

### Production (file)
View logs with:

```bash
# Tail live logs
tail -f logs/app.log

# Pretty print JSON logs
tail -f logs/app.log | jq .

# Filter by level
tail -f logs/app.log | jq 'select(.level == "ERROR")'

# Filter by request_id
tail -f logs/app.log | jq 'select(.request_id == "d7d17806-40c1-4dd0-be64-4baf644c4ac3")'
```
