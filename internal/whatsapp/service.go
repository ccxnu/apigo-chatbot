package whatsapp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	currentQR   string // In-memory QR code (not persisted to database)
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

	// Check if device is paired
	if s.client.WAClient.Store.ID == nil {
		// Not paired - need QR code. Connect will trigger QR generation
		slog.Info("Device not paired, connecting to generate QR code", "session", s.sessionName)
		if err := s.client.WAClient.Connect(); err != nil {
			return fmt.Errorf("failed to connect for QR generation: %w", err)
		}
	} else {
		// Already paired - just connect
		slog.Info("Device already paired, connecting to WhatsApp", "session", s.sessionName)
		if err := s.client.WAClient.Connect(); err != nil {
			return fmt.Errorf("failed to connect to WhatsApp: %w", err)
		}
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

// Logout logs out from WhatsApp and clears device pairing
func (s *Service) Logout(ctx context.Context) error {
	if s.client == nil || s.client.WAClient == nil {
		return fmt.Errorf("client not initialized")
	}

	// Logout from WhatsApp
	if err := s.client.Logout(); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	// Clear QR code from memory
	s.currentQR = ""

	slog.Info("WhatsApp logged out successfully", "session", s.sessionName)
	return nil
}

// Reconnect disconnects and reconnects to generate a new QR code
func (s *Service) Reconnect(ctx context.Context) error {
	if s.client == nil || s.client.WAClient == nil {
		return fmt.Errorf("client not initialized")
	}

	// Disconnect first
	s.client.WAClient.Disconnect()
	s.currentQR = ""

	slog.Info("Reconnecting to generate new QR code", "session", s.sessionName)

	// Reconnect - will trigger QR if not paired
	if err := s.client.WAClient.Connect(); err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	return nil
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

// GetCurrentQR returns the current in-memory QR code
func (s *Service) GetCurrentQR() string {
	return s.currentQR
}

// IsConnected returns the current connection status
func (s *Service) IsConnected() bool {
	if s.client == nil || s.client.WAClient == nil {
		return false
	}
	return s.client.WAClient.IsConnected()
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

// handleQRCode processes QR code events (in-memory only, not saved to database)
func (s *Service) handleQRCode(evt *events.QR) {
	if len(evt.Codes) == 0 {
		return
	}

	qrCode := evt.Codes[0]
	s.currentQR = qrCode // Store in memory

	slog.Info("QR code received - scan with WhatsApp app",
		"session", s.sessionName,
		"qr_length", len(qrCode),
	)

	// QR codes are NOT saved to database - they are ephemeral
	// Frontend should call GetCurrentQR() to retrieve the latest QR code
	slog.Debug("QR code generated and stored in memory (not saved to database)",
		"session", s.sessionName,
	)
}

// handleConnected processes connection success events
func (s *Service) handleConnected() {
	slog.Info("WhatsApp connected", "session", s.sessionName)

	s.connectTime = time.Now().Unix()

	// Clear QR code from memory once connected
	s.currentQR = ""

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

	// Clear in-memory QR code
	s.currentQR = ""

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

	// Clear QR code from database on disconnect so a fresh one is generated on reconnect
	err := s.sessionUC.UpdateQRCode(ctx, s.sessionName, "")
	if err != nil {
		slog.Error("Failed to clear QR code from database",
			"session", s.sessionName,
			"error", err,
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
