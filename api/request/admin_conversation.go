package request

import "api-chatbot/domain"

// GetConversationsRequest request for getting conversations
type GetConversationsRequest struct {
	domain.Base
	Filter string `json:"filter" validate:"omitempty,oneof=all unread blocked active" doc:"Filter conversations: all, unread, blocked, active"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100" doc:"Number of conversations to return (default: 50)"`
	Offset int    `json:"offset" validate:"omitempty,min=0" doc:"Offset for pagination (default: 0)"`
}

// GetConversationMessagesRequest request for getting conversation messages
type GetConversationMessagesRequest struct {
	domain.Base
	ConversationID int `json:"conversationId" validate:"required,min=1" doc:"Conversation ID"`
	Limit          int `json:"limit" validate:"omitempty,min=1,max=200" doc:"Number of messages to return (default: 100)"`
}

// SendAdminMessageRequest request for admin sending a message
type SendAdminMessageRequest struct {
	domain.Base
	ConversationID int    `json:"conversationId" validate:"required,min=1" doc:"Conversation ID"`
	Message        string `json:"message" validate:"required,min=1,max=4096" doc:"Message text to send"`
}

// MarkMessagesReadRequest request for marking messages as read
type MarkMessagesReadRequest struct {
	domain.Base
	ConversationID int `json:"conversationId" validate:"required,min=1" doc:"Conversation ID"`
}

// BlockUserRequest request for blocking a user
type BlockUserRequest struct {
	domain.Base
	UserID  int     `json:"userId" validate:"required,min=1" doc:"User ID to block/unblock"`
	Blocked bool    `json:"blocked" doc:"True to block, false to unblock"`
	Reason  *string `json:"reason" validate:"omitempty,max=500" doc:"Reason for blocking (optional)"`
}

// DeleteConversationRequest request for deleting a conversation
type DeleteConversationRequest struct {
	domain.Base
	ConversationID int `json:"conversationId" validate:"required,min=1" doc:"Conversation ID to delete"`
}

// SetTemporaryRequest request for setting conversation as temporary
type SetTemporaryRequest struct {
	domain.Base
	ConversationID   int  `json:"conversationId" validate:"required,min=1" doc:"Conversation ID"`
	Temporary        bool `json:"temporary" doc:"True to enable temporary, false to disable"`
	HoursUntilExpiry int  `json:"hoursUntilExpiry" validate:"omitempty,min=1,max=720" doc:"Hours until conversation expires (default: 24, max: 720/30 days)"`
}
