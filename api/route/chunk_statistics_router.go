package route

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

// Huma response types for chunk statistics
type GetChunkStatisticsResponse struct {
	Body d.Result[*d.ChunkStatistics]
}

type GetTopChunksByUsageResponse struct {
	Body d.Result[[]d.TopChunkByUsage]
}

type IncrementChunkUsageResponse struct {
	Body d.Result[d.Data]
}

type UpdateChunkQualityMetricsResponse struct {
	Body d.Result[d.Data]
}

type UpdateChunkStalenessResponse struct {
	Body d.Result[d.Data]
}

func NewChunkStatisticsRouter(statsUseCase d.ChunkStatisticsUseCase, mux *http.ServeMux, humaAPI huma.API) {
	// Huma documented routes with /api/v1/ prefix
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-chunk-statistics",
		Method:      "POST",
		Path:        "/api/v1/chunk-statistics/get-by-chunk",
		Summary:     "Get chunk statistics",
		Description: "Retrieves all statistics and quality metrics for a specific chunk",
		Tags:        []string{"Chunk Statistics", "Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetChunkStatisticsRequest
	}) (*GetChunkStatisticsResponse, error) {
		result := statsUseCase.GetByChunk(ctx, input.Body.ChunkID)
		return &GetChunkStatisticsResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-top-chunks-by-usage",
		Method:      "POST",
		Path:        "/api/v1/chunk-statistics/get-top-by-usage",
		Summary:     "Get most used chunks",
		Description: "Retrieves the most frequently used chunks for analytics. Useful for identifying popular knowledge.",
		Tags:        []string{"Chunk Statistics", "Analytics"},
	}, func(ctx context.Context, input *struct {
		Body request.GetTopChunksByUsageRequest
	}) (*GetTopChunksByUsageResponse, error) {
		result := statsUseCase.GetTopByUsage(ctx, input.Body.Limit)
		return &GetTopChunksByUsageResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "increment-chunk-usage",
		Method:      "POST",
		Path:        "/api/v1/chunk-statistics/increment-usage",
		Summary:     "Increment chunk usage counter",
		Description: "Increments the usage count and updates last used timestamp for a chunk. Call this when a chunk is used in RAG responses.",
		Tags:        []string{"Chunk Statistics"},
	}, func(ctx context.Context, input *struct {
		Body request.IncrementChunkUsageRequest
	}) (*IncrementChunkUsageResponse, error) {
		result := statsUseCase.IncrementUsage(ctx, input.Body.ChunkID)
		return &IncrementChunkUsageResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-chunk-quality-metrics",
		Method:      "POST",
		Path:        "/api/v1/chunk-statistics/update-quality-metrics",
		Summary:     "Update RAG quality metrics",
		Description: "Updates quality and relevance metrics for a chunk. Metrics include Precision@K, Recall@K, F1, MRR, MAP, and NDCG. Only provided metrics are updated.",
		Tags:        []string{"Chunk Statistics", "RAG"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdateChunkQualityMetricsRequest
	}) (*UpdateChunkQualityMetricsResponse, error) {
		params := d.UpdateChunkQualityMetricsParams{
			ChunkID:      input.Body.ChunkID,
			PrecisionAtK: input.Body.PrecisionAtK,
			RecallAtK:    input.Body.RecallAtK,
			F1AtK:        input.Body.F1AtK,
			MRR:          input.Body.MRR,
			MAP:          input.Body.MAP,
			NDCG:         input.Body.NDCG,
		}
		result := statsUseCase.UpdateQualityMetrics(ctx, params)
		return &UpdateChunkQualityMetricsResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-chunk-staleness",
		Method:      "POST",
		Path:        "/api/v1/chunk-statistics/update-staleness",
		Summary:     "Update chunk staleness tracking",
		Description: "Updates the staleness metric for a chunk to track content freshness. Higher values indicate older content.",
		Tags:        []string{"Chunk Statistics"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdateChunkStalenessRequest
	}) (*UpdateChunkStalenessResponse, error) {
		params := d.UpdateChunkStalenessParams{
			ChunkID:       input.Body.ChunkID,
			StalenessDays: input.Body.StalenessDays,
		}
		result := statsUseCase.UpdateStaleness(ctx, params)
		return &UpdateChunkStalenessResponse{Body: result}, nil
	})
}
