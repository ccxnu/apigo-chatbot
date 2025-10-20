package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetAllDocuments        = "fn_get_all_documents"
	fnGetDocumentByID        = "fn_get_document_by_id"
	fnGetDocumentsByCategory = "fn_get_documents_by_category"
	fnSearchDocumentsByTitle = "fn_search_documents_by_title"
	// Stored Procedures (Writes)
	spCreateDocument = "sp_create_document"
	spUpdateDocument = "sp_update_document"
	spDeleteDocument = "sp_delete_document"
)

type documentRepository struct {
	dal *dal.DAL
}

func NewDocumentRepository(dal *dal.DAL) domain.DocumentRepository {
	return &documentRepository{
		dal: dal,
	}
}

// GetAll retrieves all active documents with pagination
func (r *documentRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Document, error) {
	docs, err := dal.QueryRows[domain.Document](r.dal, ctx, fnGetAllDocuments, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all documents via %s: %w", fnGetAllDocuments, err)
	}
	return docs, nil
}

// GetByID retrieves a single document by ID
func (r *documentRepository) GetByID(ctx context.Context, docID int) (*domain.Document, error) {
	docs, err := dal.QueryRows[domain.Document](r.dal, ctx, fnGetDocumentByID, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document by id via %s: %w", fnGetDocumentByID, err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	return &docs[0], nil
}

// GetByCategory retrieves documents filtered by category
func (r *documentRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]domain.Document, error) {
	docs, err := dal.QueryRows[domain.Document](r.dal, ctx, fnGetDocumentsByCategory, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by category via %s: %w", fnGetDocumentsByCategory, err)
	}
	return docs, nil
}

// SearchByTitle searches documents by title pattern
func (r *documentRepository) SearchByTitle(ctx context.Context, titlePattern string, limit int) ([]domain.Document, error) {
	docs, err := dal.QueryRows[domain.Document](r.dal, ctx, fnSearchDocumentsByTitle, titlePattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents by title via %s: %w", fnSearchDocumentsByTitle, err)
	}
	return docs, nil
}

// Create creates a new document
func (r *documentRepository) Create(ctx context.Context, params domain.CreateDocumentParams) (*domain.CreateDocumentResult, error) {
	result, err := dal.ExecProc[domain.CreateDocumentResult](
		r.dal,
		ctx,
		spCreateDocument,
		params.Category,
		params.Title,
		params.Summary,
		params.Source,
		params.PublishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateDocument, err)
	}

	return result, nil
}

// Update updates an existing document
func (r *documentRepository) Update(ctx context.Context, params domain.UpdateDocumentParams) (*domain.UpdateDocumentResult, error) {
	result, err := dal.ExecProc[domain.UpdateDocumentResult](
		r.dal,
		ctx,
		spUpdateDocument,
		params.DocID,
		params.Category,
		params.Title,
		params.Summary,
		params.Source,
		params.PublishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateDocument, err)
	}

	return result, nil
}

// Delete soft deletes a document
func (r *documentRepository) Delete(ctx context.Context, docID int) (*domain.DeleteDocumentResult, error) {
	result, err := dal.ExecProc[domain.DeleteDocumentResult](
		r.dal,
		ctx,
		spDeleteDocument,
		docID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spDeleteDocument, err)
	}

	return result, nil
}
