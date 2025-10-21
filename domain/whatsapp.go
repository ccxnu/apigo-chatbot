package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
)

// WhatsAppSession represents a WhatsApp connection session
type WhatsAppSession struct {
	ID          int        `json:"id" db:"wss_id"`
	SessionName string     `json:"sessionName" db:"wss_session_name"`
	PhoneNumber string     `json:"phoneNumber" db:"wss_phone_number"`
	DeviceName  string     `json:"deviceName" db:"wss_device_name"`
	Platform    string     `json:"platform" db:"wss_platform"`
	QRCode      string     `json:"qrCode" db:"wss_qr_code"`
	Connected   bool       `json:"connected" db:"wss_connected"`
	LastSeen    *time.Time `json:"lastSeen" db:"wss_last_seen"`
	SessionData Data       `json:"sessionData" db:"wss_session_data"`
	Active      bool       `json:"active" db:"wss_active"`
	CreatedAt   time.Time  `json:"createdAt" db:"wss_created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"wss_updated_at"`
}

// Conversation represents a WhatsApp chat conversation
type Conversation struct {
	ID            int        `json:"id" db:"cnv_id"`
	UserID        *int       `json:"userId" db:"cnv_fk_user"`
	ChatID        string     `json:"chatId" db:"cnv_chat_id"`
	PhoneNumber   string     `json:"phoneNumber" db:"cnv_phone_number"`
	ContactName   string     `json:"contactName" db:"cnv_contact_name"`
	IsGroup       bool       `json:"isGroup" db:"cnv_is_group"`
	GroupName     string     `json:"groupName" db:"cnv_group_name"`
	LastMessageAt *time.Time `json:"lastMessageAt" db:"cnv_last_message_at"`
	MessageCount  int        `json:"messageCount" db:"cnv_message_count"`
	Active        bool       `json:"active" db:"cnv_active"`
	CreatedAt     time.Time  `json:"createdAt" db:"cnv_created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"cnv_updated_at"`
}

// ConversationMessage represents a single WhatsApp message
type ConversationMessage struct {
	ID             int       `json:"id" db:"cvm_id"`
	ConversationID int       `json:"conversationId" db:"cvm_fk_conversation"`
	MessageID      string    `json:"messageId" db:"cvm_message_id"`
	FromMe         bool      `json:"fromMe" db:"cvm_from_me"`
	SenderName     string    `json:"senderName" db:"cvm_sender_name"`
	MessageType    string    `json:"messageType" db:"cvm_message_type"`
	Body           string    `json:"body" db:"cvm_body"`
	MediaURL       string    `json:"mediaUrl" db:"cvm_media_url"`
	QuotedMessage  string    `json:"quotedMessage" db:"cvm_quoted_message"`
	Timestamp      int64     `json:"timestamp" db:"cvm_timestamp"`
	IsForwarded    bool      `json:"isForwarded" db:"cvm_is_forwarded"`
	Metadata       Data      `json:"metadata" db:"cvm_metadata"`
	CreatedAt      time.Time `json:"createdAt" db:"cvm_created_at"`
}

// IncomingMessage represents a WhatsApp message received by the bot
type IncomingMessage struct {
	MessageID     string
	ChatID        string
	From          string
	FromMe        bool
	SenderName    string
	Body          string
	MessageType   string
	Timestamp     int64
	IsGroup       bool
	GroupName     string
	QuotedMessage string
	MediaURL      string
	IsForwarded   bool
}

// WhatsApp Repository Params & Results

type CreateConversationParams struct {
	ChatID      string
	PhoneNumber string
	ContactName string
	IsGroup     bool
	GroupName   string
}

type CreateConversationResult struct {
	dal.DbResult
	ConversationID int `json:"conversationId" db:"o_cnv_id"`
}

type CreateMessageParams struct {
	ConversationID int
	MessageID      string
	FromMe         bool
	SenderName     string
	MessageType    string
	Body           string
	MediaURL       string
	QuotedMessage  string
	Timestamp      int64
	IsForwarded    bool
}

type CreateMessageResult struct {
	dal.DbResult
	MessageID int `json:"messageId" db:"o_cvm_id"`
}

type UpdateSessionStatusParams struct {
	SessionName string
	PhoneNumber string
	DeviceName  string
	Platform    string
	Connected   bool
}

type UpdateSessionStatusResult struct {
	dal.DbResult
}

type LinkUserToConversationParams struct {
	ChatID         string
	IdentityNumber string
}

type LinkUserToConversationResult struct {
	dal.DbResult
}

// WhatsApp Repository & UseCase Interfaces

type WhatsAppSessionRepository interface {
	GetBySessionName(ctx context.Context, sessionName string) (*WhatsAppSession, error)
	UpdateStatus(ctx context.Context, params UpdateSessionStatusParams) (*UpdateSessionStatusResult, error)
	UpdateQRCode(ctx context.Context, sessionName, qrCode string) error
	GetActiveSession(ctx context.Context) (*WhatsAppSession, error)
}

type ConversationRepository interface {
	GetByChatID(ctx context.Context, chatID string) (*Conversation, error)
	Create(ctx context.Context, params CreateConversationParams) (*CreateConversationResult, error)
	LinkUserToConversation(ctx context.Context, params LinkUserToConversationParams) (*LinkUserToConversationResult, error)
	GetConversationHistory(ctx context.Context, chatID string, limit int) ([]ConversationMessage, error)
	CreateMessage(ctx context.Context, params CreateMessageParams) (*CreateMessageResult, error)
	GetUserConversations(ctx context.Context, userID int, limit int) ([]Conversation, error)
}

type WhatsAppSessionUseCase interface {
	GetSessionStatus(ctx context.Context, sessionName string) Result[*WhatsAppSession]
	GetQRCode(ctx context.Context, sessionName string) Result[Data]
	UpdateConnectionStatus(ctx context.Context, params UpdateSessionStatusParams) Result[Data]
}

type ConversationUseCase interface {
	GetOrCreateConversation(ctx context.Context, params CreateConversationParams) Result[*Conversation]
	SaveMessage(ctx context.Context, params CreateMessageParams) Result[Data]
	GetConversationHistory(ctx context.Context, chatID string, limit int) Result[[]ConversationMessage]
	LinkUserAfterValidation(ctx context.Context, chatID, identityNumber string) Result[Data]
}
