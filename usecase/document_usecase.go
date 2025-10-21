package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type documentUseCase struct {
	docRepo        d.DocumentRepository
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewDocumentUseCase(
	docRepo d.DocumentRepository,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.DocumentUseCase {
	return &documentUseCase{
		docRepo:        docRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

func (u *documentUseCase) GetAll(c context.Context, limit, offset int) d.Result[[]d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch all documents from database", err,
			"operation", "GetAll",
			"limit", limit,
			"offset", offset,
		)
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) GetByID(c context.Context, docID int) d.Result[*d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	doc, err := u.docRepo.GetByID(ctx, docID)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch document by ID from database", err,
			"operation", "GetByID",
			"docID", docID,
		)
		return d.Error[*d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	if doc == nil {
		logger.LogWarn(ctx, "Document not found",
			"operation", "GetByID",
			"docID", docID,
		)
		return d.Error[*d.Document](u.paramCache, "ERR_DOCUMENT_NOT_FOUND")
	}

	return d.Success(doc)
}

func (u *documentUseCase) GetByCategory(c context.Context, category string, limit, offset int) d.Result[[]d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch documents by category from database", err,
			"operation", "GetByCategory",
			"category", category,
			"limit", limit,
			"offset", offset,
		)
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) SearchByTitle(c context.Context, titlePattern string, limit int) d.Result[[]d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.SearchByTitle(ctx, titlePattern, limit)
	if err != nil {
		logger.LogError(ctx, "Failed to search documents by title from database", err,
			"operation", "SearchByTitle",
			"titlePattern", titlePattern,
			"limit", limit,
		)
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) Create(c context.Context, params d.CreateDocumentParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Create(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to create document in database", err,
			"operation", "Create",
			"title", params.Title,
			"category", params.Category,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Document creation failed with business logic error",
			"operation", "Create",
			"code", result.Code,
			"title", params.Title,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{"docId": result.DocID})
}

func (u *documentUseCase) Update(c context.Context, params d.UpdateDocumentParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Update(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to update document in database", err, params)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Document update failed with business logic error",
			"operation", "Update",
			"code", result.Code,
			"docID", params.DocID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *documentUseCase) Delete(c context.Context, docID int) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Delete(ctx, docID)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to delete document from database", err,
			"operation", "Delete",
			"docID", docID,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Document deletion failed with business logic error",
			"operation", "Delete",
			"code", result.Code,
			"docID", docID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}
