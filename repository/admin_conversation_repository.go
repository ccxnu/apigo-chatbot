package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Functions
	fnGetAllConversationsForAdmin = "fn_get_all_conversations_for_admin"
	fnGetConversationMessages     = "fn_get_conversation_messages"

	// Procedures
	spBlockUser                = "sp_block_user"
	spDeleteConversation       = "sp_delete_conversation"
	spSendAdminMessage         = "sp_send_admin_message"
	spMarkMessagesAsRead       = "sp_mark_messages_as_read"
	spSetConversationTemporary = "sp_set_conversation_temporary"
)

type adminConversationRepository struct {
	dal *dal.DAL
}

func NewAdminConversationRepository(dal *dal.DAL) d.AdminConversationRepository {
	return &adminConversationRepository{
		dal: dal,
	}
}

// GetAllConversations retrieves paginated conversations for admin panel
func (r *adminConversationRepository) GetAllConversations(ctx context.Context, filter string, limit, offset int) ([]d.AdminConversationListItem, error) {
	conversations, err := dal.QueryRows[d.AdminConversationListItem](
		r.dal,
		ctx,
		fnGetAllConversationsForAdmin,
		limit,
		offset,
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations via %s: %w", fnGetAllConversationsForAdmin, err)
	}

	return conversations, nil
}

// GetConversationMessages retrieves messages for a conversation
func (r *adminConversationRepository) GetConversationMessages(ctx context.Context, conversationID int, limit int) ([]d.AdminConversationMessage, error) {
	messages, err := dal.QueryRows[d.AdminConversationMessage](
		r.dal,
		ctx,
		fnGetConversationMessages,
		conversationID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages via %s: %w", fnGetConversationMessages, err)
	}

	return messages, nil
}

// BlockUser blocks or unblocks a user
func (r *adminConversationRepository) BlockUser(ctx context.Context, params d.BlockUserParams) (*d.BlockUserResult, error) {
	var reason interface{} = nil
	if params.Reason != nil {
		reason = *params.Reason
	}

	result, err := dal.ExecProc[d.BlockUserResult](
		r.dal,
		ctx,
		spBlockUser,
		params.UserID,
		params.Blocked,
		params.AdminID,
		reason,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spBlockUser, err)
	}

	return result, nil
}

// DeleteConversation soft deletes a conversation
func (r *adminConversationRepository) DeleteConversation(ctx context.Context, conversationID int) (*d.DeleteConversationResult, error) {
	result, err := dal.ExecProc[d.DeleteConversationResult](
		r.dal,
		ctx,
		spDeleteConversation,
		conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spDeleteConversation, err)
	}

	return result, nil
}

// SendAdminMessage stores an admin message
func (r *adminConversationRepository) SendAdminMessage(ctx context.Context, params d.SendAdminMessageParams) (*d.SendAdminMessageResult, error) {
	result, err := dal.ExecProc[d.SendAdminMessageResult](
		r.dal,
		ctx,
		spSendAdminMessage,
		params.ConversationID,
		params.AdminID,
		params.MessageID,
		params.Body,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spSendAdminMessage, err)
	}

	return result, nil
}

// MarkMessagesAsRead marks all messages in conversation as read
func (r *adminConversationRepository) MarkMessagesAsRead(ctx context.Context, conversationID int) (*d.MarkReadResult, error) {
	result, err := dal.ExecProc[d.MarkReadResult](
		r.dal,
		ctx,
		spMarkMessagesAsRead,
		conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spMarkMessagesAsRead, err)
	}

	return result, nil
}

// SetConversationTemporary enables/disables temporary conversation
func (r *adminConversationRepository) SetConversationTemporary(ctx context.Context, params d.SetTemporaryParams) (*d.SetTemporaryResult, error) {
	result, err := dal.ExecProc[d.SetTemporaryResult](
		r.dal,
		ctx,
		spSetConversationTemporary,
		params.ConversationID,
		params.Temporary,
		params.HoursUntilExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spSetConversationTemporary, err)
	}

	return result, nil
}

// GetConversationByChatID retrieves conversation by WhatsApp chat ID
func (r *adminConversationRepository) GetConversationByChatID(ctx context.Context, chatID string) (*d.Conversation, error) {
	conversations, err := dal.QueryRows[d.Conversation](
		r.dal,
		ctx,
		fnGetConversationByChatID,
		chatID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation via %s: %w", fnGetConversationByChatID, err)
	}

	if len(conversations) == 0 {
		return nil, nil
	}

	return &conversations[0], nil
}
