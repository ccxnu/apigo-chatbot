package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
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
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) GetByID(c context.Context, docID int) d.Result[*d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	doc, err := u.docRepo.GetByID(ctx, docID)
	if err != nil {
		return d.Error[*d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	if doc == nil {
		return d.Error[*d.Document](u.paramCache, "ERR_DOCUMENT_NOT_FOUND")
	}

	return d.Success(doc)
}

func (u *documentUseCase) GetByCategory(c context.Context, category string, limit, offset int) d.Result[[]d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) SearchByTitle(c context.Context, titlePattern string, limit int) d.Result[[]d.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.SearchByTitle(ctx, titlePattern, limit)
	if err != nil {
		return d.Error[[]d.Document](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(docs)
}

func (u *documentUseCase) Create(c context.Context, params d.CreateDocumentParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Create(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{"docId": result.DocID})
}

func (u *documentUseCase) Update(c context.Context, params d.UpdateDocumentParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Update(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *documentUseCase) Delete(c context.Context, docID int) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Delete(ctx, docID)
	if err != nil || result == nil {
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}
