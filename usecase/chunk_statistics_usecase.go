package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type chunkStatisticsUseCase struct {
	statsRepo      d.ChunkStatisticsRepository
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewChunkStatisticsUseCase(
	statsRepo d.ChunkStatisticsRepository,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.ChunkStatisticsUseCase {
	return &chunkStatisticsUseCase{
		statsRepo:      statsRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

func (u *chunkStatisticsUseCase) GetByChunk(c context.Context, chunkID int) d.Result[*d.ChunkStatistics] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	stats, err := u.statsRepo.GetByChunk(ctx, chunkID)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch chunk statistics from database", err,
			"operation", "GetByChunk",
			"chunkID", chunkID,
		)
		return d.Error[*d.ChunkStatistics](u.paramCache, "ERR_INTERNAL_DB")
	}

	if stats == nil {
		logger.LogWarn(ctx, "Chunk statistics not found",
			"operation", "GetByChunk",
			"chunkID", chunkID,
		)
		return d.Error[*d.ChunkStatistics](u.paramCache, "ERR_CHUNK_STATS_NOT_FOUND")
	}

	return d.Success(stats)
}

func (u *chunkStatisticsUseCase) GetTopByUsage(c context.Context, limit int) d.Result[[]d.TopChunkByUsage] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	topChunks, err := u.statsRepo.GetTopByUsage(ctx, limit)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch top chunks by usage from database", err,
			"operation", "GetTopByUsage",
			"limit", limit,
		)
		return d.Error[[]d.TopChunkByUsage](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(topChunks)
}

func (u *chunkStatisticsUseCase) IncrementUsage(c context.Context, chunkID int) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.IncrementUsage(ctx, chunkID)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to increment chunk usage in database", err,
			"operation", "IncrementUsage",
			"chunkID", chunkID,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Chunk usage increment failed with business logic error",
			"operation", "IncrementUsage",
			"code", result.Code,
			"chunkID", chunkID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *chunkStatisticsUseCase) UpdateQualityMetrics(c context.Context, params d.UpdateChunkQualityMetricsParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.UpdateQualityMetrics(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to update chunk quality metrics in database", err,
			"operation", "UpdateQualityMetrics",
			"chunkID", params.ChunkID,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Chunk quality metrics update failed with business logic error",
			"operation", "UpdateQualityMetrics",
			"code", result.Code,
			"chunkID", params.ChunkID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *chunkStatisticsUseCase) UpdateStaleness(c context.Context, params d.UpdateChunkStalenessParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.statsRepo.UpdateStaleness(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to update chunk staleness in database", err,
			"operation", "UpdateStaleness",
			"chunkID", params.ChunkID,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Chunk staleness update failed with business logic error",
			"operation", "UpdateStaleness",
			"code", result.Code,
			"chunkID", params.ChunkID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}
