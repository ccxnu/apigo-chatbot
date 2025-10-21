package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	"api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetConversationByChatID = "fn_get_conversation_by_chat_id"
	fnGetConversationHistory  = "fn_get_conversation_history"
	// Stored Procedures (Writes)
	spCreateConversation        = "sp_create_conversation"
	spLinkUserToConversation    = "sp_link_user_to_conversation"
	spCreateConversationMessage = "sp_create_conversation_message"
)

type conversationRepository struct {
	dal *dal.DAL
}

func NewConversationRepository(dal *dal.DAL) domain.ConversationRepository {
	return &conversationRepository{
		dal: dal,
	}
}

// GetByChatID retrieves a conversation by WhatsApp chat ID
func (r *conversationRepository) GetByChatID(ctx context.Context, chatID string) (*domain.Conversation, error) {
	conversations, err := dal.QueryRows[domain.Conversation](r.dal, ctx, fnGetConversationByChatID, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation via %s: %w", fnGetConversationByChatID, err)
	}

	if len(conversations) == 0 {
		return nil, nil
	}

	return &conversations[0], nil
}

// Create creates a new conversation or returns existing one
func (r *conversationRepository) Create(ctx context.Context, params domain.CreateConversationParams) (*domain.CreateConversationResult, error) {
	result, err := dal.ExecProc[domain.CreateConversationResult](
		r.dal,
		ctx,
		spCreateConversation,
		params.ChatID,
		params.PhoneNumber,
		params.ContactName,
		params.IsGroup,
		params.GroupName,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateConversation, err)
	}

	return result, nil
}

// LinkUserToConversation links a validated user to a conversation
func (r *conversationRepository) LinkUserToConversation(ctx context.Context, params domain.LinkUserToConversationParams) (*domain.LinkUserToConversationResult, error) {
	result, err := dal.ExecProc[domain.LinkUserToConversationResult](
		r.dal,
		ctx,
		spLinkUserToConversation,
		params.ChatID,
		params.IdentityNumber,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spLinkUserToConversation, err)
	}

	return result, nil
}

// GetConversationHistory retrieves message history for a conversation
func (r *conversationRepository) GetConversationHistory(ctx context.Context, chatID string, limit int) ([]domain.ConversationMessage, error) {
	messages, err := dal.QueryRows[domain.ConversationMessage](r.dal, ctx, fnGetConversationHistory, chatID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history via %s: %w", fnGetConversationHistory, err)
	}

	return messages, nil
}

// CreateMessage stores a new message in a conversation
func (r *conversationRepository) CreateMessage(ctx context.Context, params domain.CreateMessageParams) (*domain.CreateMessageResult, error) {
	result, err := dal.ExecProc[domain.CreateMessageResult](
		r.dal,
		ctx,
		spCreateConversationMessage,
		params.ConversationID,
		params.MessageID,
		params.FromMe,
		params.SenderName,
		params.MessageType,
		params.Body,
		params.MediaURL,
		params.QuotedMessage,
		params.Timestamp,
		params.IsForwarded,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreateConversationMessage, err)
	}

	return result, nil
}

// GetUserConversations retrieves all conversations for a user
func (r *conversationRepository) GetUserConversations(ctx context.Context, userID int, limit int) ([]domain.Conversation, error) {
	// TODO: Implement this function if needed
	// For now, returning empty slice
	return []domain.Conversation{}, nil
}
