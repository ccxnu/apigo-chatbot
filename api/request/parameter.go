package request

import (
	"encoding/json"

	"api-chatbot/domain"
)

type GetAllParametersRequest struct {
	domain.Base
}

type GetParameterByCodeRequest struct {
	domain.Base
	Code string `json:"code" validate:"required"`
}

type AddParameterRequest struct {
	domain.Base
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Code        string          `json:"code" validate:"required,min=2,max=100"`
	Data        json.RawMessage `json:"data" validate:"required"`
	Description string          `json:"description" validate:"max=500"`
}

type UpdParameterRequest struct {
	domain.Base
	Code        string          `json:"code" validate:"required"`
	Name        string          `json:"name" validate:"required,min=3,max=100"`
	Data        json.RawMessage `json:"data" validate:"required"`
	Description string          `json:"description" validate:"max=500"`
}

type DelParameterRequest struct {
	domain.Base
	Code string `json:"code" validate:"required"`
}

type ReloadCacheRequest struct {
	domain.Base
}
