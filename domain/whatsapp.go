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
	PhoneNumber *string    `json:"phoneNumber,omitempty" db:"wss_phone_number"`
	DeviceName  *string    `json:"deviceName,omitempty" db:"wss_device_name"`
	Platform    *string    `json:"platform,omitempty" db:"wss_platform"`
	QRCode      *string    `json:"qrCode,omitempty" db:"wss_qr_code"`
	Connected   bool       `json:"connected" db:"wss_connected"`
	LastSeen    *time.Time `json:"lastSeen,omitempty" db:"wss_last_seen"`
	SessionData Data       `json:"sessionData" db:"wss_session_data"`
	Active      bool       `json:"active" db:"wss_active"`
	CreatedAt   time.Time  `json:"createdAt" db:"wss_created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"wss_updated_at"`
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

type UpdateSessionStatusParams struct {
	SessionName string
	PhoneNumber *string
	DeviceName  *string
	Platform    *string
	Connected   bool
}

type UpdateSessionStatusResult struct {
	dal.DbResult
}

// WhatsApp Repository & UseCase Interfaces

type WhatsAppSessionRepository interface {
	GetBySessionName(ctx context.Context, sessionName string) (*WhatsAppSession, error)
	UpdateStatus(ctx context.Context, params UpdateSessionStatusParams) (*UpdateSessionStatusResult, error)
	UpdateQRCode(ctx context.Context, sessionName, qrCode string) error
	GetActiveSession(ctx context.Context) (*WhatsAppSession, error)
}

type WhatsAppSessionUseCase interface {
	GetSessionStatus(ctx context.Context, sessionName string) Result[*WhatsAppSession]
	GetQRCode(ctx context.Context, sessionName string) Result[Data]
	UpdateConnectionStatus(ctx context.Context, params UpdateSessionStatusParams) Result[Data]
	UpdateQRCode(ctx context.Context, sessionName, qrCode string) error
}

// User represents a WhatsApp user (student or professor)
type WhatsAppUser struct {
	ID             int       `json:"id" db:"usr_id"`
	IdentityNumber string    `json:"identityNumber" db:"usr_identity_number"`
	Name           string    `json:"name" db:"usr_name"`
	Email          string    `json:"email" db:"usr_email"`
	Phone          string    `json:"phone" db:"usr_phone"`
	Role           string    `json:"role" db:"usr_rol"`
	Details        Data      `json:"details" db:"usr_details"`
	WhatsApp       string    `json:"whatsapp" db:"usr_whatsapp"`
	Active         bool      `json:"active" db:"usr_active"`
	CreatedAt      time.Time `json:"createdAt" db:"usr_created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"usr_updated_at"`
}

// CreateUserParams parameters for creating a new user
type CreateUserParams struct {
	IdentityNumber string
	Name           string
	Email          string
	Phone          string
	Role           string
	WhatsApp       string
	Details        Data
}

type CreateUserResult struct {
	dal.DbResult
	UserID int `json:"userId" db:"o_usr_id"`
}

type UpdateUserWhatsAppParams struct {
	IdentityNumber string
	WhatsApp       string
}

// WhatsAppUserRepository interface for user data access
type WhatsAppUserRepository interface {
	GetByIdentity(ctx context.Context, identityNumber string) (*WhatsAppUser, error)
	GetByWhatsApp(ctx context.Context, whatsapp string) (*WhatsAppUser, error)
	Create(ctx context.Context, params CreateUserParams) (*CreateUserResult, error)
	UpdateWhatsApp(ctx context.Context, params UpdateUserWhatsAppParams) error
}

// WhatsAppUserUseCase interface for user business logic
type WhatsAppUserUseCase interface {
	GetOrRegisterUser(ctx context.Context, whatsapp string, identityNumber string) Result[*WhatsAppUser]
	ValidateWithInstituteAPI(ctx context.Context, identityNumber string) Result[*InstituteUserData]
	GetUserByWhatsApp(ctx context.Context, whatsapp string) Result[*WhatsAppUser]
}

// InstituteUserData represents data from institute validation API
type InstituteUserData struct {
	IdentityNumber string `json:"identityNumber"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Role           string `json:"role"` // "ROLE_STUDENT" or "ROLE_PROFESSOR"
	IsValid        bool   `json:"isValid"`
}
