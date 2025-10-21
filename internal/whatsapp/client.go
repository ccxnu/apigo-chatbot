package whatsapp

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Client wraps whatsmeow client for WhatsApp integration
type Client struct {
	WAClient  *whatsmeow.Client
	Container *sqlstore.Container
	SessionID string
}

// Config holds configuration for WhatsApp client
type Config struct {
	DBDialect string // "postgres" or "sqlite3"
	DBAddress string // connection string
	SessionID string // unique session identifier
	LogLevel  string // "ERROR", "WARN", "INFO", "DEBUG"
}

// NewClient creates a new WhatsApp client
func NewClient(cfg Config) (*Client, error) {
	// Create database container for session storage
	container, err := sqlstore.New(cfg.DBDialect, cfg.DBAddress, waLog.Noop)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQL store: %w", err)
	}

	// Get device store (creates if doesn't exist)
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Create whatsmeow client
	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)

	return &Client{
		WAClient:  waClient,
		Container: container,
		SessionID: cfg.SessionID,
	}, nil
}

// Connect establishes connection to WhatsApp
func (c *Client) Connect(ctx context.Context) error {
	if c.WAClient.Store.ID == nil {
		// Not logged in, need QR code
		return fmt.Errorf("device not paired, QR code required")
	}

	return c.WAClient.Connect()
}

// Disconnect closes the WhatsApp connection
func (c *Client) Disconnect() {
	if c.WAClient != nil {
		c.WAClient.Disconnect()
	}
}

// IsConnected checks if client is connected
func (c *Client) IsConnected() bool {
	return c.WAClient != nil && c.WAClient.IsConnected()
}

// IsLoggedIn checks if device is paired
func (c *Client) IsLoggedIn() bool {
	return c.WAClient != nil && c.WAClient.Store.ID != nil
}

// GetQRChannel returns channel for QR code updates
func (c *Client) GetQRChannel() (<-chan whatsmeow.QRChannelItem, error) {
	if c.WAClient.Store.ID != nil {
		return nil, fmt.Errorf("already logged in")
	}

	qrChan, err := c.WAClient.GetQRChannel(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get QR channel: %w", err)
	}

	// Connect to start QR generation
	err = c.WAClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect for QR: %w", err)
	}

	return qrChan, nil
}

// SendTextMessage sends a text message to a chat
func (c *Client) SendTextMessage(chatJID types.JID, text string) (types.MessageID, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected to WhatsApp")
	}

	msg := &waE2E.Message{
		Conversation: &text,
	}

	resp, err := c.WAClient.SendMessage(context.Background(), chatJID, msg)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return resp.ID, nil
}

// GetDeviceInfo returns information about the connected device
func (c *Client) GetDeviceInfo() *DeviceInfo {
	if c.WAClient.Store.ID == nil {
		return nil
	}

	return &DeviceInfo{
		PhoneNumber: c.WAClient.Store.ID.User,
		DeviceName:  c.WAClient.Store.PushName,
		Platform:    c.WAClient.Store.Platform,
		Connected:   c.IsConnected(),
	}
}

// DeviceInfo holds information about the WhatsApp device
type DeviceInfo struct {
	PhoneNumber string
	DeviceName  string
	Platform    string
	Connected   bool
}

// Logout logs out from WhatsApp
func (c *Client) Logout() error {
	if c.WAClient.Store.ID == nil {
		return fmt.Errorf("not logged in")
	}

	err := c.WAClient.Logout()
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}
