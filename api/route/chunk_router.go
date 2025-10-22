package route

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

type GetChunksByDocumentResponse struct {
	Body d.Result[[]d.Chunk]
}

type GetChunkByIDResponse struct {
	Body d.Result[*d.Chunk]
}

type SimilaritySearchResponse struct {
	Body d.Result[[]d.ChunkWithSimilarity]
}

type HybridSearchResponse struct {
	Body d.Result[[]d.ChunkWithHybridSimilarity]
}

type CreateChunkResponse struct {
	Body d.Result[d.Data]
}

type UpdateChunkContentResponse struct {
	Body d.Result[d.Data]
}

type DeleteChunkResponse struct {
	Body d.Result[d.Data]
}

type BulkCreateChunksResponse struct {
	Body d.Result[d.Data]
}

func NewChunkRouter(chunkUseCase d.ChunkUseCase, mux *http.ServeMux, humaAPI huma.API) {
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-chunks-by-document",
		Method:      "POST",
		Path:        "/api/v1/chunks/get-by-document",
		Summary:     "Get chunks by document",
		Description: "Retrieves all chunks for a specific document",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.GetChunksByDocumentRequest
	}) (*GetChunksByDocumentResponse, error) {
		result := chunkUseCase.GetByDocument(ctx, input.Body.DocID)
		return &GetChunksByDocumentResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-chunk-by-id",
		Method:      "POST",
		Path:        "/api/v1/chunks/get-by-id",
		Summary:     "Get chunk by ID",
		Description: "Retrieves a specific chunk by its ID including embedding",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.GetChunkByIDRequest
	}) (*GetChunkByIDResponse, error) {
		result := chunkUseCase.GetByID(ctx, input.Body.ChunkID)
		return &GetChunkByIDResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "similarity-search-chunks",
		Method:      "POST",
		Path:        "/api/v1/chunks/similarity-search",
		Summary:     "Vector similarity search",
		Description: "Performs semantic search using text query (automatically converted to embedding). Returns top K most similar chunks ordered by cosine similarity.",
		Tags:        []string{"Chunks", "RAG"},
	}, func(ctx context.Context, input *struct { Body request.SimilaritySearchRequest }) (*SimilaritySearchResponse, error) {
		result := chunkUseCase.SimilaritySearch(ctx, input.Body.QueryText, input.Body.Limit, input.Body.MinSimilarity)
		return &SimilaritySearchResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "hybrid-search-chunks",
		Method:      "POST",
		Path:        "/api/v1/chunks/hybrid-search",
		Summary:     "Hybrid search (vector + full-text)",
		Description: "Performs hybrid search combining semantic similarity (embeddings) and keyword matching (full-text search). Returns top K chunks ordered by combined score.",
		Tags:        []string{"Chunks", "RAG"},
	}, func(ctx context.Context, input *struct {
		Body request.HybridSearchRequest
	}) (*HybridSearchResponse, error) {
		result := chunkUseCase.HybridSearch(ctx, input.Body.QueryText, input.Body.Limit, input.Body.MinSimilarity, input.Body.KeywordWeight)
		return &HybridSearchResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "create-chunk",
		Method:      "POST",
		Path:        "/api/v1/chunks/create",
		Summary:     "Create chunk",
		Description: "Creates a new chunk from text content (embedding generated automatically). Initializes statistics record.",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.CreateChunkRequest
	}) (*CreateChunkResponse, error) {
		result := chunkUseCase.Create(ctx, input.Body.DocumentID, input.Body.Content)
		return &CreateChunkResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-chunk-content",
		Method:      "POST",
		Path:        "/api/v1/chunks/update-content",
		Summary:     "Update chunk content",
		Description: "Updates the content and regenerates the embedding vector automatically",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdateChunkContentRequest
	}) (*UpdateChunkContentResponse, error) {
		result := chunkUseCase.UpdateContent(ctx, input.Body.ChunkID, input.Body.Content)
		return &UpdateChunkContentResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "delete-chunk",
		Method:      "POST",
		Path:        "/api/v1/chunks/delete",
		Summary:     "Delete chunk",
		Description: "Hard deletes a chunk (cascades to statistics)",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.DeleteChunkRequest
	}) (*DeleteChunkResponse, error) {
		result := chunkUseCase.Delete(ctx, input.Body.ChunkID)
		return &DeleteChunkResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "bulk-create-chunks",
		Method:      "POST",
		Path:        "/api/v1/chunks/bulk-create",
		Summary:     "Bulk create chunks",
		Description: "Creates multiple chunks from text contents (embeddings generated automatically). Efficient for document ingestion.",
		Tags:        []string{"Chunks"},
	}, func(ctx context.Context, input *struct {
		Body request.BulkCreateChunksRequest
	}) (*BulkCreateChunksResponse, error) {
		result := chunkUseCase.BulkCreate(ctx, input.Body.DocumentID, input.Body.Contents)
		return &BulkCreateChunksResponse{Body: result}, nil
	})
}
