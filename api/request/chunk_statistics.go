package request

import (
	"api-chatbot/domain"
)

// Chunk Statistics Requests

type GetChunkStatisticsRequest struct {
	domain.Base
	ChunkID int `json:"chunkId" validate:"required,gte=1"`
}

type GetTopChunksByUsageRequest struct {
	domain.Base
	Limit int `json:"limit" validate:"omitempty,gte=1,lte=100"`
}

type IncrementChunkUsageRequest struct {
	domain.Base
	ChunkID int `json:"chunkId" validate:"required,gte=1"`
}

type UpdateChunkQualityMetricsRequest struct {
	domain.Base
	ChunkID      int      `json:"chunkId" validate:"required,gte=1"`
	PrecisionAtK *float64 `json:"precisionAtK" validate:"omitempty,gte=0,lte=1"`
	RecallAtK    *float64 `json:"recallAtK" validate:"omitempty,gte=0,lte=1"`
	F1AtK        *float64 `json:"f1AtK" validate:"omitempty,gte=0,lte=1"`
	MRR          *float64 `json:"mrr" validate:"omitempty,gte=0,lte=1"`
	MAP          *float64 `json:"map" validate:"omitempty,gte=0,lte=1"`
	NDCG         *float64 `json:"ndcg" validate:"omitempty,gte=0,lte=1"`
}

type UpdateChunkStalenessRequest struct {
	domain.Base
	ChunkID       int `json:"chunkId" validate:"required,gte=1"`
	StalenessDays int `json:"stalenessDays" validate:"required,gte=0"`
}
