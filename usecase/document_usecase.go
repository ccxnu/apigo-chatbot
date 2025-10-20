package usecase

import (
	"context"
	"time"

	"api-chatbot/domain"
)

type documentUseCase struct {
	docRepo        domain.DocumentRepository
	paramCache     domain.ParameterCache
	contextTimeout time.Duration
}

func NewDocumentUseCase(
	docRepo domain.DocumentRepository,
	paramCache domain.ParameterCache,
	timeout time.Duration,
) domain.DocumentUseCase {
	return &documentUseCase{
		docRepo:        docRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// getErrorMessage retrieves error message from parameter cache
func (u *documentUseCase) getErrorMessage(errorCode string) string {
	if param, exists := u.paramCache.Get(errorCode); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if message, ok := data["message"].(string); ok {
				return message
			}
		}
	}
	return "Ha ocurrido un error"
}

func (u *documentUseCase) GetAll(c context.Context, limit, offset int) domain.Result[[]domain.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return domain.Result[[]domain.Document]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    []domain.Document{},
		}
	}

	return domain.Result[[]domain.Document]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    docs,
	}
}

func (u *documentUseCase) GetByID(c context.Context, docID int) domain.Result[*domain.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	doc, err := u.docRepo.GetByID(ctx, docID)
	if err != nil {
		return domain.Result[*domain.Document]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if doc == nil {
		return domain.Result[*domain.Document]{
			Success: false,
			Code:    "ERR_DOCUMENT_NOT_FOUND",
			Info:    u.getErrorMessage("ERR_DOCUMENT_NOT_FOUND"),
			Data:    nil,
		}
	}

	return domain.Result[*domain.Document]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    doc,
	}
}

func (u *documentUseCase) GetByCategory(c context.Context, category string, limit, offset int) domain.Result[[]domain.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return domain.Result[[]domain.Document]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    []domain.Document{},
		}
	}

	return domain.Result[[]domain.Document]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    docs,
	}
}

func (u *documentUseCase) SearchByTitle(c context.Context, titlePattern string, limit int) domain.Result[[]domain.Document] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	docs, err := u.docRepo.SearchByTitle(ctx, titlePattern, limit)
	if err != nil {
		return domain.Result[[]domain.Document]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    []domain.Document{},
		}
	}

	return domain.Result[[]domain.Document]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    docs,
	}
}

func (u *documentUseCase) Create(c context.Context, params domain.CreateDocumentParams) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Create(ctx, params)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data: map[string]any{
			"docId": result.DocID,
		},
	}
}

func (u *documentUseCase) Update(c context.Context, params domain.UpdateDocumentParams) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Update(ctx, params)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}

func (u *documentUseCase) Delete(c context.Context, docID int) domain.Result[map[string]any] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.docRepo.Delete(ctx, docID)
	if err != nil || result == nil {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[map[string]any]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[map[string]any]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}
