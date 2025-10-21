package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type conversationUseCase struct {
	convRepo       d.ConversationRepository
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewConversationUseCase(
	convRepo d.ConversationRepository,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.ConversationUseCase {
	return &conversationUseCase{
		convRepo:       convRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

func (u *conversationUseCase) GetOrCreateConversation(c context.Context, params d.CreateConversationParams) d.Result[*d.Conversation] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Try to get existing conversation
	conversation, err := u.convRepo.GetByChatID(ctx, params.ChatID)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch conversation by chat ID from database", err,
			"operation", "GetOrCreateConversation",
			"chatID", params.ChatID,
		)
		return d.Error[*d.Conversation](u.paramCache, "ERR_INTERNAL_DB")
	}

	// If found, return it
	if conversation != nil {
		return d.Success(conversation)
	}

	// Create new conversation
	result, err := u.convRepo.Create(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to create conversation in database", err,
			"operation", "GetOrCreateConversation",
			"chatID", params.ChatID,
		)
		return d.Error[*d.Conversation](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Conversation creation failed with business logic error",
			"operation", "GetOrCreateConversation",
			"code", result.Code,
			"chatID", params.ChatID,
		)
		return d.Error[*d.Conversation](u.paramCache, result.Code)
	}

	// Get the created conversation
	conversation, err = u.convRepo.GetByChatID(ctx, params.ChatID)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch newly created conversation from database", err,
			"operation", "GetOrCreateConversation",
			"chatID", params.ChatID,
		)
		return d.Error[*d.Conversation](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(conversation)
}

func (u *conversationUseCase) SaveMessage(c context.Context, params d.CreateMessageParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.convRepo.CreateMessage(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to save message in database", err,
			"operation", "SaveMessage",
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Message save failed with business logic error",
			"operation", "SaveMessage",
			"code", result.Code,
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{"messageId": result.MessageID})
}

func (u *conversationUseCase) GetConversationHistory(c context.Context, chatID string, limit int) d.Result[[]d.ConversationMessage] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	messages, err := u.convRepo.GetConversationHistory(ctx, chatID, limit)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch conversation history from database", err,
			"operation", "GetConversationHistory",
			"chatID", chatID,
			"limit", limit,
		)
		return d.Error[[]d.ConversationMessage](u.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(messages)
}

func (u *conversationUseCase) LinkUserAfterValidation(c context.Context, chatID, identityNumber string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	params := d.LinkUserToConversationParams{
		ChatID:         chatID,
		IdentityNumber: identityNumber,
	}

	result, err := u.convRepo.LinkUserToConversation(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to link user to conversation in database", err,
			"operation", "LinkUserAfterValidation",
			"chatID", chatID,
			"identityNumber", identityNumber,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "User linking failed with business logic error",
			"operation", "LinkUserAfterValidation",
			"code", result.Code,
			"chatID", chatID,
			"identityNumber", identityNumber,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}
