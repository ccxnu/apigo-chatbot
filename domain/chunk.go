package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
	"github.com/pgvector/pgvector-go"
)

type Chunk struct {
	ID         int        `json:"id" db:"chk_id"`
	DocumentID int        `json:"documentId" db:"chk_fk_document"`
	Content    string     `json:"content" db:"chk_content"`
	Embedding  *[]float32 `json:"embedding,omitempty" db:"chk_embedding"`
	CreatedAt  time.Time  `json:"createdAt" db:"chk_created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"chk_updated_at"`
}

// ChunkWithSimilarity extends Chunk for similarity search results
type ChunkWithSimilarity struct {
	ID              int     `json:"id" db:"chk_id"`
	DocumentID      int     `json:"documentId" db:"chk_fk_document"`
	Content         string  `json:"content" db:"chk_content"`
	SimilarityScore float64 `json:"similarityScore" db:"similarity_score"`
	DocTitle        string  `json:"docTitle" db:"doc_title"`
	DocCategory     string  `json:"docCategory" db:"doc_category"`
}

// Chunk Repository Params & Results
type CreateChunkParams struct {
	DocumentID int
	Content    string
	Embedding  *[]float32
}

type CreateChunkResult struct {
	dal.DbResult
	ChunkID int `json:"chunkId" db:"o_chk_id"`
}

type UpdateChunkEmbeddingParams struct {
	ChunkID   int
	Embedding []float32
}

type UpdateChunkEmbeddingResult struct {
	dal.DbResult
}

type DeleteChunkResult struct {
	dal.DbResult
}

type BulkCreateChunksParams struct {
	DocumentID int
	Contents   []string
	Embeddings *[]pgvector.Vector
}

type BulkCreateChunksResult struct {
	dal.DbResult
	ChunksCreated int `json:"chunksCreated" db:"o_chunks_created"`
}

type SimilaritySearchParams struct {
	QueryEmbedding []float32
	Limit          int
	MinSimilarity  float64
}

// Chunk Repository & UseCase Interfaces
type ChunkRepository interface {
	GetByDocument(ctx context.Context, docID int) ([]Chunk, error)
	GetByID(ctx context.Context, chunkID int) (*Chunk, error)
	SimilaritySearch(ctx context.Context, params SimilaritySearchParams) ([]ChunkWithSimilarity, error)
	Create(ctx context.Context, params CreateChunkParams) (*CreateChunkResult, error)
	UpdateEmbedding(ctx context.Context, params UpdateChunkEmbeddingParams) (*UpdateChunkEmbeddingResult, error)
	Delete(ctx context.Context, chunkID int) (*DeleteChunkResult, error)
	BulkCreate(ctx context.Context, params BulkCreateChunksParams) (*BulkCreateChunksResult, error)
}

type ChunkUseCase interface {
	GetByDocument(ctx context.Context, docID int) Result[[]Chunk]
	GetByID(ctx context.Context, chunkID int) Result[*Chunk]
	SimilaritySearch(ctx context.Context, queryText string, limit int, minSimilarity float64) Result[[]ChunkWithSimilarity]
	Create(ctx context.Context, documentID int, content string) Result[Data]
	UpdateContent(ctx context.Context, chunkID int, content string) Result[Data]
	Delete(ctx context.Context, chunkID int) Result[Data]
	BulkCreate(ctx context.Context, documentID int, contents []string) Result[Data]
}
