# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based WhatsApp chatbot API built with Huma router and PostgreSQL. The codebase follows a clean architecture pattern with clear separation between layers (controller, usecase, repository, domain). The system uses RAG (Retrieval Augmented Generation) with vector embeddings for intelligent responses.

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
   - `whatsapp/`: WhatsApp client and message handlers
   - `llm/`: LLM integration (Groq, OpenAI compatible)
   - `embedding/`: Vector embedding services (OpenAI, Ollama)
   - `httpclient/`: HTTP client for external APIs
   - `migration/`: Database migration system (golang-migrate)
   - `reports/`: Report generation (Excel, PDF)

6. **Config** (`config/`): Configuration management
   - Reads from `config.json` using Viper
   - Database connection pooling with pgxpool
   - Centralized environment configuration
   - Migration configuration (AUTO_MIGRATE, VERBOSE)

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
- `Database`: PostgreSQL connection details (host, port, user, password, name)
- `Migration`: Migration settings (AUTO_MIGRATE, VERBOSE)
- **Parameters**: All other configuration is stored in the database `cht_parameters` table and cached at startup

Important: Don't add new sections to `config.json`. Use the parameter system instead:
```sql
INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
VALUES ('Config Name', 'CONFIG_CODE', '{"key": "value"}', 'Description');
```

### Database Migration System

**IMPORTANT**: The project uses automatic database migrations with golang-migrate.

**DO NOT manually execute SQL files from `db/` directory - they are legacy reference only.**

Active migrations are in: `internal/migration/migrations/`

Migration files follow the pattern: `{version}_{description}.{up|down}.sql`
- Example: `000001_database_setup.up.sql`, `000001_database_setup.down.sql`

#### Running Migrations

**Automatic (Development)**:
```bash
# Set in config.json: "Migration": {"AUTO_MIGRATE": true, "VERBOSE": true}
go build -o main cmd/main.go
./main  # Migrations run automatically on startup
```

**Manual Control**:
```bash
# Build migration CLI
go build -o migrate cmd/migrate/main.go

# Check current version
./migrate -version

# Apply pending migrations
./migrate -up

# Rollback last migration
./migrate -down

# Fix dirty state (use with caution)
./migrate -force 14
```

#### Creating New Migrations

1. Find next version number: `./migrate -version` (e.g., shows 14)
2. Create migration files:
```bash
touch internal/migration/migrations/000015_add_feature.up.sql
touch internal/migration/migrations/000015_add_feature.down.sql
```
3. Write SQL (always use `IF NOT EXISTS` / `IF EXISTS`)
4. Rebuild app (migrations are embedded): `go build -o main cmd/main.go`
5. Test: `./migrate -up` then `./migrate -down` then `./migrate -up`

#### Migration Version Tracking

Migrations are tracked in `schema_migrations` table:
```sql
SELECT * FROM schema_migrations;
-- version | dirty
--      14 | f     (f = clean, t = needs manual fix)
```

#### Database Schemas

- `public`: Application tables and functions
- `ex`: PostgreSQL extensions (vector, pgcrypto, uuid-ossp) - isolated to avoid conflicts

### Route Organization

All routes use Huma for automatic OpenAPI 3.1 documentation and validation:
- Routes are split by feature in `api/route/` (e.g., `parameter_router.go`, `document_router.go`)
- Each router file registers Huma operations with the API
- OpenAPI docs available at `/docs`, spec at `/openapi`
- JWT middleware available but currently commented out

To add a new feature:
1. **Database First**: Create migration with tables and stored procedures
   ```bash
   touch internal/migration/migrations/000015_add_feature.up.sql
   touch internal/migration/migrations/000015_add_feature.down.sql
   ```
2. Create domain interface and model in `domain/`
3. Create repository implementation in `repository/` (calls PostgreSQL functions/procedures)
4. Create use case in `usecase/`
5. Create request DTOs in `api/request/` with validation tags
6. Create response types and Huma operations in router file (`api/route/`)
7. Register router in `api/route/route.go`
8. Rebuild app: `go build -o main cmd/main.go`

### SQL Naming Conventions

**IMPORTANT**: Follow these strict naming conventions for all SQL code:

#### Variables (PLpgSQL):
- Input parameters: `p_` prefix (e.g., `p_user_id`, `p_status`)
- Output parameters: `o_` prefix (e.g., `o_parameter_id`)
  - Standard outputs: `success BOOLEAN`, `code VARCHAR`
- Local variables: `v_` prefix (e.g., `v_mod_id`, `v_count`)
- Records: `r_` prefix (e.g., `r_user`)
- Counters: `i_` prefix (e.g., `i_counter`)
- Booleans: `bl_` or `is_` prefix (e.g., `bl_exists`, `is_active`)
- Constants: `c_` prefix (e.g., `c_default_status`)

#### Functions and Procedures:
- Stored procedures: `sp_` prefix (e.g., `sp_add_user`, `sp_update_status`)
  - Always return: `success BOOLEAN`, `code VARCHAR` at first values (e.g., 'OK', 'ERR_NOT_FOUND')
- Functions: `fn_` prefix (e.g., `fn_get_parameters`, `fn_search_chunks`)
- Views: `vw_` prefix (e.g., `vw_user_permissions`)
- Triggers: `tr_` prefix (e.g., `tr_update_timestamp`)

#### General Rules:
- Use `snake_case` for all identifiers
- SQL keywords in lowercase
- Always comment code blocks
- Scripts must be idempotent (use `IF NOT EXISTS`, `IF EXISTS`)
- 4-space or tab indentation

Example:
```sql
CREATE OR REPLACE PROCEDURE sp_create_user(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_username VARCHAR,
    IN p_email VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_user_id INT;
BEGIN
    o_success := TRUE;
    o_code := 'OK';

    -- Check if user exists
    IF EXISTS (SELECT 1 FROM cht_users WHERE usr_email = p_email) THEN
        o_success := FALSE;
        o_code := 'ERR_USER_EXISTS';
        RETURN;
    END IF;

    -- Insert user
    INSERT INTO cht_users (usr_username, usr_email)
    VALUES (p_username, p_email)
    RETURNING usr_id INTO v_user_id;

EXCEPTION
    WHEN OTHERS THEN
        o_success := FALSE;
        o_code := 'ERR_CREATE_USER';
END;
$$;
```

### Documentation

**Single Source of Truth**: `docs/manual_programador.typ`

This Typst document contains all programmer documentation:
- Architecture overview
- Database schema and migration guide
- SQL naming conventions and best practices
- API documentation
- Development workflows

To compile:
```bash
typst compile docs/manual_programador.typ
# Generates: docs/manual_programador.pdf
```

**Do not create new .md documentation files** - add content to the Typst manual instead.
