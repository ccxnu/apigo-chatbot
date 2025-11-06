package domain

import (
	"context"
	"time"
)

// AdminConversationListItem represents a conversation in the admin panel list
type AdminConversationListItem struct {
	// Conversation data
	ID              int        `json:"id" db:"cnv_id"`
	ChatID          string     `json:"chatId" db:"cnv_chat_id"`
	PhoneNumber     string     `json:"phoneNumber" db:"cnv_phone_number"`
	ContactName     *string    `json:"contactName,omitempty" db:"cnv_contact_name"`
	IsGroup         bool       `json:"isGroup" db:"cnv_is_group"`
	GroupName       *string    `json:"groupName,omitempty" db:"cnv_group_name"`
	LastMessageAt   *time.Time `json:"lastMessageAt,omitempty" db:"cnv_last_message_at"`
	MessageCount    int        `json:"messageCount" db:"cnv_message_count"`
	UnreadCount     int        `json:"unreadCount" db:"cnv_unread_count"`
	Blocked         bool       `json:"blocked" db:"cnv_blocked"`
	AdminIntervened bool       `json:"adminIntervened" db:"cnv_admin_intervened"`
	Temporary       bool       `json:"temporary" db:"cnv_temporary"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty" db:"cnv_expires_at"`

	// User data
	UserID             *int    `json:"userId,omitempty" db:"usr_id"`
	UserName           *string `json:"userName,omitempty" db:"usr_name"`
	UserIdentityNumber *string `json:"userIdentityNumber,omitempty" db:"usr_identity_number"`
	UserRole           *string `json:"userRole,omitempty" db:"usr_rol"`
	UserBlocked        *bool   `json:"userBlocked,omitempty" db:"usr_blocked"`

	// Last message preview
	LastMessagePreview *string `json:"lastMessagePreview,omitempty" db:"last_message_preview"`
	LastMessageFromMe  *bool   `json:"lastMessageFromMe,omitempty" db:"last_message_from_me"`
}

// AdminConversationMessage represents a message in the admin panel
type AdminConversationMessage struct {
	ID            int       `json:"id" db:"cvm_id"`
	MessageID     string    `json:"messageId" db:"cvm_message_id"`
	FromMe        bool      `json:"fromMe" db:"cvm_from_me"`
	SenderName    *string   `json:"senderName,omitempty" db:"cvm_sender_name"`
	SenderType    string    `json:"senderType" db:"cvm_sender_type"` // user, admin, bot
	MessageType   string    `json:"messageType" db:"cvm_message_type"`
	Body          *string   `json:"body,omitempty" db:"cvm_body"`
	MediaURL      *string   `json:"mediaUrl,omitempty" db:"cvm_media_url"`
	QuotedMessage *string   `json:"quotedMessage,omitempty" db:"cvm_quoted_message"`
	Timestamp     int64     `json:"timestamp" db:"cvm_timestamp"`
	IsForwarded   bool      `json:"isForwarded" db:"cvm_is_forwarded"`
	Read          bool      `json:"read" db:"cvm_read"`
	CreatedAt     time.Time `json:"createdAt" db:"cvm_created_at"`
	AdminName     *string   `json:"adminName,omitempty" db:"admin_name"`
}

// BlockUserParams parameters for blocking a user
type BlockUserParams struct {
	UserID  int
	Blocked bool
	AdminID int
	Reason  *string
}

// BlockUserResult result of blocking a user
type BlockUserResult struct {
	Success bool   `json:"success" db:"success"`
	Code    string `json:"code" db:"code"`
}

// SendAdminMessageParams parameters for admin sending a message
type SendAdminMessageParams struct {
	ConversationID int
	AdminID        int
	MessageID      string
	Body           string
}

// SendAdminMessageResult result of sending admin message
type SendAdminMessageResult struct {
	Success   bool   `json:"success" db:"success"`
	Code      string `json:"code" db:"code"`
	MessageID *int   `json:"messageId,omitempty" db:"o_message_id"`
}

// SetTemporaryParams parameters for setting conversation as temporary
type SetTemporaryParams struct {
	ConversationID   int
	Temporary        bool
	HoursUntilExpiry int
}

// SetTemporaryResult result of setting temporary
type SetTemporaryResult struct {
	Success bool   `json:"success" db:"success"`
	Code    string `json:"code" db:"code"`
}

// MarkReadResult result of marking messages as read
type MarkReadResult struct {
	Success bool   `json:"success" db:"success"`
	Code    string `json:"code" db:"code"`
}

// DeleteConversationResult result of deleting conversation
type DeleteConversationResult struct {
	Success bool   `json:"success" db:"success"`
	Code    string `json:"code" db:"code"`
}

// AdminConversationRepository interface for admin conversation data access
type AdminConversationRepository interface {
	// GetAllConversations retrieves paginated conversations for admin panel
	GetAllConversations(ctx context.Context, filter string, limit, offset int) ([]AdminConversationListItem, error)

	// GetConversationMessages retrieves messages for a conversation
	GetConversationMessages(ctx context.Context, conversationID int, limit int) ([]AdminConversationMessage, error)

	// BlockUser blocks or unblocks a user
	BlockUser(ctx context.Context, params BlockUserParams) (*BlockUserResult, error)

	// DeleteConversation soft deletes a conversation
	DeleteConversation(ctx context.Context, conversationID int) (*DeleteConversationResult, error)

	// SendAdminMessage stores an admin message
	SendAdminMessage(ctx context.Context, params SendAdminMessageParams) (*SendAdminMessageResult, error)

	// MarkMessagesAsRead marks all messages in conversation as read
	MarkMessagesAsRead(ctx context.Context, conversationID int) (*MarkReadResult, error)

	// SetConversationTemporary enables/disables temporary conversation
	SetConversationTemporary(ctx context.Context, params SetTemporaryParams) (*SetTemporaryResult, error)

	// GetConversationByChatID retrieves conversation by WhatsApp chat ID
	GetConversationByChatID(ctx context.Context, chatID string) (*Conversation, error)
}

// AdminConversationUseCase interface for admin conversation business logic
type AdminConversationUseCase interface {
	// GetAllConversations retrieves paginated conversations
	GetAllConversations(ctx context.Context, filter string, limit, offset int) Result[[]AdminConversationListItem]

	// GetConversationMessages retrieves messages for a conversation
	GetConversationMessages(ctx context.Context, conversationID int, limit int) Result[[]AdminConversationMessage]

	// BlockUser blocks or unblocks a user
	BlockUser(ctx context.Context, params BlockUserParams) Result[Data]

	// DeleteConversation soft deletes a conversation
	DeleteConversation(ctx context.Context, conversationID int) Result[Data]

	// SendAdminMessage sends a message as admin (also sends via WhatsApp)
	SendAdminMessage(ctx context.Context, params SendAdminMessageParams) Result[Data]

	// MarkMessagesAsRead marks all messages as read
	MarkMessagesAsRead(ctx context.Context, conversationID int) Result[Data]

	// SetConversationTemporary enables/disables temporary conversation
	SetConversationTemporary(ctx context.Context, params SetTemporaryParams) Result[Data]
}
