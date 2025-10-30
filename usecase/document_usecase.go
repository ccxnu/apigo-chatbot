package usecase

import (
	"context"
	"fmt"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
	"api-chatbot/internal/pdfprocessor"
	"api-chatbot/internal/textchunker"
)

type documentUseCase struct {
	docRepo        d.DocumentRepository
	chunkUseCase   d.ChunkUseCase
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewDocumentUseCase(
	docRepo d.DocumentRepository,
	chunkUseCase d.ChunkUseCase,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.DocumentUseCase {
	return &documentUseCase{
		docRepo:        docRepo,
		chunkUseCase:   chunkUseCase,
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

func (u *documentUseCase) UploadPDF(c context.Context, params d.UploadPDFDocumentParams) d.Result[d.Data] {
	// Use longer timeout for PDF processing
	ctx, cancel := context.WithTimeout(c, 2*time.Minute)
	defer cancel()

	logger.LogInfo(ctx, "Starting PDF upload",
		"operation", "UploadPDF",
		"title", params.Title,
		"category", params.Category,
		"chunkSize", params.ChunkSize,
		"chunkOverlap", params.ChunkOverlap,
	)

	// Step 1: Extract text from PDF
	text, err := pdfprocessor.ExtractTextFromBase64PDF(params.FileBase64)
	if err != nil {
		logger.LogError(ctx, "Failed to extract text from PDF", err,
			"operation", "UploadPDF",
			"title", params.Title,
		)
		return d.Error[d.Data](u.paramCache, "ERR_PDF_PROCESSING")
	}

	// Generate summary from first 500 characters
	summary := text
	if len(text) > 500 {
		summary = text[:500] + "..."
	}

	logger.LogInfo(ctx, "Text extracted from PDF",
		"operation", "UploadPDF",
		"textLength", len(text),
		"title", params.Title,
	)

	// Step 2: Create the document
	docParams := d.CreateDocumentParams{
		Category: params.Category,
		Title:    params.Title,
		Summary:  &summary,
		Source:   params.Source,
	}

	docResult, err := u.docRepo.Create(ctx, docParams)
	if err != nil || docResult == nil {
		logger.LogError(ctx, "Failed to create document in database", err,
			"operation", "UploadPDF",
			"title", params.Title,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !docResult.Success {
		logger.LogWarn(ctx, "Document creation failed with business logic error",
			"operation", "UploadPDF",
			"code", docResult.Code,
			"title", params.Title,
		)
		return d.Error[d.Data](u.paramCache, docResult.Code)
	}

	docID := docResult.DocID
	logger.LogInfo(ctx, "Document created successfully",
		"operation", "UploadPDF",
		"docID", docID,
		"title", params.Title,
	)

	// Step 3: Split text into chunks
	chunks := textchunker.ChunkText(text, params.ChunkSize, params.ChunkOverlap)

	logger.LogInfo(ctx, "Text split into chunks",
		"operation", "UploadPDF",
		"docID", docID,
		"chunksCount", len(chunks),
	)

	if len(chunks) == 0 {
		logger.LogWarn(ctx, "No chunks created from PDF text",
			"operation", "UploadPDF",
			"docID", docID,
			"textLength", len(text),
		)
		return d.Success(d.Data{
			"docId":        docID,
			"chunksCreated": 0,
			"message":      "Document created but no chunks generated (text might be empty)",
		})
	}

	// Step 4: Create chunks using ChunkUseCase
	chunkResult := u.chunkUseCase.BulkCreate(ctx, docID, chunks)
	if !chunkResult.Success {
		logger.LogError(ctx, "Failed to create chunks",
			fmt.Errorf("chunk creation failed: %s", chunkResult.Code),
			"operation", "UploadPDF",
			"docID", docID,
			"chunksCount", len(chunks),
		)
		return d.Error[d.Data](u.paramCache, "ERR_CHUNK_CREATION")
	}

	chunksCreated, _ := chunkResult.Data["chunksCreated"].(int)

	logger.LogInfo(ctx, "PDF upload completed successfully",
		"operation", "UploadPDF",
		"docID", docID,
		"chunksCreated", chunksCreated,
	)

	return d.Success(d.Data{
		"docId":         docID,
		"chunksCreated": chunksCreated,
		"message":       "PDF uploaded and processed successfully",
	})
}
