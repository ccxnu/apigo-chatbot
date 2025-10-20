package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetChunksByDocument    = "fn_get_chunks_by_document"
	fnGetChunkByID           = "fn_get_chunk_by_id"
	fnSimilaritySearchChunks = "fn_similarity_search_chunks"
	// Stored Procedures (Writes)
	spCreateChunk          = "sp_create_chunk"
	spUpdateChunkEmbedding = "sp_update_chunk_embedding"
	spDeleteChunk          = "sp_delete_chunk"
	spBulkCreateChunks     = "sp_bulk_create_chunks"
)

type chunkRepository struct {
	dal *dal.DAL
}

func NewChunkRepository(dal *dal.DAL) domain.ChunkRepository {
	return &chunkRepository{
		dal: dal,
	}
}

// GetByDocument retrieves all chunks for a specific document
func (r *chunkRepository) GetByDocument(ctx context.Context, docID int) ([]domain.Chunk, error) {
	chunks, err := dal.QueryRows[domain.Chunk](r.dal, ctx, fnGetChunksByDocument, docID)

	if err != nil {
		return nil, fmt.Errorf("failed to get chunks by document via %s: %w", fnGetChunksByDocument, err)
	}

	return chunks, nil

}

// GetByID retrieves a single chunk by ID
func (r *chunkRepository) GetByID(ctx context.Context, chunkID int) (*domain.Chunk, error) {
	chunks, err := dal.QueryRows[domain.Chunk](r.dal, ctx, fnGetChunkByID, chunkID)

	if err != nil {
		return nil, fmt.Errorf("failed to get chunk by id via %s: %w", fnGetChunkByID, err)
	}

	if len(chunks) == 0 {
		return nil, nil // No technical error, just no row found
	}

	return &chunks[0], nil
}

// SimilaritySearch performs vector similarity search for RAG
func (r *chunkRepository) SimilaritySearch(ctx context.Context, params domain.SimilaritySearchParams) ([]domain.ChunkWithSimilarity, error) {
	chunks, err := dal.QueryRows[domain.ChunkWithSimilarity](
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

// Create creates a new chunk
func (r *chunkRepository) Create(ctx context.Context, params domain.CreateChunkParams) (*domain.CreateChunkResult, error) {
	result, err := dal.ExecProc[domain.CreateChunkResult](
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
func (r *chunkRepository) UpdateEmbedding(ctx context.Context, params domain.UpdateChunkEmbeddingParams) (*domain.UpdateChunkEmbeddingResult, error) {
	result, err := dal.ExecProc[domain.UpdateChunkEmbeddingResult](
		r.dal,
		ctx,
		spUpdateChunkEmbedding,
		params.ChunkID,
		params.Embedding,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateChunkEmbedding, err)
	}

	return result, nil
}

func (r *chunkRepository) Delete(ctx context.Context, chunkID int) (*domain.DeleteChunkResult, error) {
	result, err := dal.ExecProc[domain.DeleteChunkResult](
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
func (r *chunkRepository) BulkCreate(ctx context.Context, params domain.BulkCreateChunksParams) (*domain.BulkCreateChunksResult, error) {
	result, err := dal.ExecProc[domain.BulkCreateChunksResult](
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
