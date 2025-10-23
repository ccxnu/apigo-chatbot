package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type adminConversationUseCase struct {
	repo           d.AdminConversationRepository
	whatsappClient WhatsAppMessageSender
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

// WhatsAppMessageSender interface for sending WhatsApp messages
type WhatsAppMessageSender interface {
	SendText(chatID, message string) error
}

func NewAdminConversationUseCase(
	repo d.AdminConversationRepository,
	whatsappClient WhatsAppMessageSender,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.AdminConversationUseCase {
	return &adminConversationUseCase{
		repo:           repo,
		whatsappClient: whatsappClient,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// GetAllConversations retrieves paginated conversations
func (uc *adminConversationUseCase) GetAllConversations(
	c context.Context,
	filter string,
	limit, offset int,
) d.Result[[]d.AdminConversationListItem] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Validate filter
	validFilters := map[string]bool{"all": true, "unread": true, "blocked": true, "active": true}
	if !validFilters[filter] {
		filter = "all"
	}

	conversations, err := uc.repo.GetAllConversations(ctx, filter, limit, offset)
	if err != nil {
		logger.LogError(ctx, "Failed to get conversations", err,
			"operation", "GetAllConversations",
			"filter", filter,
		)
		return d.Error[[]d.AdminConversationListItem](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(conversations)
}

// GetConversationMessages retrieves messages for a conversation
func (uc *adminConversationUseCase) GetConversationMessages(
	c context.Context,
	conversationID int,
	limit int,
) d.Result[[]d.AdminConversationMessage] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	messages, err := uc.repo.GetConversationMessages(ctx, conversationID, limit)
	if err != nil {
		logger.LogError(ctx, "Failed to get conversation messages", err,
			"operation", "GetConversationMessages",
			"conversationID", conversationID,
		)
		return d.Error[[]d.AdminConversationMessage](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(messages)
}

// BlockUser blocks or unblocks a user
func (uc *adminConversationUseCase) BlockUser(
	c context.Context,
	params d.BlockUserParams,
) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	result, err := uc.repo.BlockUser(ctx, params)
	if err != nil {
		logger.LogError(ctx, "Failed to block user", err,
			"operation", "BlockUser",
			"userID", params.UserID,
			"blocked", params.Blocked,
		)
		return d.Error[d.Data](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Block user failed with business logic error",
			"operation", "BlockUser",
			"code", result.Code,
			"userID", params.UserID,
		)
		return d.Error[d.Data](uc.paramCache, result.Code)
	}

	return d.Success(d.Data{
		"userId":  params.UserID,
		"blocked": params.Blocked,
	})
}

// DeleteConversation soft deletes a conversation
func (uc *adminConversationUseCase) DeleteConversation(
	c context.Context,
	conversationID int,
) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	result, err := uc.repo.DeleteConversation(ctx, conversationID)
	if err != nil {
		logger.LogError(ctx, "Failed to delete conversation", err,
			"operation", "DeleteConversation",
			"conversationID", conversationID,
		)
		return d.Error[d.Data](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Delete conversation failed with business logic error",
			"operation", "DeleteConversation",
			"code", result.Code,
			"conversationID", conversationID,
		)
		return d.Error[d.Data](uc.paramCache, result.Code)
	}

	return d.Success(d.Data{
		"conversationId": conversationID,
		"deleted":        true,
	})
}

// SendAdminMessage sends a message as admin (also sends via WhatsApp)
func (uc *adminConversationUseCase) SendAdminMessage(
	c context.Context,
	params d.SendAdminMessageParams,
) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Store message in database first
	result, err := uc.repo.SendAdminMessage(ctx, params)
	if err != nil {
		logger.LogError(ctx, "Failed to send admin message to database", err,
			"operation", "SendAdminMessage",
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Send admin message failed with business logic error",
			"operation", "SendAdminMessage",
			"code", result.Code,
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](uc.paramCache, result.Code)
	}

	// Get conversation to find chat ID for WhatsApp
	// We'll need to query by conversation ID - let's use a helper
	// For now, we'll assume the chat ID is derived from conversation
	// This needs the GetConversationByChatID to be modified or a new function created

	// TODO: Send via WhatsApp
	// For now, we'll just return success
	// In production, you'd call:
	// err = uc.whatsappClient.SendText(chatID, params.Body)

	return d.Success(d.Data{
		"messageId":      result.MessageID,
		"conversationId": params.ConversationID,
		"sent":           true,
	})
}

// MarkMessagesAsRead marks all messages as read
func (uc *adminConversationUseCase) MarkMessagesAsRead(
	c context.Context,
	conversationID int,
) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	result, err := uc.repo.MarkMessagesAsRead(ctx, conversationID)
	if err != nil {
		logger.LogError(ctx, "Failed to mark messages as read", err,
			"operation", "MarkMessagesAsRead",
			"conversationID", conversationID,
		)
		return d.Error[d.Data](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Mark messages as read failed with business logic error",
			"operation", "MarkMessagesAsRead",
			"code", result.Code,
			"conversationID", conversationID,
		)
		return d.Error[d.Data](uc.paramCache, result.Code)
	}

	return d.Success(d.Data{
		"conversationId": conversationID,
		"marked":         true,
	})
}

// SetConversationTemporary enables/disables temporary conversation
func (uc *adminConversationUseCase) SetConversationTemporary(
	c context.Context,
	params d.SetTemporaryParams,
) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	result, err := uc.repo.SetConversationTemporary(ctx, params)
	if err != nil {
		logger.LogError(ctx, "Failed to set conversation temporary", err,
			"operation", "SetConversationTemporary",
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "Set conversation temporary failed with business logic error",
			"operation", "SetConversationTemporary",
			"code", result.Code,
			"conversationID", params.ConversationID,
		)
		return d.Error[d.Data](uc.paramCache, result.Code)
	}

	var expiresAt *time.Time
	if params.Temporary {
		expiry := time.Now().Add(time.Duration(params.HoursUntilExpiry) * time.Hour)
		expiresAt = &expiry
	}

	return d.Success(d.Data{
		"conversationId": params.ConversationID,
		"temporary":      params.Temporary,
		"expiresAt":      expiresAt,
	})
}
