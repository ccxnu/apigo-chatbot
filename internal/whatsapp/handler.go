package whatsapp

import (
	"context"

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
	handlers []MessageHandler
}

// NewMessageDispatcher creates a new message dispatcher
func NewMessageDispatcher(handlers []MessageHandler) *MessageDispatcher {
	dispatcher := &MessageDispatcher{
		handlers: handlers,
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
	for _, handler := range d.handlers {
		if handler.Match(ctx, msg) {
			return handler.Handle(ctx, msg)
		}
	}

	// No handler matched - could use a fallback handler here
	return nil
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
	Client     *Client
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
