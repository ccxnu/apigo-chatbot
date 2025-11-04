package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetChunksByDocument          = "fn_get_chunks_by_document"
	fnGetChunkByID                 = "fn_get_chunk_by_id"
	fnSimilaritySearchChunks       = "fn_similarity_search_chunks"
	fnSimilaritySearchChunksHybrid = "fn_similarity_search_chunks_hybrid"
	// Stored Procedures (Writes)
	spCreateChunk          = "sp_create_chunk"
	spUpdateChunkEmbedding = "sp_update_chunk_embedding"
	spDeleteChunk          = "sp_delete_chunk"
	spBulkCreateChunks     = "sp_bulk_create_chunks"
)

type chunkRepository struct {
	dal *dal.DAL
}

func NewChunkRepository(dal *dal.DAL) d.ChunkRepository {
	return &chunkRepository{
		dal: dal,
	}
}

// GetByDocument retrieves all chunks for a specific document
func (r *chunkRepository) GetByDocument(ctx context.Context, docID int) ([]d.Chunk, error) {
	chunks, err := dal.QueryRows[d.Chunk](r.dal, ctx, fnGetChunksByDocument, docID)

	if err != nil {
		return nil, fmt.Errorf("failed to get chunks by document via %s: %w", fnGetChunksByDocument, err)
	}

	return chunks, nil

}

// GetByID retrieves a single chunk by ID
func (r *chunkRepository) GetByID(ctx context.Context, chunkID int) (*d.Chunk, error) {
	chunks, err := dal.QueryRows[d.Chunk](r.dal, ctx, fnGetChunkByID, chunkID)

	if err != nil {
		return nil, fmt.Errorf("failed to get chunk by id via %s: %w", fnGetChunkByID, err)
	}

	if len(chunks) == 0 {
		return nil, nil // No technical error, just no row found
	}

	return &chunks[0], nil
}

// SimilaritySearch performs vector similarity search for RAG
func (r *chunkRepository) SimilaritySearch(ctx context.Context, params d.SimilaritySearchParams) ([]d.ChunkWithSimilarity, error) {
	chunks, err := dal.QueryRows[d.ChunkWithSimilarity](
		r.dal,
		ctx,
		fnSimilaritySearchChunks,
		params.QueryEmbedding,
		params.Limit,
		params.MinSimilarity,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search via %s: %w", fnSimilaritySearchChunks, err)
	}

	return chunks, nil
}

// HybridSearch performs hybrid search combining vector similarity and full-text search
func (r *chunkRepository) HybridSearch(ctx context.Context, params d.HybridSearchParams) ([]d.ChunkWithHybridSimilarity, error) {
	chunks, err := dal.QueryRows[d.ChunkWithHybridSimilarity](
		r.dal,
		ctx,
		fnSimilaritySearchChunksHybrid,
		params.QueryEmbedding,
		params.QueryText,
		params.Limit,
		params.MinSimilarity,
		params.KeywordWeight,
		params.Category, // Pass category filter (can be nil)
	)

	if err != nil {
		return nil, fmt.Errorf("failed to perform hybrid search via %s: %w", fnSimilaritySearchChunksHybrid, err)
	}

	return chunks, nil
}

// Create creates a new chunk
func (r *chunkRepository) Create(ctx context.Context, params d.CreateChunkParams) (*d.CreateChunkResult, error) {
	result, err := dal.ExecProc[d.CreateChunkResult](
		r.dal,
		ctx,
		spCreateChunk,
		params.DocumentID,
		params.Content,
		params.Embedding,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateChunk, err)
	}

	return result, nil
}

// UpdateEmbedding updates the embedding for a chunk
func (r *chunkRepository) UpdateEmbedding(ctx context.Context, params d.UpdateChunkEmbeddingParams) (*d.UpdateChunkEmbeddingResult, error) {
	result, err := dal.ExecProc[d.UpdateChunkEmbeddingResult](
		r.dal,
		ctx,
		spUpdateChunkEmbedding,
		params.ChunkID,
		params.Embedding,
		params.Content,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateChunkEmbedding, err)
	}

	return result, nil
}

func (r *chunkRepository) Delete(ctx context.Context, chunkID int) (*d.DeleteChunkResult, error) {
	result, err := dal.ExecProc[d.DeleteChunkResult](
		r.dal,
		ctx,
		spDeleteChunk,
		chunkID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spDeleteChunk, err)
	}

	return result, nil
}

// BulkCreate creates multiple chunks at once
func (r *chunkRepository) BulkCreate(ctx context.Context, params d.BulkCreateChunksParams) (*d.BulkCreateChunksResult, error) {
	result, err := dal.ExecProc[d.BulkCreateChunksResult](
		r.dal,
		ctx,
		spBulkCreateChunks,
		params.DocumentID,
		params.Contents,
		params.Embeddings,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spBulkCreateChunks, err)
	}

	return result, nil
}
