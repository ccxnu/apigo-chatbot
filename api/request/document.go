package request

import (
	"api-chatbot/domain"
)

// Document Requests

type GetAllDocumentsRequest struct {
	domain.Base
	Limit  int `json:"limit" validate:"omitempty,gte=1,lte=1000"`
	Offset int `json:"offset" validate:"omitempty,gte=0"`
}

type GetDocumentByIDRequest struct {
	domain.Base
	DocID int `json:"docId" validate:"required,gte=1"`
}

type GetDocumentsByCategoryRequest struct {
	domain.Base
	Category string `json:"category" validate:"required"`
	Limit    int    `json:"limit" validate:"omitempty,gte=1,lte=1000"`
	Offset   int    `json:"offset" validate:"omitempty,gte=0"`
}

type SearchDocumentsByTitleRequest struct {
	domain.Base
	TitlePattern string `json:"titlePattern" validate:"required,min=1"`
	Limit        int    `json:"limit" validate:"omitempty,gte=1,lte=1000"`
}

type CreateDocumentRequest struct {
	domain.Base
	Category    string  `json:"category" validate:"required"`
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Summary     *string `json:"summary" validate:"omitempty,max=5000"`
	Source      *string `json:"source" validate:"omitempty,max=500"`
	PublishedAt *string `json:"publishedAt" validate:"omitempty"` // ISO 8601 timestamp string
}

type UpdateDocumentRequest struct {
	domain.Base
	DocID       int     `json:"docId" validate:"required,gte=1"`
	Category    string  `json:"category" validate:"required"`
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Summary     *string `json:"summary" validate:"omitempty,max=5000"`
	Source      *string `json:"source" validate:"omitempty,max=500"`
	PublishedAt *string `json:"publishedAt" validate:"omitempty"` // ISO 8601 timestamp string
}

type DeleteDocumentRequest struct {
	domain.Base
	DocID int `json:"docId" validate:"required,gte=1"`
}

type UploadPDFDocumentRequest struct {
	domain.Base
	Category    string  `json:"category" validate:"required"`
	Title       string  `json:"title" validate:"required,min=1,max=200"`
	Source      *string `json:"source" validate:"omitempty,max=500"`
	FileBase64  string  `json:"fileBase64" validate:"required"`
	ChunkSize   *int    `json:"chunkSize" validate:"omitempty,gte=100,lte=5000"`
	ChunkOverlap *int   `json:"chunkOverlap" validate:"omitempty,gte=0,lte=500"`
}
