package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
)

type Document struct {
	ID          int        `json:"id" db:"doc_id"`
	Category    string     `json:"category" db:"doc_category"`
	Title       string     `json:"title" db:"doc_title"`
	Summary     *string    `json:"summary" db:"doc_summary"`
	Source      *string    `json:"source" db:"doc_source"`
	PublishedAt *time.Time `json:"publishedAt" db:"doc_published_at"`
	Active      bool       `json:"active" db:"doc_active"`
	CreatedAt   time.Time  `json:"createdAt" db:"doc_created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"doc_updated_at"`
}

// Document Repository Params & Results
type CreateDocumentParams struct {
	Category    string
	Title       string
	Summary     *string
	Source      *string
	PublishedAt *time.Time
}

type CreateDocumentResult struct {
	dal.DbResult
	DocID int `json:"docId" db:"o_doc_id"`
}

type UpdateDocumentParams struct {
	DocID       int
	Category    string
	Title       string
	Summary     *string
	Source      *string
	PublishedAt *time.Time
}

type UpdateDocumentResult struct {
	dal.DbResult
}

type DeleteDocumentResult struct {
	dal.DbResult
}

// Document Repository & UseCase Interfaces
type DocumentRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]Document, error)
	GetByID(ctx context.Context, docID int) (*Document, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]Document, error)
	SearchByTitle(ctx context.Context, titlePattern string, limit int) ([]Document, error)
	Create(ctx context.Context, params CreateDocumentParams) (*CreateDocumentResult, error)
	Update(ctx context.Context, params UpdateDocumentParams) (*UpdateDocumentResult, error)
	Delete(ctx context.Context, docID int) (*DeleteDocumentResult, error)
}

type UploadPDFDocumentParams struct {
	Category     string
	Title        string
	Source       *string
	FileBase64   string
	ChunkSize    int
	ChunkOverlap int
}

type DocumentUseCase interface {
	GetAll(ctx context.Context, limit, offset int) Result[[]Document]
	GetByID(ctx context.Context, docID int) Result[*Document]
	GetByCategory(ctx context.Context, category string, limit, offset int) Result[[]Document]
	SearchByTitle(ctx context.Context, titlePattern string, limit int) Result[[]Document]
	Create(ctx context.Context, params CreateDocumentParams) Result[Data]
	Update(ctx context.Context, params UpdateDocumentParams) Result[Data]
	Delete(ctx context.Context, docID int) Result[Data]
	UploadPDF(ctx context.Context, params UploadPDFDocumentParams) Result[Data]
}
