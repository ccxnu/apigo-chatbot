package whatsapp

import (
	"context"
	"log/slog"

	"api-chatbot/domain"
)

// MessageHandler interface defines the contract for message handlers
// Similar to NestJS handler pattern with Match() and Handle()
type MessageHandler interface {
	// Match determines if this handler should process the message
	Match(ctx context.Context, msg *domain.IncomingMessage) bool

	// Handle processes the message
	Handle(ctx context.Context, msg *domain.IncomingMessage) error

	// Priority returns handler priority (higher = executed first)
	Priority() int
}

// MessageDispatcher routes messages to appropriate handlers
type MessageDispatcher struct {
	handlers   []MessageHandler
	paramCache domain.ParameterCache
	client     *Client
}

// NewMessageDispatcher creates a new message dispatcher
func NewMessageDispatcher(handlers []MessageHandler, paramCache domain.ParameterCache, client *Client) *MessageDispatcher {
	dispatcher := &MessageDispatcher{
		handlers:   handlers,
		paramCache: paramCache,
		client:     client,
	}
	dispatcher.sortHandlers()
	return dispatcher
}

// RegisterHandler adds a handler to the dispatcher
func (d *MessageDispatcher) RegisterHandler(handler MessageHandler) {
	d.handlers = append(d.handlers, handler)
	// Sort by priority (highest first)
	d.sortHandlers()
}

// Dispatch routes a message to the first matching handler
func (d *MessageDispatcher) Dispatch(ctx context.Context, msg *domain.IncomingMessage) error {
	// Check if chatbot is active (hot-reloadable via parameter cache)
	if !d.isChatbotActive() {
		slog.Info("Chatbot is deactivated, skipping message processing", "messageID", msg.MessageID)
		d.sendDeactivatedMessage(msg.ChatID)
		return nil
	}

	slog.Info("Dispatching message to handlers",
		"handlersCount", len(d.handlers),
		"messageID", msg.MessageID,
	)

	for i, handler := range d.handlers {
		matched := handler.Match(ctx, msg)
		slog.Info("Handler check",
			"handlerIndex", i,
			"priority", handler.Priority(),
			"matched", matched,
		)

		if matched {
			slog.Info("Handler matched, executing",
				"handlerIndex", i,
				"priority", handler.Priority(),
			)
			return handler.Handle(ctx, msg)
		}
	}

	// No handler matched - could use a fallback handler here
	slog.Warn("No handler matched for message", "messageID", msg.MessageID)
	return nil
}

// isChatbotActive checks if the chatbot is enabled via parameter cache
func (d *MessageDispatcher) isChatbotActive() bool {
	param, exists := d.paramCache.Get("CHATBOT_ACTIVE")
	if !exists {
		// Default to active if parameter doesn't exist
		return true
	}

	data, err := param.GetDataAsMap()
	if err != nil {
		slog.Warn("Failed to parse CHATBOT_ACTIVE parameter, defaulting to active", "error", err)
		return true
	}

	active, ok := data["active"].(bool)
	if !ok {
		slog.Warn("CHATBOT_ACTIVE parameter has invalid format, defaulting to active")
		return true
	}

	return active
}

// sendDeactivatedMessage sends a message when chatbot is deactivated
func (d *MessageDispatcher) sendDeactivatedMessage(chatID string) {
	if d.client == nil {
		return
	}

	param, exists := d.paramCache.Get("CHATBOT_DEACTIVATED_MESSAGE")
	message := "🔧 El chatbot está temporalmente desactivado por mantenimiento. Por favor, intenta más tarde."

	if exists {
		data, err := param.GetDataAsMap()
		if err == nil {
			if msg, ok := data["message"].(string); ok {
				message = msg
			}
		}
	}

	err := d.client.SendText(chatID, message)
	if err != nil {
		slog.Error("Failed to send deactivated message", "error", err, "chatID", chatID)
	}
}

// sortHandlers sorts handlers by priority (descending)
func (d *MessageDispatcher) sortHandlers() {
	for i := 0; i < len(d.handlers); i++ {
		for j := i + 1; j < len(d.handlers); j++ {
			if d.handlers[j].Priority() > d.handlers[i].Priority() {
				d.handlers[i], d.handlers[j] = d.handlers[j], d.handlers[i]
			}
		}
	}
}

// BaseHandler provides common functionality for handlers
type BaseHandler struct {
	Client      *Client
	ConvUseCase domain.ConversationUseCase
}

// MessageFilter provides utility functions for filtering messages
type MessageFilter struct{}

// NewMessageFilter creates a new message filter
func NewMessageFilter() *MessageFilter {
	return &MessageFilter{}
}

// IsCommand checks if message starts with a command prefix
func (f *MessageFilter) IsCommand(msg *domain.IncomingMessage, command string) bool {
	if len(msg.Body) == 0 {
		return false
	}

	// Check for /command format
	if msg.Body[0] == '/' && len(msg.Body) > 1 {
		cmd := msg.Body[1:]
		if len(cmd) >= len(command) {
			return cmd[:len(command)] == command
		}
	}

	return false
}

// IsFromMe checks if message was sent by the bot
func (f *MessageFilter) IsFromMe(msg *domain.IncomingMessage) bool {
	return msg.FromMe
}

// IsGroup checks if message is from a group
func (f *MessageFilter) IsGroup(msg *domain.IncomingMessage) bool {
	return msg.IsGroup
}

// IsTextMessage checks if message is a text message
func (f *MessageFilter) IsTextMessage(msg *domain.IncomingMessage) bool {
	return msg.MessageType == "text" || msg.MessageType == ""
}

// HasBody checks if message has text content
func (f *MessageFilter) HasBody(msg *domain.IncomingMessage) bool {
	return len(msg.Body) > 0
}
