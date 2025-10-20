package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
)

type chunkUseCase struct {
	chunkRepo        d.ChunkRepository
	cache            d.ParameterCache
	embeddingService d.EmbeddingService
	contextTimeout   time.Duration
}

func NewChunkUseCase(
	chunkRepo d.ChunkRepository,
	cache d.ParameterCache,
	embeddingService d.EmbeddingService,
	timeout time.Duration,
) d.ChunkUseCase {
	return &chunkUseCase{
		chunkRepo:        chunkRepo,
		cache:            cache,
		embeddingService: embeddingService,
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

	return d.Success(chunks)
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
