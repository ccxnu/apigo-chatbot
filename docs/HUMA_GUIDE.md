# Huma OpenAPI 3.1 Integration Guide

This project uses [Huma v2](https://github.com/danielgtaylor/huma) for automatic OpenAPI 3.1 documentation generation with a code-first approach.

## What's Included

✅ **Huma v2** - Modern OpenAPI 3.1 framework for Go
✅ **Chi Adapter** - Seamlessly works with existing chi router
✅ **Auto Documentation** - OpenAPI spec generated from Go code
✅ **Interactive UI** - Swagger UI, Redoc, and RapiDoc built-in
✅ **Type Safety** - Uses Go generics for request/response validation

## Quick Start

### 1. Access the Documentation

Start the server:
```bash
go run cmd/main.go
```

Visit the OpenAPI documentation:
- **Docs UI**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/openapi.json
- **OpenAPI YAML**: http://localhost:8080/openapi.yaml

### 2. Test the Example Endpoint

Health check endpoint (Huma-powered):
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2025-01-19T10:30:00Z",
  "version": "1.0.0"
}
```

## Architecture

### Current Setup

Your API has **two types of endpoints**:

1. **Existing Endpoints** (Original chi handlers)
   - Located in `api/controller/`
   - Use custom middleware for validation
   - Paths: `/parameters/*`, `/documents/*`, `/chunks/*`, etc.
   - **Continue working as before** ✅

2. **Huma Endpoints** (New OpenAPI documented)
   - Located in `api/huma/`
   - Auto-documented with OpenAPI 3.1
   - Auto-validated from Go types
   - Example path: `/health`, `/huma/parameters/*`

### File Structure

```
api/
├── huma/                    # NEW: Huma-powered endpoints
│   ├── health.go           # Health check example
│   └── parameter.go        # Parameter endpoints example
├── controller/             # EXISTING: Original controllers
├── route/                  # EXISTING: Original routes
└── ...
```

## Response Format Standard

**ALL endpoints MUST use `common.Result[T]` type:**

```go
type Result[T any] struct {
    Success bool   `json:"success"` // true for success, false for errors
    Code    string `json:"code"`    // "OK" for success, "ERR_..." for errors
    Info    string `json:"info"`    // Human-readable message from parameter cache
    Data    T      `json:"data"`    // Actual response data
}
```

**Success Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Operation successful",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "code": "ERR_PARAM_NOT_FOUND",
  "info": "Parameter not found in the system",
  "data": null
}
```

## How to Add Huma Endpoints

### Example 1: Simple GET Endpoint

```go
package huma

import (
    "context"
    "api-chatbot/api/common"
    "github.com/danielgtaylor/huma/v2"
)

type HelloData struct {
    Message string `json:"message" example:"Hello World" doc:"Greeting message"`
}

type HelloResponse struct {
    Body common.Result[HelloData]
}

func RegisterHelloRoute(api huma.API) {
    huma.Register(api, huma.Operation{
        OperationID: "say-hello",
        Method:      "GET",
        Path:        "/hello",
        Summary:     "Say hello",
        Tags:        []string{"Greetings"},
    }, func(ctx context.Context, input *struct{}) (*HelloResponse, error) {
        resp := &HelloResponse{}
        resp.Body = common.Result[HelloData]{
            Success: true,
            Code:    "OK",
            Info:    "Greeting generated successfully",
            Data: HelloData{
                Message: "Hello World",
            },
        }
        return resp, nil
    })
}
```

### Example 2: POST with Request Body

```go
type CreateUserInput struct {
    Body struct {
        Name  string `json:"name" minLength:"1" maxLength:"100" doc:"User name"`
        Email string `json:"email" format:"email" doc:"User email"`
    }
}

type CreateUserData struct {
    ID int `json:"id" doc:"Created user ID"`
}

type CreateUserResponse struct {
    Body common.Result[CreateUserData]
}

func RegisterUserRoute(api huma.API, paramCache domain.ParameterCache) {
    huma.Register(api, huma.Operation{
        OperationID: "create-user",
        Method:      "POST",
        Path:        "/users",
        Summary:     "Create new user",
        Tags:        []string{"Users"},
    }, func(ctx context.Context, input *CreateUserInput) (*CreateUserResponse, error) {
        // Your business logic here
        // If error occurs, return error response
        if input.Body.Name == "" {
            return &CreateUserResponse{
                Body: common.Result[CreateUserData]{
                    Success: false,
                    Code:    "ERR_INVALID_NAME",
                    Info:    "Name cannot be empty",
                    Data:    CreateUserData{},
                },
            }, nil
        }

        // On success
        resp := &CreateUserResponse{}
        resp.Body = common.Result[CreateUserData]{
            Success: true,
            Code:    "OK",
            Info:    "User created successfully",
            Data: CreateUserData{
                ID: 123,
            },
        }
        return resp, nil
    })
}
```

### Example 3: Using Existing Use Cases

See `api/huma/parameter.go` for a complete example of wrapping existing use cases.

**Key Points:**
- Use cases already return `common.Result[T]`
- Simply assign the result to `resp.Body`
- No need to manually construct success/error responses
- Error messages come from parameter cache automatically

```go
type GetParameterByCodeInput struct {
    Body struct {
        Code string `json:"code" minLength:"1" doc:"Parameter code"`
    }
}

type GetParameterByCodeResponse struct {
    Body common.Result[*domain.Parameter]
}

func RegisterParameterRoutes(api huma.API, paramUseCase domain.ParameterUseCase) {
    huma.Register(api, huma.Operation{
        OperationID: "get-parameter-by-code",
        Method:      "POST",
        Path:        "/huma/parameters/get-by-code",
        Summary:     "Get parameter by code",
        Description: "Returns Success: true with Code: OK if found, or Success: false with ERR_PARAM_NOT_FOUND",
        Tags:        []string{"Parameters"},
    }, func(ctx context.Context, input *GetParameterByCodeInput) (*GetParameterByCodeResponse, error) {
        // Use case already returns common.Result[*domain.Parameter]
        result := paramUseCase.GetByCode(ctx, input.Body.Code)

        // Simply assign to response
        resp := &GetParameterByCodeResponse{}
        resp.Body = result // Already has Success, Code, Info, Data

        return resp, nil
    })
}
```

## Validation

Huma automatically validates requests based on struct tags:

### Available Validation Tags

```go
type ExampleInput struct {
    Body struct {
        // String validation
        Name string `json:"name" minLength:"3" maxLength:"100"`

        // Number validation
        Age int `json:"age" minimum:"0" maximum:"150"`

        // Format validation
        Email string `json:"email" format:"email"`
        URL   string `json:"url" format:"uri"`

        // Pattern validation
        Code string `json:"code" pattern:"^[A-Z0-9_]+$"`

        // Required/Optional
        Required string  `json:"required"`
        Optional *string `json:"optional,omitempty"`

        // Enums
        Status string `json:"status" enum:"active,inactive,pending"`

        // Arrays
        Tags []string `json:"tags" minItems:"1" maxItems:"10"`
    }
}
```

## Migration Strategy

You have **3 options** for migrating existing endpoints:

### Option 1: Hybrid Approach (Recommended)
Keep existing endpoints as-is, add new Huma endpoints gradually:
- ✅ No breaking changes
- ✅ Existing endpoints keep working
- ✅ New features get OpenAPI docs
- ✅ Migrate when convenient

### Option 2: Parallel Endpoints
Create Huma versions alongside existing endpoints:
- Existing: `/parameters/get-all`
- Huma: `/huma/parameters/get-all`
- Both work simultaneously
- Deprecate old ones later

### Option 3: Full Migration
Convert all endpoints to Huma:
- Most work upfront
- Full OpenAPI coverage
- Consistent API design

## Benefits of Huma

### 1. Automatic Documentation
Write code → Get OpenAPI spec automatically ✨

### 2. Request Validation
No need for manual validation middleware:
```go
// Before (manual validation)
if req.Name == "" || len(req.Name) < 3 {
    return errors.New("name too short")
}

// After (automatic with Huma)
type Input struct {
    Body struct {
        Name string `json:"name" minLength:"3"` // Auto-validated!
    }
}
```

### 3. Type Safety
Compile-time safety for requests and responses.

### 4. Multiple UI Options
Built-in support for:
- Swagger UI (default at `/docs`)
- Redoc
- RapiDoc
- Stoplight Elements

### 5. Auto Content Negotiation
Supports JSON, YAML, CBOR, MessagePack automatically.

## Configuration

Edit `cmd/main.go` to customize:

```go
humaConfig := humaLib.DefaultConfig("Your API Name", "2.0.0")
humaConfig.Info.Description = "Your description"
humaConfig.Info.Contact = &humaLib.Contact{
    Name:  "Your Team",
    Email: "api@example.com",
    URL:   "https://example.com",
}
humaConfig.Servers = []*humaLib.Server{
    {URL: "https://api.example.com", Description: "Production"},
    {URL: "http://localhost:8080", Description: "Development"},
}
```

## Resources

- **Huma Documentation**: https://huma.rocks
- **GitHub**: https://github.com/danielgtaylor/huma
- **Examples**: https://github.com/danielgtaylor/huma/tree/main/examples

## Next Steps

1. ✅ Access http://localhost:8080/docs
2. ✅ Test the `/health` endpoint
3. ✅ Review example in `api/huma/parameter.go`
4. Create your first Huma endpoint
5. Gradually migrate existing endpoints (optional)

Happy coding! 🚀
