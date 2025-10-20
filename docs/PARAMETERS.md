# Configuration Parameters

This application stores most of its configuration in the database as parameters, allowing for dynamic updates without restarting the server.

## Configuration Locations

### Database Parameters (Updateable at runtime)
All application configuration except database connection details are stored in `cht_parameters` table.

### File Configuration (Requires restart)
Only database connection settings remain in `config.json`:
- Database host, port, user, password, name
- Database connection pool size

## Available Parameters

### APP_CONFIG
General application settings:
```json
{
  "application": "wsChatbot",
  "appEnv": "development",
  "contextTimeout": 2,
  "basicAuth": "X-Auth wiaAchcHks3rBxIhJQem1nLoMDwdoQ=="
}
```

### LOG_CONFIG
Logging system configuration:
```json
{
  "path": "/logs",
  "level": "debug",
  "maxSize": "20m",
  "maxFiles": "7d",
  "pattern": "YYYY-MM-DD",
  "saveToFile": true
}
```

### JWT_CONFIG
JWT token configuration:
```json
{
  "accessTokenSecret": "Y2xhdmVzdXBlcnNlY3JldGFAMTc2NjQ1NDU2Nw==",
  "accessTokenExpiryHour": 2,
  "refreshTokenSecret": "ttahoeuthaorcuhntuhoatneuh",
  "refreshTokenExpiryHour": 24
}
```

### EMAIL_CONFIG
Email sender configuration:
```json
{
  "sender": "noreply@example.com"
}
```

### WPP_CONNECT_CONFIG
WhatsApp Connect service configuration:
```json
{
  "baseUrl": "http://0.0.0.0:21465"
}
```

### EMBEDDING_CONFIG
Embedding models configuration (OpenAI and Ollama):
```json
{
  "openaiUrl": "https://api.openai.com/v1/embeddings",
  "openaiApiKey": "sk-proj-xxx",
  "openaiModel": "text-embedding-3-small",
  "ollamaUrl": "https://localhost:11434/embeddings/:model",
  "ollamaModel": "nomic"
}
```

### LLM_CONFIG
Large Language Model API configuration:
```json
{
  "url": "https://api.groq.com/openai/v1",
  "model": "llama-3.1-8b-instant",
  "apiKey": "gsk_xxx"
}
```

### RAG_CONFIG
Retrieval Augmented Generation settings:
```json
{
  "topK": 5,
  "similarityThreshold": 0.7,
  "chunkSize": 1000,
  "chunkOverlap": 200
}
```

## How to Update Configuration

### Via API (Recommended - No restart needed)

```bash
# Update a parameter
curl -X POST http://localhost:8080/api/v1/parameters/update \
  -H "Content-Type: application/json" \
  -H "X-App-Authorization: your-auth-token" \
  -d '{
    "code": "APP_CONFIG",
    "name": "APP_CONFIGURATION",
    "data": {
      "application": "wsChatbot",
      "appEnv": "production",
      "contextTimeout": 5,
      "basicAuth": "new-auth-token"
    },
    "description": "General application configuration",
    "idSession": "admin-session",
    "idRequest": "550e8400-e29b-41d4-a716-446655440000",
    "process": "update-config",
    "idDevice": "admin-device",
    "deviceAdress": "192.168.1.1",
    "dateProcess": "2025-10-19T20:00:00Z"
  }'

# Reload cache (if needed for immediate effect)
curl -X POST http://localhost:8080/api/v1/parameters/reload-cache \
  -H "Content-Type: application/json" \
  -H "X-App-Authorization: your-auth-token" \
  -d '{
    "idSession": "admin-session",
    "idRequest": "550e8400-e29b-41d4-a716-446655440001",
    "process": "reload-cache",
    "idDevice": "admin-device",
    "deviceAdress": "192.168.1.1",
    "dateProcess": "2025-10-19T20:00:00Z"
  }'
```

### Via Database (Direct SQL)

```sql
-- Update a parameter
UPDATE cht_parameters
SET prm_data = '{
  "application": "wsChatbot",
  "appEnv": "production",
  "contextTimeout": 5,
  "basicAuth": "new-auth-token"
}'::jsonb
WHERE prm_code = 'APP_CONFIG';

-- View all configuration parameters
SELECT prm_code, prm_name, prm_data, prm_description
FROM cht_parameters
WHERE prm_name LIKE '%CONFIGURATION';
```

## Error Codes

All error codes and messages are also stored as parameters:
- `ERR_BAD_REQUEST` - "Solicitud incorrecta. Verifique los datos enviados"
- `ERR_UNAUTHORIZED` - "No autorizado. Se requiere autenticación"
- `ERR_FORBIDDEN` - "Acceso prohibido. No tiene permisos suficientes"
- `ERR_NOT_FOUND` - "Recurso no encontrado"
- `ERR_VALIDATION_FAILED` - "Error de validación. Revise los campos enviados"
- `ERR_INTERNAL_SERVER` - "Error interno del servidor"
- And many more...

## Benefits

1. **Dynamic Updates**: Change configuration without restarting the server
2. **API Management**: Update via REST API endpoints
3. **Centralized**: All configuration in one place (database)
4. **Versioned**: Database migrations track configuration changes
5. **Cached**: Fast access via in-memory parameter cache
6. **Secure**: Sensitive data encrypted in database
