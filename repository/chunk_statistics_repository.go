package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetChunkStatistics  = "fn_get_chunk_statistics"
	fnGetTopChunksByUsage = "fn_get_top_chunks_by_usage"
	// Stored Procedures (Writes)
	spIncrementChunkUsage       = "sp_increment_chunk_usage"
	spUpdateChunkQualityMetrics = "sp_update_chunk_quality_metrics"
	spUpdateChunkStaleness      = "sp_update_chunk_staleness"
)

type chunkStatisticsRepository struct {
	dal *dal.DAL
}

func NewChunkStatisticsRepository(dal *dal.DAL) domain.ChunkStatisticsRepository {
	return &chunkStatisticsRepository{
		dal: dal,
	}
}

// GetByChunk retrieves statistics for a specific chunk
func (r *chunkStatisticsRepository) GetByChunk(ctx context.Context, chunkID int) (*domain.ChunkStatistics, error) {
	stats, err := dal.QueryRows[domain.ChunkStatistics](r.dal, ctx, fnGetChunkStatistics, chunkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunk statistics via %s: %w", fnGetChunkStatistics, err)
	}

	if len(stats) == 0 {
		return nil, nil
	}

	return &stats[0], nil
}

// GetTopByUsage retrieves most frequently used chunks
func (r *chunkStatisticsRepository) GetTopByUsage(ctx context.Context, limit int) ([]domain.TopChunkByUsage, error) {
	topChunks, err := dal.QueryRows[domain.TopChunkByUsage](r.dal, ctx, fnGetTopChunksByUsage, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top chunks by usage via %s: %w", fnGetTopChunksByUsage, err)
	}
	return topChunks, nil
}

// IncrementUsage increments usage count for a chunk
func (r *chunkStatisticsRepository) IncrementUsage(ctx context.Context, chunkID int) (*domain.IncrementChunkUsageResult, error) {
	result, err := dal.ExecProc[domain.IncrementChunkUsageResult](
		r.dal,
		ctx,
		spIncrementChunkUsage,
		chunkID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spIncrementChunkUsage, err)
	}

	return result, nil
}

// UpdateQualityMetrics updates RAG quality metrics for a chunk
func (r *chunkStatisticsRepository) UpdateQualityMetrics(ctx context.Context, params domain.UpdateChunkQualityMetricsParams) (*domain.UpdateChunkQualityMetricsResult, error) {
	result, err := dal.ExecProc[domain.UpdateChunkQualityMetricsResult](
		r.dal,
		ctx,
		spUpdateChunkQualityMetrics,
		params.ChunkID,
		params.PrecisionAtK,
		params.RecallAtK,
		params.F1AtK,
		params.MRR,
		params.MAP,
		params.NDCG,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateChunkQualityMetrics, err)
	}

	return result, nil
}

// UpdateStaleness updates staleness tracking for a chunk
func (r *chunkStatisticsRepository) UpdateStaleness(ctx context.Context, params domain.UpdateChunkStalenessParams) (*domain.UpdateChunkStalenessResult, error) {
	result, err := dal.ExecProc[domain.UpdateChunkStalenessResult](
		r.dal,
		ctx,
		spUpdateChunkStaleness,
		params.ChunkID,
		params.StalenessDays,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateChunkStaleness, err)
	}

	return result, nil
}
