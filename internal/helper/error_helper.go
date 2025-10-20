package helper

import (
	"context"

	"api-chatbot/domain"
)

type ErrorHelper struct {
	paramUseCase domain.ParameterUseCase
}

func NewErrorHelper(paramUseCase domain.ParameterUseCase) *ErrorHelper {
	return &ErrorHelper{
		paramUseCase: paramUseCase,
	}
}

// GetErrorMessage retrieves the error message from parameters by error code
func (h *ErrorHelper) GetErrorMessage(ctx context.Context, errorCode string) string {
	// Try to get from cache first
	param, err := h.paramUseCase.GetByCode(ctx, errorCode)
	if err != nil || param == nil {
		// Fallback to generic message
		return "Ha ocurrido un error"
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		return "Ha ocurrido un error"
	}

	if message, ok := data["message"].(string); ok {
		return message
	}

	return "Ha ocurrido un error"
}
