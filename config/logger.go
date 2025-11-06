package config

import (
	"io"
	"log/slog"
	"os"

	"api-chatbot/domain"
)

type LogConfig struct {
	Level      string `json:"level"`      // debug, info, warn, error
	Format     string `json:"format"`     // json, text
	Output     string `json:"output"`     // stdout, file, both
	FilePath   string `json:"filePath"`   // logs/app.log
	MaxSizeMB  int    `json:"maxSizeMB"`  // max size in MB before rotation
	MaxBackups int    `json:"maxBackups"` // max number of old log files
	MaxAgeDays int    `json:"maxAgeDays"` // max days to retain old logs
}

func SetupLogger(cache domain.ParameterCache) (*slog.Logger, func() error) {
	config := LogConfig{
		Level:      "info",
		Format:     "json",
		Output:     "both",
		FilePath:   "logs/app.log",
		MaxSizeMB:  100,
		MaxBackups: 5,
		MaxAgeDays: 30,
	}

	// Check environment from APP_CONFIG to determine output mode
	isDevelopment := false
	if param, exists := cache.Get("APP_CONFIG"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if appEnv, ok := data["appEnv"].(string); ok {
				isDevelopment = (appEnv == "development")
			}
		}
	}

	// Try to load from cache
	if param, exists := cache.Get("LOG_CONFIG"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if level, ok := data["level"].(string); ok {
				config.Level = level
			}
			if format, ok := data["format"].(string); ok {
				config.Format = format
			}
			if output, ok := data["output"].(string); ok {
				config.Output = output
			}
			if filePath, ok := data["filePath"].(string); ok {
				config.FilePath = filePath
			}
			if maxSizeMB, ok := data["maxSizeMB"].(float64); ok {
				config.MaxSizeMB = int(maxSizeMB)
			}
			if maxBackups, ok := data["maxBackups"].(float64); ok {
				config.MaxBackups = int(maxBackups)
			}
			if maxAgeDays, ok := data["maxAgeDays"].(float64); ok {
				config.MaxAgeDays = int(maxAgeDays)
			}
		}
	}

	// Override output based on environment
	// Development: stdout only (easier for local debugging)
	// Production: file only (persistent logs, less noise in container logs)
	// if isDevelopment {
	config.Output = "stdout"
	// } else {
	// config.Output = "file"
	// }

	// Parse log level
	var level slog.Level
	switch config.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false, // Remove source for cleaner logs
	}

	var writers []io.Writer
	var fileWriter io.WriteCloser

	// Setup output destinations
	switch config.Output {
	case "stdout":
		writers = append(writers, os.Stdout)
	default:
		writers = append(writers, os.Stdout)
	}

	// Combine writers
	multiWriter := io.MultiWriter(writers...)

	// Create handler based on format
	var handler slog.Handler
	if config.Format == "text" {
		handler = slog.NewTextHandler(multiWriter, opts)
	} else {
		handler = slog.NewJSONHandler(multiWriter, opts)
	}

	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	// Log startup info
	logMsg := "Logger initialized"
	if isDevelopment {
		logMsg += " (development: stdout only)"
	} else {
		logMsg += " (production: file only)"
	}
	logger.Info(logMsg,
		"level", config.Level,
		"format", config.Format,
		"output", config.Output,
		"filePath", config.FilePath,
	)

	// Return cleanup function
	cleanup := func() error {
		if fileWriter != nil {
			logger.Info("Closing log file")
			return fileWriter.Close()
		}
		return nil
	}

	return logger, cleanup
}
