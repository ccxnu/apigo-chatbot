package domain

import "context"

// EmbeddingService generates vector embeddings from text
type EmbeddingService interface {
	// GenerateEmbedding generates a single embedding vector from text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateEmbeddings generates embeddings for multiple texts (batch)
	GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}
