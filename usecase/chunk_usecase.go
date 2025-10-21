package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/metrics"
)

type chunkUseCase struct {
	chunkRepo        d.ChunkRepository
	statsRepo        d.ChunkStatisticsRepository
	cache            d.ParameterCache
	embeddingService d.EmbeddingService
	metricsCalc      *metrics.RAGMetrics
	contextTimeout   time.Duration
}

func NewChunkUseCase(
	chunkRepo d.ChunkRepository,
	statsRepo d.ChunkStatisticsRepository,
	cache d.ParameterCache,
	embeddingService d.EmbeddingService,
	timeout time.Duration,
) d.ChunkUseCase {
	return &chunkUseCase{
		chunkRepo:        chunkRepo,
		statsRepo:        statsRepo,
		cache:            cache,
		embeddingService: embeddingService,
		metricsCalc:      metrics.NewRAGMetrics(),
		contextTimeout:   timeout,
	}
}

func (u *chunkUseCase) GetByDocument(c context.Context, docID int) d.Result[[]d.Chunk] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	chunks, err := u.chunkRepo.GetByDocument(ctx, docID)
	if err != nil {
		return d.Error[[]d.Chunk](u.cache, "ERR_INTERNAL_DB")
	}

	return d.Success(chunks)
}

func (u *chunkUseCase) GetByID(c context.Context, chunkID int) d.Result[*d.Chunk] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	chunk, err := u.chunkRepo.GetByID(ctx, chunkID)
	if err != nil {
		return d.Error[*d.Chunk](u.cache, "ERR_INTERNAL_DB")
	}

	if chunk == nil {
		return d.Error[*d.Chunk](u.cache, "ERR_CHUNK_NOT_FOUND")
	}

	return d.Success(chunk)
}

func (u *chunkUseCase) SimilaritySearch(c context.Context, queryText string, limit int, minSimilarity float64) d.Result[[]d.ChunkWithSimilarity] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Generate embedding from query text
	queryEmbedding, err := u.embeddingService.GenerateEmbedding(ctx, queryText)
	if err != nil {
		return d.Error[[]d.ChunkWithSimilarity](u.cache, "ERR_EMBEDDING_GENERATION")
	}

	// Create params with generated embedding
	params := d.SimilaritySearchParams{
		QueryEmbedding: queryEmbedding,
		Limit:          limit,
		MinSimilarity:  minSimilarity,
	}

	chunks, err := u.chunkRepo.SimilaritySearch(ctx, params)
	if err != nil {
		return d.Error[[]d.ChunkWithSimilarity](u.cache, "ERR_INTERNAL_DB")
	}

	// Automatically update statistics for each retrieved chunk
	// This happens asynchronously to not block the response
	go u.updateChunkStatistics(chunks)

	return d.Success(chunks)
}

// updateChunkStatistics updates usage statistics and calculates quality metrics
func (u *chunkUseCase) updateChunkStatistics(chunks []d.ChunkWithSimilarity) {
	asyncCtx, asyncCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer asyncCancel()

	// Build metrics input from retrieved chunks
	metricsChunks := make([]metrics.RetrievedChunk, len(chunks))
	relevanceThreshold := 0.75 // Chunks with similarity >= 0.75 are considered relevant

	for i, chunk := range chunks {
		metricsChunks[i] = metrics.RetrievedChunk{
			ChunkID:         chunk.ID,
			SimilarityScore: chunk.SimilarityScore,
			Position:        i + 1, // 1-based position
			IsRelevant:      metrics.EstimateRelevanceFromSimilarity(chunk.SimilarityScore, relevanceThreshold),
		}
	}

	// Count total relevant chunks (those above threshold)
	totalRelevant := 0
	for _, mc := range metricsChunks {
		if mc.IsRelevant {
			totalRelevant++
		}
	}

	// Calculate all quality metrics
	result := u.metricsCalc.CalculateAllMetrics(metricsChunks, totalRelevant)

	// Update each chunk's statistics
	for i, chunk := range chunks {
		chunkID := chunk.ID
		position := i + 1

		// Increment usage count
		_, _ = u.statsRepo.IncrementUsage(asyncCtx, chunkID)

		// Only update quality metrics for chunks in top positions (more accurate)
		// Update metrics more aggressively for top 3 results
		if position <= 3 {
			// Calculate individual chunk metrics
			// Use exponential moving average to smooth metrics over time
			params := d.UpdateChunkQualityMetricsParams{
				ChunkID: chunkID,
			}

			// For top-ranked chunks, update with calculated metrics
			if metricsChunks[i].IsRelevant {
				params.PrecisionAtK = &result.PrecisionAtK
				params.RecallAtK = &result.RecallAtK
				params.F1AtK = &result.F1AtK
				params.MRR = &result.MRR
				params.MAP = &result.MAP
				params.NDCG = &result.NDCG

				_, _ = u.statsRepo.UpdateQualityMetrics(asyncCtx, params)
			}
		}
	}
}

func (u *chunkUseCase) Create(c context.Context, documentID int, content string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Generate embedding from content
	embedding, err := u.embeddingService.GenerateEmbedding(ctx, content)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_EMBEDDING_GENERATION")
	}

	// Create params with generated embedding
	params := d.CreateChunkParams{
		DocumentID: documentID,
		Content:    content,
		Embedding:  &embedding,
	}

	result, err := u.chunkRepo.Create(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	return d.Success(d.Data{"chunkId": result.ChunkID})
}

func (u *chunkUseCase) UpdateContent(c context.Context, chunkID int, content string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Generate new embedding from updated content
	embedding, err := u.embeddingService.GenerateEmbedding(ctx, content)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_EMBEDDING_GENERATION")
	}

	// Create params with generated embedding
	params := d.UpdateChunkEmbeddingParams{
		ChunkID:   chunkID,
		Embedding: embedding,
	}

	result, err := u.chunkRepo.UpdateEmbedding(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *chunkUseCase) Delete(c context.Context, chunkID int) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.chunkRepo.Delete(ctx, chunkID)
	if err != nil || result == nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *chunkUseCase) BulkCreate(c context.Context, documentID int, contents []string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Generate embeddings for all contents (batch processing)
	embeddings, err := u.embeddingService.GenerateEmbeddings(ctx, contents)
	if err != nil {
		return d.Error[d.Data](u.cache, "ERR_EMBEDDING_GENERATION")
	}

	// Create params with generated embeddings
	params := d.BulkCreateChunksParams{
		DocumentID: documentID,
		Contents:   contents,
		Embeddings: &embeddings,
	}

	result, err := u.chunkRepo.BulkCreate(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.cache, result.Code)
	}

	return d.Success(d.Data{"chunksCreated": result.ChunksCreated})
}
