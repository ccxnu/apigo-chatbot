package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"

	d "api-chatbot/domain"
)

// Service manages WhatsApp client lifecycle and event handling
type Service struct {
	client      *Client
	dispatcher  *MessageDispatcher
	sessionUC   d.WhatsAppSessionUseCase
	sessionName string
	connectTime	int64
}

// NewService creates a new WhatsApp service
func NewService(
	deviceStore *store.Device,
	sessionName string,
	sessionUC d.WhatsAppSessionUseCase,
	handlers []MessageHandler,
	paramCache d.ParameterCache,
) (*Service, error) {
	cfg := Config{
		SessionName: sessionName,
		DeviceStore: deviceStore,
		LogLevel:    "INFO",
	}

	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create WhatsApp client: %w", err)
	}

	dispatcher := NewMessageDispatcher(handlers, paramCache, client)

	return &Service{
		client:      client,
		dispatcher:  dispatcher,
		sessionUC:   sessionUC,
		sessionName: sessionName,
	}, nil
}

// NewServiceWithClient creates a new WhatsApp service with an existing client
func NewServiceWithClient(
	client *Client,
	sessionName string,
	sessionUC d.WhatsAppSessionUseCase,
	handlers []MessageHandler,
	paramCache d.ParameterCache,
) (*Service, error) {
	dispatcher := NewMessageDispatcher(handlers, paramCache, client)

	return &Service{
		client:      client,
		dispatcher:  dispatcher,
		sessionUC:   sessionUC,
		sessionName: sessionName,
	}, nil
}

// Start initializes the WhatsApp connection and starts listening for events
func (s *Service) Start(ctx context.Context) error {
	// Register event handler
	s.client.WAClient.AddEventHandler(s.handleEvent)

	// Connect to WhatsApp
	if err := s.client.WAClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to WhatsApp: %w", err)
	}

	slog.Info("WhatsApp service started", "session", s.sessionName)
	return nil
}

// Stop gracefully disconnects the WhatsApp client
func (s *Service) Stop() {
	if s.client != nil && s.client.WAClient != nil {
		s.client.WAClient.Disconnect()
		slog.Info("WhatsApp service stopped", "session", s.sessionName)
	}
}

// GetQRChannel returns the channel for QR code events
func (s *Service) GetQRChannel() <-chan string {
	qrChan := make(chan string, 1)

	s.client.WAClient.AddEventHandler(func(evt any) {
		if qrEvt, ok := evt.(*events.QR); ok {
			qrChan <- qrEvt.Codes[0]
		}
	})

	return qrChan
}

// handleEvent processes all WhatsApp events
func (s *Service) handleEvent(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		s.handleIncomingMessage(v)

	case *events.QR:
		s.handleQRCode(v)

	case *events.Connected:
		s.handleConnected()

	case *events.Disconnected:
		s.handleDisconnected()

	case *events.PairSuccess:
		s.handlePairSuccess(v)
	}
}

// handleIncomingMessage processes incoming WhatsApp messages
func (s *Service) handleIncomingMessage(evt *events.Message) {

	msg := convertEventToMessage(evt)

	const timeThreshold = 5000
	if !evt.Info.IsFromMe && evt.Info.Timestamp.Unix() < (s.connectTime - timeThreshold) {
		slog.Info("Ignoring old message",
			"messageID", evt.Info.ID,
			"chatID", msg.ChatID,
			"timestamp", evt.Info.Timestamp.Unix(),
			"connectTime", s.connectTime,
		)
		return
	}

	ctx := context.Background()

	// Dispatch to handlers
	if err := s.dispatcher.Dispatch(ctx, msg); err != nil {
		slog.Error("Failed to dispatch message",
			"messageID", msg.MessageID,
			"chatID", msg.ChatID,
			"error", err,
		)
	}
}

// handleQRCode processes QR code events and updates database
func (s *Service) handleQRCode(evt *events.QR) {
	if len(evt.Codes) == 0 {
		return
	}

	qrCode := evt.Codes[0]

	slog.Info("QR code received - scan with WhatsApp app",
		"session", s.sessionName,
		"qr_length", len(qrCode),
	)

	// Save QR code to database
	ctx := context.Background()
	err := s.sessionUC.UpdateQRCode(ctx, s.sessionName, qrCode)
	if err != nil {
		slog.Error("Failed to save QR code to database",
			"session", s.sessionName,
			"error", err,
		)
	}

	// Print QR code to console for easy scanning
	printQRCodeToConsole(qrCode)
}

