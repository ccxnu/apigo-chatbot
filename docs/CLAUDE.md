# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based chatbot API built with chi router and PostgreSQL. The codebase follows a clean architecture pattern with clear separation between layers (controller, usecase, repository, domain).

## Development Commands

### Running the Application

```bash
# Run the application (requires config.json in root)
go run cmd/main.go

# Server runs on port 8080 by default
# http://localhost:8080
```

### Building

```bash
# Build the binary
go build -o main cmd/main.go

# Build with optimizations (used in Docker)
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o main cmd/main.go
```

### Docker Deployment

```bash
# Build and run container (production deployment)
./deploy.sh

# Container exposes port 3434 externally, maps to 8080 internally
# Requires config.json at /config/chatbot/config.json on host
```

### Dependencies

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

## Architecture

### Layer Structure

The codebase follows a clean architecture with these layers:

1. **Domain Layer** (`domain/`): Core business entities and interfaces
   - Defines interfaces for repositories, use cases, and caches
   - Contains domain models (Parameter, Document, Logger)
   - No dependencies on other layers

2. **Repository Layer** (`repository/`): Data access implementations
   - Implements domain interfaces for database operations
   - Uses the DAL (Data Access Layer) abstraction
   - All database calls go through PostgreSQL functions/procedures

3. **Use Case Layer** (`usecase/`): Business logic
   - Implements domain use case interfaces
   - Orchestrates repositories and caches
   - Handles timeout contexts and error mapping

4. **API Layer** (`api/`):
   - `route/`: Huma-based route definitions with built-in validation and OpenAPI documentation
   - `request/`: Request DTOs with validation tags (validated automatically by Huma)
   - `middleware/`: JWT auth, logger
   - `common/`: Common result types and error handling
   - `dal/`: Database abstraction layer

5. **Internal** (`internal/`):
   - `cache/`: In-memory caching implementations (thread-safe with sync.RWMutex)
   - `helper/`: Utility functions
   - `jwttoken/`: JWT token operations

6. **Config** (`config/`): Configuration management
   - Reads from `config.json` using Viper
   - Database connection pooling with pgxpool
   - Centralized environment configuration

### Database Access Pattern

The DAL provides generic functions that ALL database operations use:

- `QueryRows[T]()`: Execute PostgreSQL functions that return multiple rows
- `QueryRow[T]()`: Execute PostgreSQL functions that return single row
- `ExecProc[T]()`: Execute PostgreSQL stored procedures (CALL statement)

All database operations are defined as PostgreSQL functions or procedures in `db/` directory. The Go code never writes raw SQL queries - it only calls database functions/procedures by name.

Example:
```go
// Calls: SELECT * FROM fn_get_all_parameters()
params, err := dal.QueryRows[domain.Parameter](r.dal, ctx, "fn_get_all_parameters")

// Calls: CALL sp_create_parameter($1, $2, $3, $4)
result, err := dal.ExecProc[domain.AddParameterResult](r.dal, ctx, "sp_create_parameter", name, code, data, desc)
```

### Request Validation

All routes use Huma's built-in validation:

```go
// In route files - Huma handles validation automatically
huma.Register(humaAPI, huma.Operation{
	OperationID: "get-parameter-by-code",
	Method:      "POST",
	Path:        "/api/v1/parameters/get-by-code",
	Summary:     "Get parameter by code",
	Tags:        []string{"Parameters"},
}, func(ctx context.Context, input *struct {
	Body request.GetParameterByCodeRequest // Validated automatically
}) (*GetParameterByCodeResponse, error) {
	// input.Body is already validated - use it directly
	result := paramUseCase.GetByCode(ctx, input.Body.Code)
	return &GetParameterByCodeResponse{Body: result}, nil
})
```

Huma validation:
1. Automatically decodes JSON body
2. Validates using struct tags (`validate:"required"`, `validate:"min=3,max=100"`, etc.)
3. Returns 400 errors with detailed validation messages if validation fails
4. Passes validated data to your handler function

### Caching Strategy

The parameter system uses a two-tier caching approach:
- In-memory cache (thread-safe with sync.RWMutex)
- Cache-aside pattern: check cache first, load from DB on miss
- Cache updates happen on Add/Update/Delete operations
- Supports manual cache reload via `/reload-cache` endpoint

### Configuration

The application requires a `config.json` file in the root directory with these sections:
- `App`: Application settings (name, environment, timeout)
- `Log`: Logging configuration
- `Database`: PostgreSQL connection details
- `Jwt`: JWT token secrets and expiry
- `Email`: Email sender configuration
- `WppConnect`: WhatsApp connector base URL
- `Embedding`: OpenAI and Ollama embedding models
- `Llm`: LLM service configuration

### Database Setup

Database schema and functions are in `db/` directory:
- `00_database_setup.sql`: Creates database, schemas, and extensions (vector, pgcrypto, uuid-ossp)
- `01_create_tables.sql`: Table definitions
- `02_parameters_procedures.sql`: PostgreSQL functions and procedures for parameter operations
- `initial_data.sql`: Seed data

The database uses two schemas:
- `public`: Application tables and functions
- `ex`: PostgreSQL extensions (isolated to avoid naming conflicts)

### Route Organization

All routes use Huma for automatic OpenAPI 3.1 documentation and validation:
- Routes are split by feature in `api/route/` (e.g., `parameter_router.go`, `document_router.go`)
- Each router file registers Huma operations with the API
- OpenAPI docs available at `/docs`, spec at `/openapi`
- JWT middleware available but currently commented out

To add a new feature:
1. Create domain interface and model in `domain/`
2. Create repository implementation in `repository/`
3. Create use case in `usecase/`
4. Create request DTOs in `api/request/` with validation tags
5. Create response types and Huma operations in router file (`api/route/`)
6. Register router in `api/route/route.go`
