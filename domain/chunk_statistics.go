package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
)

type ChunkStatistics struct {
	ID                    int        `json:"id" db:"cst_id"`
	ChunkID               int        `json:"chunkId" db:"cst_fk_chunk"`
	UsageCount            int        `json:"usageCount" db:"cst_usage_count"`
	LastUsedAt            *time.Time `json:"lastUsedAt" db:"cst_last_used_at"`
	PrecisionAtK          *float64   `json:"precisionAtK" db:"cst_precision_atk"`
	RecallAtK             *float64   `json:"recallAtK" db:"cst_recall_atk"`
	F1AtK                 *float64   `json:"f1AtK" db:"cst_f1_atk"`
	MRR                   *float64   `json:"mrr" db:"cst_mrr"`
	MAP                   *float64   `json:"map" db:"cst_map"`
	NDCG                  *float64   `json:"ndcg" db:"cst_ndcg"`
	StalenessDays         *int       `json:"stalenessDays" db:"cst_staleness_days"`
	LastRefreshAt         *time.Time `json:"lastRefreshAt" db:"cst_last_refresh_at"`
	CurriculumCoveragePct *float64   `json:"curriculumCoveragePct" db:"cst_curriculum_coverage_pct"`
	CreatedAt             time.Time  `json:"createdAt" db:"cst_created_at"`
	UpdatedAt             time.Time  `json:"updatedAt" db:"cst_updated_at"`
}

// TopChunkByUsage for analytics queries
type TopChunkByUsage struct {
	ChunkID    int        `json:"chunkId" db:"chk_id"`
	Content    string     `json:"content" db:"chk_content"`
	DocTitle   string     `json:"docTitle" db:"doc_title"`
	UsageCount int        `json:"usageCount" db:"usage_count"`
	LastUsedAt *time.Time `json:"lastUsedAt" db:"last_used_at"`
	F1Score    *float64   `json:"f1Score" db:"f1_score"`
}

type IncrementChunkUsageResult struct {
	dal.DbResult
}

type UpdateChunkQualityMetricsParams struct {
	ChunkID      int
	PrecisionAtK *float64
	RecallAtK    *float64
	F1AtK        *float64
	MRR          *float64
	MAP          *float64
	NDCG         *float64
}

type UpdateChunkQualityMetricsResult struct {
	dal.DbResult
}

type UpdateChunkStalenessParams struct {
	ChunkID       int
	StalenessDays int
}

type UpdateChunkStalenessResult struct {
	dal.DbResult
}

// Chunk Statistics Repository & UseCase Interfaces
type ChunkStatisticsRepository interface {
	GetByChunk(ctx context.Context, chunkID int) (*ChunkStatistics, error)
	GetTopByUsage(ctx context.Context, limit int) ([]TopChunkByUsage, error)
	IncrementUsage(ctx context.Context, chunkID int) (*IncrementChunkUsageResult, error)
	UpdateQualityMetrics(ctx context.Context, params UpdateChunkQualityMetricsParams) (*UpdateChunkQualityMetricsResult, error)
	UpdateStaleness(ctx context.Context, params UpdateChunkStalenessParams) (*UpdateChunkStalenessResult, error)
}

type ChunkStatisticsUseCase interface {
	GetByChunk(ctx context.Context, chunkID int) Result[*ChunkStatistics]
	GetTopByUsage(ctx context.Context, limit int) Result[[]TopChunkByUsage]
	IncrementUsage(ctx context.Context, chunkID int) Result[Data]
	UpdateQualityMetrics(ctx context.Context, params UpdateChunkQualityMetricsParams) Result[Data]
	UpdateStaleness(ctx context.Context, params UpdateChunkStalenessParams) Result[Data]
}
