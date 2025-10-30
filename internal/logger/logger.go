package logger

import (
	"context"
	"log/slog"

	"api-chatbot/internal/contextutil"
)

// LogError logs an error with context information including request ID
func LogError(ctx context.Context, msg string, err error, args ...any) {
	logArgs := []any{"error", err.Error()}

	// Add request ID if available
	if reqID := contextutil.GetRequestID(ctx); reqID != "" {
		logArgs = append(logArgs, "request_id", reqID)
	}

	// Add any additional arguments
	logArgs = append(logArgs, args...)

	slog.ErrorContext(ctx, msg, logArgs...)
}

// LogInfo logs an informational message with context
func LogInfo(ctx context.Context, msg string, args ...any) {
	logArgs := []any{}

	// Add request ID if available
	if reqID := contextutil.GetRequestID(ctx); reqID != "" {
		logArgs = append(logArgs, "request_id", reqID)
	}

	// Add any additional arguments
	logArgs = append(logArgs, args...)

	slog.InfoContext(ctx, msg, logArgs...)
}

// LogWarn logs a warning message with context
func LogWarn(ctx context.Context, msg string, args ...any) {
	logArgs := []any{}

	// Add request ID if available
	if reqID := contextutil.GetRequestID(ctx); reqID != "" {
		logArgs = append(logArgs, "request_id", reqID)
	}

	// Add any additional arguments
	logArgs = append(logArgs, args...)

	slog.WarnContext(ctx, msg, logArgs...)
}

// LogDebug logs a debug message with context
func LogDebug(ctx context.Context, msg string, args ...any) {
	logArgs := []any{}

	// Add request ID if available
	if reqID := contextutil.GetRequestID(ctx); reqID != "" {
		logArgs = append(logArgs, "request_id", reqID)
	}

	// Add any additional arguments
	logArgs = append(logArgs, args...)

	slog.DebugContext(ctx, msg, logArgs...)
}
