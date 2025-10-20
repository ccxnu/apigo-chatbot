package usecase

import (
	"context"
	"time"

	"api-chatbot/domain"
)

type chunkStatisticsUseCase struct {
	statsRepo      domain.ChunkStatisticsRepository
	paramCache     domain.ParameterCache
	contextTimeout time.Duration
}

func NewChunkStatisticsUseCase(
	statsRepo domain.ChunkStatisticsRepository,
	paramCache domain.ParameterCache,
	timeout time.Duration,
) domain.ChunkStatisticsUseCase {
	return &chunkStatisticsUseCase{
		statsRepo:      statsRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// getErrorMessage retrieves error message from parameter cache
func (u *chunkStatisticsUseCase) getErrorMessage(errorCode string) string {
	if param, exists := u.paramCache.Get(errorCode); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if message, ok := data["message"].(string); ok {
				return message
			}
		}
	}
	return "Ha ocurrido un error"
}

func (u *chunkStatisticsUseCase) GetByChunk(c context.Context, chunkID int) domain.Result[*domain.ChunkStatistics] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	stats, err := u.statsRepo.GetByChunk(ctx, chunkID)
	if err != nil {
		return domain.Result[*domain.ChunkStatistics]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if stats == nil {
		return domain.Result[*domain.ChunkStatistics]{
			Success: false,
			Code:    "ERR_CHUNK_STATS_NOT_FOUND",
			Info:    u.getErrorMessage("ERR_CHUNK_STATS_NOT_FOUND"),
			Data:    nil,
		}
	}

	return domain.Result[*domain.ChunkStatistics]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    stats,
	}
}

func (u *chunkStatisticsUseCase) GetTopByUsage(c context.Context, limit int) domain.Result[[]domain.TopChunkByUsage] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	topChunks, err := u.statsRepo.GetTopByUsage(ctx, limit)
	if err != nil {
		return domain.Result[[]domain.TopChunkByUsage]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    []domain.TopChunkByUsage{},
		}
	}

	return domain.Result[[]domain.TopChunkByUsage]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    topChunks,
	}
}

func (u *chunkStatisticsUseCase) IncrementUsage(c context.Context, chunkID int) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.IncrementUsage(ctx, chunkID)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}

func (u *chunkStatisticsUseCase) UpdateQualityMetrics(c context.Context, params domain.UpdateChunkQualityMetricsParams) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.UpdateQualityMetrics(ctx, params)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}

func (u *chunkStatisticsUseCase) UpdateStaleness(c context.Context, params domain.UpdateChunkStalenessParams) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.UpdateStaleness(ctx, params)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}
