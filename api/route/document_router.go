package route

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

// Huma response types for documents
type GetAllDocumentsResponse struct {
	Body d.Result[[]d.Document]
}

type GetDocumentByIDResponse struct {
	Body d.Result[*d.Document]
}

type GetDocumentsByCategoryResponse struct {
	Body d.Result[[]d.Document]
}

type SearchDocumentsByTitleResponse struct {
	Body d.Result[[]d.Document]
}

type CreateDocumentResponse struct {
	Body d.Result[d.Data]
}

type UpdateDocumentResponse struct {
	Body d.Result[d.Data]
}

type DeleteDocumentResponse struct {
	Body d.Result[d.Data]
}

type UploadPDFDocumentResponse struct {
	Body d.Result[d.Data]
}

func NewDocumentRouter(docUseCase d.DocumentUseCase, mux *http.ServeMux, humaAPI huma.API) {
	// Huma documented routes with /api/v1/ prefix
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-all-documents",
		Method:      "POST",
		Path:        "/api/v1/documents/get-all",
		Summary:     "Get all documents",
		Description: "Retrieves all active documents with pagination",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.GetAllDocumentsRequest
	}) (*GetAllDocumentsResponse, error) {
		result := docUseCase.GetAll(ctx, input.Body.Limit, input.Body.Offset)
		return &GetAllDocumentsResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-document-by-id",
		Method:      "POST",
		Path:        "/api/v1/documents/get-by-id",
		Summary:     "Get document by ID",
		Description: "Retrieves a specific document by its ID",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.GetDocumentByIDRequest
	}) (*GetDocumentByIDResponse, error) {
		result := docUseCase.GetByID(ctx, input.Body.DocID)
		return &GetDocumentByIDResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-documents-by-category",
		Method:      "POST",
		Path:        "/api/v1/documents/get-by-category",
		Summary:     "Get documents by category",
		Description: "Retrieves documents filtered by category with pagination",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.GetDocumentsByCategoryRequest
	}) (*GetDocumentsByCategoryResponse, error) {
		result := docUseCase.GetByCategory(ctx, input.Body.Category, input.Body.Limit, input.Body.Offset)
		return &GetDocumentsByCategoryResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "search-documents-by-title",
		Method:      "POST",
		Path:        "/api/v1/documents/search-by-title",
		Summary:     "Search documents by title",
		Description: "Searches documents by title pattern (case-insensitive)",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.SearchDocumentsByTitleRequest
	}) (*SearchDocumentsByTitleResponse, error) {
		result := docUseCase.SearchByTitle(ctx, input.Body.TitlePattern, input.Body.Limit)
		return &SearchDocumentsByTitleResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "create-document",
		Method:      "POST",
		Path:        "/api/v1/documents/create",
		Summary:     "Create document",
		Description: "Creates a new document in the knowledge base",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.CreateDocumentRequest
	}) (*CreateDocumentResponse, error) {
		// Parse publishedAt if provided
		var publishedAt *time.Time
		if input.Body.PublishedAt != nil && *input.Body.PublishedAt != "" {
			t, err := time.Parse(time.RFC3339, *input.Body.PublishedAt)
			if err == nil {
				publishedAt = &t
			}
		}

		params := d.CreateDocumentParams{
			Category:    input.Body.Category,
			Title:       input.Body.Title,
			Summary:     input.Body.Summary,
			Source:      input.Body.Source,
			PublishedAt: publishedAt,
		}
		result := docUseCase.Create(ctx, params)
		return &CreateDocumentResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-document",
		Method:      "POST",
		Path:        "/api/v1/documents/update",
		Summary:     "Update document",
		Description: "Updates an existing document",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdateDocumentRequest
	}) (*UpdateDocumentResponse, error) {
		// Parse publishedAt if provided
		var publishedAt *time.Time
		if input.Body.PublishedAt != nil && *input.Body.PublishedAt != "" {
			t, err := time.Parse(time.RFC3339, *input.Body.PublishedAt)
			if err == nil {
				publishedAt = &t
			}
		}

		params := d.UpdateDocumentParams{
			DocID:       input.Body.DocID,
			Category:    input.Body.Category,
			Title:       input.Body.Title,
			Summary:     input.Body.Summary,
			Source:      input.Body.Source,
			PublishedAt: publishedAt,
		}
		result := docUseCase.Update(ctx, params)
		return &UpdateDocumentResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "delete-document",
		Method:      "POST",
		Path:        "/api/v1/documents/delete",
		Summary:     "Delete document",
		Description: "Soft deletes a document (sets active = false)",
		Tags:        []string{"Documents"},
	}, func(ctx context.Context, input *struct {
		Body request.DeleteDocumentRequest
	}) (*DeleteDocumentResponse, error) {
		result := docUseCase.Delete(ctx, input.Body.DocID)
		return &DeleteDocumentResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID:  "upload-pdf-document",
		Method:       "POST",
		Path:         "/api/v1/documents/upload-pdf",
		Summary:      "Upload PDF document",
		Description:  "Uploads a PDF file (base64 encoded), extracts text using OCR, creates a document, and generates chunks",
		Tags:         []string{"Documents"},
		MaxBodyBytes: 20 * 1024 * 1024, // 20MB limit for PDF uploads
	}, func(ctx context.Context, input *struct {
		Body request.UploadPDFDocumentRequest
	}) (*UploadPDFDocumentResponse, error) {
		// Set default chunk size and overlap if not provided
		chunkSize := 1000
		chunkOverlap := 200
		if input.Body.ChunkSize != nil {
			chunkSize = *input.Body.ChunkSize
		}
		if input.Body.ChunkOverlap != nil {
			chunkOverlap = *input.Body.ChunkOverlap
		}

		params := d.UploadPDFDocumentParams{
			Category:     input.Body.Category,
			Title:        input.Body.Title,
			Source:       input.Body.Source,
			FileBase64:   input.Body.FileBase64,
			ChunkSize:    chunkSize,
			ChunkOverlap: chunkOverlap,
		}
		result := docUseCase.UploadPDF(ctx, params)
		return &UploadPDFDocumentResponse{Body: result}, nil
	})
}
