# Using ParameterCache Directly

The application now uses `ParameterCache` directly instead of the old `Env` struct. This makes configuration more flexible and allows updates without restarting the server.

## How It Works

1. **On Startup**: All parameters are loaded from the database into `ParameterCache`
2. **During Runtime**: Access parameters directly from the cache
3. **On Update**: Call `/parameters/reload-cache` to refresh without restart

## Accessing Parameters

### In Your Code

```go
// Get a parameter value
if param, exists := paramCache.Get("APP_CONFIG"); exists {
    if data, err := param.GetDataAsMap(); err == nil {
        // Access fields
        appEnv := data["appEnv"].(string)
        timeout := int(data["contextTimeout"].(float64))
        basicAuth := data["basicAuth"].(string)
    }
}
```

### Available Parameter Codes

| Code | Description | Example Fields |
|------|-------------|----------------|
| `APP_CONFIG` | Application settings | `application`, `appEnv`, `contextTimeout`, `basicAuth` |
| `LOG_CONFIG` | Logging configuration | `path`, `level`, `maxSize`, `maxFiles`, `pattern`, `saveToFile` |
| `JWT_CONFIG` | JWT tokens | `accessTokenSecret`, `accessTokenExpiryHour`, `refreshTokenSecret`, `refreshTokenExpiryHour` |
| `EMAIL_CONFIG` | Email settings | `sender` |
| `WPP_CONNECT_CONFIG` | WhatsApp Connect | `baseUrl` |
| `EMBEDDING_CONFIG` | Embedding models | `openaiUrl`, `openaiApiKey`, `openaiModel`, `ollamaUrl`, `ollamaModel` |
| `LLM_CONFIG` | LLM API | `url`, `model`, `apiKey` |
| `RAG_CONFIG` | RAG settings | `topK`, `similarityThreshold`, `chunkSize`, `chunkOverlap` |

### Helper Function (Recommended)

Create a helper function to safely get values:

```go
func getConfigString(cache domain.ParameterCache, code, key, defaultValue string) string {
    if param, exists := cache.Get(code); exists {
        if data, err := param.GetDataAsMap(); err == nil {
            if val, ok := data[key].(string); ok {
                return val
            }
        }
    }
    return defaultValue
}

// Usage
basicAuth := getConfigString(paramCache, "APP_CONFIG", "basicAuth", "")
```

## Examples

### Example 1: Auth Middleware (Already Implemented)

```go
// api/middleware/auth.go
func AuthMiddleware(next http.Handler, paramCache domain.ParameterCache) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get basicAuth from parameter cache
        var basicAuth string
        if param, exists := paramCache.Get("APP_CONFIG"); exists {
            if data, err := param.GetDataAsMap(); err == nil {
                if auth, ok := data["basicAuth"].(string); ok {
                    basicAuth = auth
                }
            }
        }

        // Use basicAuth for validation
        // ...
    })
}
```

### Example 2: Email Service

```go
type EmailService struct {
    paramCache domain.ParameterCache
}

func (s *EmailService) SendEmail(to, subject, body string) error {
    var sender string
    if param, exists := s.paramCache.Get("EMAIL_CONFIG"); exists {
        if data, err := param.GetDataAsMap(); err == nil {
            sender = data["sender"].(string)
        }
    }

    // Use sender to send email
    // ...
}
```

### Example 3: LLM Client

```go
type LLMClient struct {
    paramCache domain.ParameterCache
}

func (c *LLMClient) MakeRequest(prompt string) (string, error) {
    // Get LLM config from cache
    if param, exists := c.paramCache.Get("LLM_CONFIG"); exists {
        if data, err := param.GetDataAsMap(); err == nil {
            url := data["url"].(string)
            model := data["model"].(string)
            apiKey := data["apiKey"].(string)

            // Make API request with these values
            // ...
        }
    }

    return "", nil
}
```

## Benefits

1. **No Struct Mapping**: Access values directly without intermediate structs
2. **Runtime Updates**: Change config via API without restart
3. **Flexible**: Easy to add new parameters
4. **Cached**: Fast in-memory access
5. **Simple**: Less code, more direct

## Migration from Old Env

**Before** (old way with Env struct):
```go
func SomeFunction(env *config.Env) {
    timeout := env.App.ContextTimeOut
    basicAuth := env.App.BasicAuth
}
```

**After** (new way with ParameterCache):
```go
func SomeFunction(paramCache domain.ParameterCache) {
    var timeout int
    var basicAuth string

    if param, exists := paramCache.Get("APP_CONFIG"); exists {
        if data, err := param.GetDataAsMap(); err == nil {
            timeout = int(data["contextTimeout"].(float64))
            basicAuth = data["basicAuth"].(string)
        }
    }
}
```

## Updating Configuration

```bash
# Update a parameter
curl -X POST http://localhost:8080/api/v1/parameters/update \
  -H "Content-Type: application/json" \
  -H "X-App-Authorization: your-token" \
  -d '{
    "code": "APP_CONFIG",
    "name": "APP_CONFIGURATION",
    "data": {
      "application": "wsChatbot",
      "appEnv": "production",
      "contextTimeout": 5,
      "basicAuth": "new-token"
    },
    ...
  }'

# Reload cache (instant effect)
curl -X POST http://localhost:8080/api/v1/parameters/reload-cache \
  -H "Content-Type: application/json" \
  -H "X-App-Authorization: your-token" \
  -d '{...}'
```
