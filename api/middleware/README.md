# Middleware

This directory contains global middleware for the API.

## Available Middlewares

### 1. CORSMiddleware (`cors.go`)
Handles Cross-Origin Resource Sharing (CORS) headers for all requests.

**Features:**
- Allows requests from any origin (`*`)
- Supports common HTTP methods: GET, POST, PUT, DELETE, OPTIONS
- Handles preflight requests (OPTIONS)
- Configures credential sharing and cache time

**Usage:**
Applied globally in `main.go`

---

### 2. AuthMiddleware (`auth.go`)
Validates custom Authorization header for API security.

**Features:**
- Validates Bearer token in Authorization header
- Skips validation for `/docs`, `/openapi` endpoints
- Uses environment variable `API_AUTH_TOKEN` for token validation
- Returns 401 Unauthorized if validation fails
- Development mode: If `API_AUTH_TOKEN` is not set, auth is skipped

**Configuration:**
Set the environment variable:
```bash
export API_AUTH_TOKEN="your-secret-token-here"
```

**Request Example:**
```bash
curl -H "Authorization: Bearer your-secret-token-here" \
  http://localhost:8080/api/v1/parameters/get-all
```

**Skip Auth (Development):**
Don't set `API_AUTH_TOKEN` environment variable, and auth will be bypassed.

---

### 3. Logger Middleware (`logger.go`)
Logs HTTP requests with details (method, path, duration, status).

---

### 4. JWT Auth Middleware (`jwt_auth_middleware.go`)
JWT token-based authentication (currently not in use, available for future implementation).

---

## Middleware Order

Middlewares are applied in this order (defined in `main.go`):
1. **CORS** - First, to handle preflight and set headers
2. **Auth** - Second, to validate authorization
3. **Handler** - Finally, the actual route handler

```go
handler := middleware.CORSMiddleware(
    middleware.AuthMiddleware(mux),
)
```

Order matters! CORS must be first to properly handle OPTIONS requests.
