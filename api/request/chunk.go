package request

import (
	"api-chatbot/domain"
)

// Chunk Requests

type GetChunksByDocumentRequest struct {
	domain.Base
	DocID int `json:"docId" validate:"required,gte=1" doc:"Document ID to retrieve chunks from"`
}

type GetChunkByIDRequest struct {
	domain.Base
	ChunkID int `json:"chunkId" validate:"required,gte=1" doc:"Chunk ID to retrieve"`
}

type SimilaritySearchRequest struct {
	domain.Base
	QueryText     string  `json:"queryText" validate:"required,min=1" doc:"Text query (automatically converted to embedding)"`
	Limit         int     `json:"limit" validate:"omitempty,gte=1,lte=100" doc:"Maximum number of results (default: 10)"`
	MinSimilarity float64 `json:"minSimilarity" validate:"omitempty,gte=0,lte=1" doc:"Minimum similarity score 0-1 (default: 0.7)"`
}

type CreateChunkRequest struct {
	domain.Base
	DocumentID int    `json:"documentId" validate:"required,gte=1" doc:"Document ID this chunk belongs to"`
	Content    string `json:"content" validate:"required,min=1" doc:"Chunk text content (embedding generated automatically)"`
}

type UpdateChunkContentRequest struct {
	domain.Base
	ChunkID int    `json:"chunkId" validate:"required,gte=1" doc:"Chunk ID to update"`
	Content string `json:"content" validate:"required,min=1" doc:"New content (embedding regenerated automatically)"`
}

type DeleteChunkRequest struct {
	domain.Base
	ChunkID int `json:"chunkId" validate:"required,gte=1" doc:"Chunk ID to delete"`
}

type BulkCreateChunksRequest struct {
	domain.Base
	DocumentID int      `json:"documentId" validate:"required,gte=1" doc:"Document ID for all chunks"`
	Contents   []string `json:"contents" validate:"required,min=1,dive,required,min=1" doc:"Array of chunk contents (embeddings generated automatically)"`
}