// handleConnected processes connection success events
func (s *Service) handleConnected() {
	slog.Info("WhatsApp connected", "session", s.sessionName)

	s.connectTime = time.Now().Unix()
	ctx := context.Background()

	// Get device info
	device := s.client.DeviceStore

	connected := true
	phoneNumber := device.ID.User
	params := d.UpdateSessionStatusParams{
		SessionName: s.sessionName,
		PhoneNumber: &phoneNumber,
		Connected:   connected,
	}

	result := s.sessionUC.UpdateConnectionStatus(ctx, params)
	if !result.Success {
		slog.Error("Failed to update connection status",
			"session", s.sessionName,
			"code", result.Code,
		)
	}
}

// handleDisconnected processes disconnection events
func (s *Service) handleDisconnected() {
	slog.Warn("WhatsApp disconnected", "session", s.sessionName)

	ctx := context.Background()

	connected := false
	params := d.UpdateSessionStatusParams{
		SessionName: s.sessionName,
		Connected:   connected,
	}

	result := s.sessionUC.UpdateConnectionStatus(ctx, params)
	if !result.Success {
		slog.Error("Failed to update disconnection status",
			"session", s.sessionName,
			"code", result.Code,
		)
	}
}

// handlePairSuccess processes successful pairing events
func (s *Service) handlePairSuccess(evt *events.PairSuccess) {
	slog.Info("WhatsApp pairing successful",
		"session", s.sessionName,
		"phone", evt.ID.User,
		"platform", evt.Platform,
	)

	ctx := context.Background()

	connected := true
	phoneNumber := evt.ID.User
	platform := evt.Platform
	params := d.UpdateSessionStatusParams{
		SessionName: s.sessionName,
		PhoneNumber: &phoneNumber,
		Platform:    &platform,
		Connected:   connected,
	}

	result := s.sessionUC.UpdateConnectionStatus(ctx, params)
	if !result.Success {
		slog.Error("Failed to update pairing status",
			"session", s.sessionName,
			"code", result.Code,
		)
	}
}

// convertEventToMessage converts whatsmeow event to domain IncomingMessage
func convertEventToMessage(evt *events.Message) *d.IncomingMessage {
	msg := &d.IncomingMessage{
		MessageID:   evt.Info.ID,
		ChatID:      evt.Info.Chat.String(),
		From:        evt.Info.Sender.String(),
		FromMe:      evt.Info.IsFromMe,
		Timestamp:   evt.Info.Timestamp.Unix(),
		IsGroup:     evt.Info.IsGroup,
		MessageType: string(evt.Message.GetConversation()),
	}

	// Extract text content
	if evt.Message.GetConversation() != "" {
		msg.Body = evt.Message.GetConversation()
		msg.MessageType = "text"
	} else if ext := evt.Message.GetExtendedTextMessage(); ext != nil {
		msg.Body = ext.GetText()
		msg.MessageType = "text"

		// Handle quoted messages
		if quoted := ext.GetContextInfo().GetQuotedMessage(); quoted != nil {
			msg.QuotedMessage = quoted.GetConversation()
		}
	}

	// Extract group info
	if evt.Info.IsGroup {
		msg.GroupName = evt.Info.PushName
	}

	// Extract media (images, videos, documents, etc.)
	if img := evt.Message.GetImageMessage(); img != nil {
		msg.MessageType = "image"
		msg.MediaURL = img.GetURL()
		msg.Body = img.GetCaption()
	} else if vid := evt.Message.GetVideoMessage(); vid != nil {
		msg.MessageType = "video"
		msg.MediaURL = vid.GetURL()
		msg.Body = vid.GetCaption()
	} else if doc := evt.Message.GetDocumentMessage(); doc != nil {
		msg.MessageType = "document"
		msg.MediaURL = doc.GetURL()
		msg.Body = doc.GetFileName()
	} else if aud := evt.Message.GetAudioMessage(); aud != nil {
		msg.MessageType = "audio"
		msg.MediaURL = aud.GetURL()
	}

	return msg
}

// printQRCodeToConsole displays the QR code in the terminal for easy scanning
func printQRCodeToConsole(code string) {
	qr, err := qrcode.New(code, qrcode.Medium)
	if err != nil {
		slog.Error("Failed to generate QR code", "error", err)
		return
	}

	fmt.Println("\n========================================")
	fmt.Println("SCAN THIS QR CODE WITH WHATSAPP:")
	fmt.Println("========================================")
	fmt.Println(qr.ToSmallString(false))
	fmt.Println("========================================")
}
