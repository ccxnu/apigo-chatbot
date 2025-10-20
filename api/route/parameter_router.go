package route

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

// Huma response types for parameters
type GetAllParametersResponse struct {
	Body d.Result[[]d.Parameter]
}

type GetParameterByCodeResponse struct {
	Body d.Result[*d.Parameter]
}

type AddParameterResponse struct {
	Body d.Result[d.Data]
}

type UpdParameterResponse struct {
	Body d.Result[d.Data]
}

type DelParameterResponse struct {
	Body d.Result[d.Data]
}

type ReloadCacheResponse struct {
	Body d.Result[d.Data]
}

func NewParameterRouter(paramUseCase d.ParameterUseCase, mux *http.ServeMux, humaAPI huma.API) {
	// Huma documented routes with /api/v1/ prefix
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-all-parameters",
		Method:      "POST",
		Path:        "/api/v1/parameters/get-all",
		Summary:     "Get all parameters",
		Description: "Retrieves all active system parameters from cache or database",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.GetAllParametersRequest
	}) (*GetAllParametersResponse, error) {
		result := paramUseCase.GetAll(ctx)
		return &GetAllParametersResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-parameter-by-code",
		Method:      "POST",
		Path:        "/api/v1/parameters/get-by-code",
		Summary:     "Get parameter by code",
		Description: "Retrieves a specific parameter by its unique code",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.GetParameterByCodeRequest
	}) (*GetParameterByCodeResponse, error) {
		result := paramUseCase.GetByCode(ctx, input.Body.Code)
		return &GetParameterByCodeResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "add-parameter",
		Method:      "POST",
		Path:        "/api/v1/parameters/add",
		Summary:     "Add parameter",
		Description: "Creates a new system parameter with validation",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.AddParameterRequest
	}) (*AddParameterResponse, error) {
		params := d.AddParameterParams{
			Name:        input.Body.Name,
			Code:        input.Body.Code,
			Data:        input.Body.Data,
			Description: input.Body.Description,
		}
		result := paramUseCase.AddParameter(ctx, params)
		return &AddParameterResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-parameter",
		Method:      "POST",
		Path:        "/api/v1/parameters/update",
		Summary:     "Update parameter",
		Description: "Updates an existing parameter",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdParameterRequest
	}) (*UpdParameterResponse, error) {
		params := d.UpdParameterParams{
			Code:        input.Body.Code,
			Name:        input.Body.Name,
			Data:        input.Body.Data,
			Description: input.Body.Description,
		}
		result := paramUseCase.UpdParameter(ctx, params)
		return &UpdParameterResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "delete-parameter",
		Method:      "POST",
		Path:        "/api/v1/parameters/delete",
		Summary:     "Delete parameter",
		Description: "Soft deletes a parameter (sets active = false)",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.DelParameterRequest
	}) (*DelParameterResponse, error) {
		result := paramUseCase.DelParameter(ctx, input.Body.Code)
		return &DelParameterResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "reload-parameter-cache",
		Method:      "POST",
		Path:        "/api/v1/parameters/reload-cache",
		Summary:     "Reload parameter cache",
		Description: "Reloads parameter cache from database",
		Tags:        []string{"Parameters"},
	}, func(ctx context.Context, input *struct {
		Body request.ReloadCacheRequest
	}) (*ReloadCacheResponse, error) {
		result := paramUseCase.ReloadCache(ctx)
		return &ReloadCacheResponse{Body: result}, nil
	})
}
