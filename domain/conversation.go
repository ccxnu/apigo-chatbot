package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
)

// Conversation represents a WhatsApp conversation
type Conversation struct {
	ID            int        `json:"id" db:"cnv_id"`
	UserID        *int       `json:"userId,omitempty" db:"cnv_fk_user"`
	ChatID        string     `json:"chatId" db:"cnv_chat_id"`
	PhoneNumber   string     `json:"phoneNumber" db:"cnv_phone_number"`
	ContactName   *string    `json:"contactName,omitempty" db:"cnv_contact_name"`
	IsGroup       bool       `json:"isGroup" db:"cnv_is_group"`
	GroupName     *string    `json:"groupName,omitempty" db:"cnv_group_name"`
	LastMessageAt *time.Time `json:"lastMessageAt,omitempty" db:"cnv_last_message_at"`
	MessageCount  int        `json:"messageCount" db:"cnv_message_count"`
	Active        bool       `json:"active" db:"cnv_active"`
	CreatedAt     time.Time  `json:"createdAt" db:"cnv_created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"cnv_updated_at"`
}

// ConversationMessage represents a message in a conversation
type ConversationMessage struct {
	ID               int       `json:"id" db:"cvm_id"`
	ConversationID   int       `json:"conversationId" db:"cvm_fk_conversation"`
	MessageID        string    `json:"messageId" db:"cvm_message_id"`
	FromMe           bool      `json:"fromMe" db:"cvm_from_me"`
	SenderName       *string   `json:"senderName,omitempty" db:"cvm_sender_name"`
	MessageType      string    `json:"messageType" db:"cvm_message_type"`
	Body             *string   `json:"body,omitempty" db:"cvm_body"`
	MediaURL         *string   `json:"mediaUrl,omitempty" db:"cvm_media_url"`
	QuotedMessage    *string   `json:"quotedMessage,omitempty" db:"cvm_quoted_message"`
	Timestamp        int64     `json:"timestamp" db:"cvm_timestamp"`
	IsForwarded      bool      `json:"isForwarded" db:"cvm_is_forwarded"`
	Metadata         Data      `json:"metadata,omitempty" db:"cvm_metadata"`
	QueueTimeMs      *int      `json:"queueTimeMs,omitempty" db:"cvm_queue_time_ms"`
	PromptTokens     *int      `json:"promptTokens,omitempty" db:"cvm_prompt_tokens"`
	PromptTimeMs     *int      `json:"promptTimeMs,omitempty" db:"cvm_prompt_time_ms"`
	CompletionTokens *int      `json:"completionTokens,omitempty" db:"cvm_completion_tokens"`
	CompletionTimeMs *int      `json:"completionTimeMs,omitempty" db:"cvm_completion_time_ms"`
	TotalTokens      *int      `json:"totalTokens,omitempty" db:"cvm_total_tokens"`
	TotalTimeMs      *int      `json:"totalTimeMs,omitempty" db:"cvm_total_time_ms"`
	CreatedAt        time.Time `json:"createdAt" db:"cvm_created_at"`
}

// Conversation Repository Params & Results
type CreateConversationParams struct {
	ChatID      string
	PhoneNumber string
	ContactName *string
	IsGroup     bool
	GroupName   *string
}

type CreateConversationResult struct {
	dal.DbResult
	ConversationID int `json:"conversationId" db:"o_cnv_id"`
}

type LinkUserToConversationParams struct {
	ChatID         string
	IdentityNumber string
}

type LinkUserToConversationResult struct {
	dal.DbResult
}

type CreateConversationMessageParams struct {
	ConversationID   int
	MessageID        string
	FromMe           bool
	SenderName       *string
	MessageType      string
	Body             *string
	MediaURL         *string
	QuotedMessage    *string
	Timestamp        int64
	IsForwarded      bool
	QueueTimeMs      *int
	PromptTokens     *int
	PromptTimeMs     *int
	CompletionTokens *int
	CompletionTimeMs *int
	TotalTokens      *int
	TotalTimeMs      *int
}

type CreateConversationMessageResult struct {
	dal.DbResult
	MessageID int `json:"messageId" db:"o_cvm_id"`
}

// Conversation Repository & UseCase Interfaces
type ConversationRepository interface {
	GetByChatID(ctx context.Context, chatID string) (*Conversation, error)
	Create(ctx context.Context, params CreateConversationParams) (*CreateConversationResult, error)
	LinkUser(ctx context.Context, params LinkUserToConversationParams) (*LinkUserToConversationResult, error)
	GetHistory(ctx context.Context, chatID string, limit int) ([]ConversationMessage, error)
	CreateMessage(ctx context.Context, params CreateConversationMessageParams) (*CreateConversationMessageResult, error)
}

type ConversationUseCase interface {
	GetOrCreateConversation(ctx context.Context, chatID, phoneNumber string, contactName *string, isGroup bool, groupName *string) Result[*Conversation]
	LinkUserToConversation(ctx context.Context, chatID, identityNumber string) Result[Data]
	GetConversationHistory(ctx context.Context, chatID string, limit int) Result[[]ConversationMessage]
	StoreMessage(ctx context.Context, conversationID int, messageID string, fromMe bool, body string, timestamp int64) Result[Data]
	StoreMessageWithStats(ctx context.Context, params CreateConversationMessageParams) Result[Data]
}
