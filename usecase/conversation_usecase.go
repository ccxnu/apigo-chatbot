package usecase

import (
	"context"
	"time"

	"api-chatbot/domain"
)

type conversationUseCase struct {
	convRepo       domain.ConversationRepository
	paramCache     domain.ParameterCache
	contextTimeout time.Duration
}

func NewConversationUseCase(
	convRepo domain.ConversationRepository,
	paramCache domain.ParameterCache,
	timeout time.Duration,
) domain.ConversationUseCase {
	return &conversationUseCase{
		convRepo:       convRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// getErrorMessage retrieves error message from parameter cache
func (u *conversationUseCase) getErrorMessage(errorCode string) string {
	if param, exists := u.paramCache.Get(errorCode); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if message, ok := data["message"].(string); ok {
				return message
			}
		}
	}
	return "Ha ocurrido un error"
}

func (u *conversationUseCase) GetOrCreateConversation(c context.Context, params domain.CreateConversationParams) domain.Result[*domain.Conversation] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Try to get existing conversation
	conversation, err := u.convRepo.GetByChatID(ctx, params.ChatID)
	if err != nil {
		return domain.Result[*domain.Conversation]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	// If found, return it
	if conversation != nil {
		return domain.Result[*domain.Conversation]{
			Success: true,
			Code:    "OK",
			Info:    u.getErrorMessage("OK"),
			Data:    conversation,
		}
	}

	// Create new conversation
	result, err := u.convRepo.Create(ctx, params)
	if err != nil || result == nil {
		return domain.Result[*domain.Conversation]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[*domain.Conversation]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	// Get the created conversation
	conversation, err = u.convRepo.GetByChatID(ctx, params.ChatID)
	if err != nil {
		return domain.Result[*domain.Conversation]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	return domain.Result[*domain.Conversation]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    conversation,
	}
}

func (u *conversationUseCase) SaveMessage(c context.Context, params domain.CreateMessageParams) domain.Result[domain.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.convRepo.CreateMessage(ctx, params)
	if err != nil || result == nil {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[domain.Data]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data: domain.Data{
			"messageId": result.MessageID,
		},
	}
}

func (u *conversationUseCase) GetConversationHistory(c context.Context, chatID string, limit int) domain.Result[[]domain.ConversationMessage] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	messages, err := u.convRepo.GetConversationHistory(ctx, chatID, limit)
	if err != nil {
		return domain.Result[[]domain.ConversationMessage]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    []domain.ConversationMessage{},
		}
	}

	return domain.Result[[]domain.ConversationMessage]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    messages,
	}
}

func (u *conversationUseCase) LinkUserAfterValidation(c context.Context, chatID, identityNumber string) domain.Result[domain.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	params := domain.LinkUserToConversationParams{
		ChatID:         chatID,
		IdentityNumber: identityNumber,
	}

	result, err := u.convRepo.LinkUserToConversation(ctx, params)
	if err != nil || result == nil {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[domain.Data]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}
