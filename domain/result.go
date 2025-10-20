package domain

import (
	"encoding/json"
	"net/http"
)

type Data map[string]any

type Result[T any] struct {
	Success bool   `json:"success"` // Indicates if the operation was successful
	Code    string `json:"code"`    // Error/success code (e.g., "OK", "ERR_INTERNAL_DB")
	Info    string `json:"info"`    // Human-readable message
	Data    T      `json:"data"`    // Response payload
}

func Error[T any](cache ParameterCache, errorCode string) Result[T] {
	info := "Ha ocurrido un error"

	if param, exists := cache.Get(errorCode); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if msg, ok := data["message"].(string); ok {
				info = msg
			}
		}
	}
	return Result[T]{
		Success: false,
		Code:    errorCode,
		Info:    info,
		Data:    *new(T),
	}
}

func Success[T any](data T) Result[T] {
	return Result[T]{
		Success: true,
		Code:    "OK",
		Info:    "Operación exitosa",
		Data:    data,
	}
}

// HTTP Response Helper Functions (for middleware and error handling)
func writeResponse[T any](w http.ResponseWriter, statusCode int, res Result[T]) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"success":false,"code":"ERR_INTERNAL_SERVER","info":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// AppError sends an HTTP error response with Result format
func AppError(w http.ResponseWriter, statusCode int, errorCode string, info string) {
	res := Result[Data]{
		Success: false,
		Code:    errorCode,
		Info:    info,
		Data:    make(Data),
	}
	writeResponse(w, statusCode, res)
}
