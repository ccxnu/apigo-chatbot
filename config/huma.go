package config

import (
	"fmt"
	"net/http"

	"api-chatbot/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// https://github.com/danielgtaylor/huma/blob/main/defaults.go#L53
func CustomHumaConfig(title, version string) huma.Config {
	schemaPrefix := "#/components/schemas/"
	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	config := huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   title,
				Version: version,
			},
			Components: &huma.Components{
				Schemas: registry,
			},
		},
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs",
		SchemasPath:   "",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks:   []func(huma.Config) huma.Config{},
	}

	// Add custom error transformer to convert Huma errors to our Result format
	config.Transformers = []huma.Transformer{customErrorTransformer}

	return config
}

func customErrorTransformer(ctx huma.Context, status string, v any) (any, error) {
	if errorModel, ok := v.(*huma.ErrorModel); ok {
		errorCode := "ERR_UNKNOWN"
		switch errorModel.Status {
		case http.StatusBadRequest:
			errorCode = "ERR_BAD_REQUEST"
		case http.StatusUnauthorized:
			errorCode = "ERR_UNAUTHORIZED"
		case http.StatusForbidden:
			errorCode = "ERR_FORBIDDEN"
		case http.StatusNotFound:
			errorCode = "ERR_NOT_FOUND"
		case http.StatusUnprocessableEntity:
			errorCode = "ERR_VALIDATION_FAILED"
		case http.StatusInternalServerError:
			errorCode = "ERR_INTERNAL_SERVER"
		case http.StatusServiceUnavailable:
			errorCode = "ERR_SERVICE_UNAVAILABLE"
		}

		// Build error details for the data field
		data := make(map[string]any)
		if len(errorModel.Errors) > 0 {
			// Put validation error details in data field
			errors := make([]map[string]string, 0, len(errorModel.Errors))
			for _, err := range errorModel.Errors {
				errors = append(errors, map[string]string{
					"location": err.Location,
					"message":  err.Message,
					"value":    fmt.Sprintf("%v", err.Value),
				})
			}
			data["errors"] = errors
		}

		// Info comes from database parameter (error code will be looked up)
		// For now, use the detail as fallback
		info := errorModel.Detail

		return domain.Result[map[string]any]{
			Success: false,
			Code:    errorCode,
			Info:    info,
			Data:    data,
		}, nil
	}

	return v, nil
}

func NewHumaAPI(mux *http.ServeMux, paramCache domain.ParameterCache) huma.API {
	humaConfig := CustomHumaConfig("ISTS Chatbot API", "1.0.0")

	humaConfig.Info.Description = "RAG-based chatbot API for institute knowledge management"
	humaConfig.Info.Contact = &huma.Contact{
		Name: "ISTS Development Team",
	}

	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8080", Description: "Development server"},
	}

	// Add security scheme for X-App-Authorization header
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"AppAuth": {
			Type:        "apiKey",
			In:          "header",
			Name:        "X-App-Authorization",
			Description: "Custom authorization token",
		},
	}

	// Apply security globally (except for /docs, /openapi which are whitelisted in middleware)
	humaConfig.Security = []map[string][]string{
		{"AppAuth": {}},
	}

	return humago.New(mux, humaConfig)
}
