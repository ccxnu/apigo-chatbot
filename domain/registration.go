package domain

import (
	"context"
	"time"

	"api-chatbot/api/dal"
)

// PendingRegistration represents a user awaiting OTP verification
type PendingRegistration struct {
	ID               int       `json:"id" db:"pending_id"`
	IdentityNumber   string    `json:"identityNumber" db:"identity_number"`
	WhatsApp         string    `json:"whatsapp" db:"whatsapp"`
	Name             string    `json:"name" db:"name"`
	Email            string    `json:"email" db:"email"`
	Phone            string    `json:"phone" db:"phone"`
	Role             string    `json:"role" db:"role"`
	UserType         string    `json:"userType" db:"user_type"` // 'institute' or 'external'
	Details          Data      `json:"details" db:"details"`
	RegistrationStep string    `json:"registrationStep" db:"registration_step"` // Current step in registration flow
	OTPExpiresAt     time.Time `json:"otpExpiresAt" db:"otp_expires_at"`
	OTPAttempts      int       `json:"otpAttempts" db:"otp_attempts"`
	Verified         bool      `json:"verified" db:"verified"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
}

// OTPVerificationResult represents the result of OTP verification
type OTPVerificationResult struct {
	Success        bool    `json:"success" db:"success"`
	Code           string  `json:"code" db:"code"`
	Message        string  `json:"message" db:"message"`
	PendingID      *int    `json:"pendingId,omitempty" db:"pending_id"`
	IdentityNumber *string `json:"identityNumber,omitempty" db:"identity_number"`
	Name           *string `json:"name,omitempty" db:"name"`
	Email          *string `json:"email,omitempty" db:"email"`
	Phone          *string `json:"phone,omitempty" db:"phone"`
	Role           *string `json:"role,omitempty" db:"role"`
	UserType       *string `json:"userType,omitempty" db:"user_type"`
	Details        Data    `json:"details,omitempty" db:"details"`
}

// CreatePendingRegistrationParams parameters for creating pending registration
type CreatePendingRegistrationParams struct {
	IdentityNumber   string
	WhatsApp         string
	Name             string
	Email            string
	Phone            string
	Role             string
	UserType         string // 'institute' or 'external'
	Details          Data
	RegistrationStep string // Current step in registration flow
	OTPCode          string
	OTPExpiresAt     time.Time
}

type CreatePendingRegistrationResult struct {
	dal.DbResult
	PendingID int `json:"pendingId" db:"o_pending_id"`
}

type DeletePendingRegistrationResult struct {
	dal.DbResult
}

// VerifyOTPParams parameters for OTP verification
type VerifyOTPParams struct {
	WhatsApp  string
	OTPCode   string
	IPAddress *string
}

// RegistrationRepository interface for registration data access
type RegistrationRepository interface {
	CreatePendingRegistration(ctx context.Context, params CreatePendingRegistrationParams) (*CreatePendingRegistrationResult, error)
	VerifyOTP(ctx context.Context, params VerifyOTPParams) (*OTPVerificationResult, error)
	GetPendingByWhatsApp(ctx context.Context, whatsapp string) (*PendingRegistration, error)
	UpdateRegistrationStep(ctx context.Context, whatsapp, step string) error
	DeletePendingRegistration(ctx context.Context, pendingID int) error
	CleanupExpiredRegistrations(ctx context.Context) (int, error)
}

// RegistrationUseCase interface for registration business logic
type RegistrationUseCase interface {
	// InitiateRegistration starts the registration process
	// Returns OTP code that should be sent to user's email
	InitiateRegistration(ctx context.Context, whatsapp, identityNumber string) Result[*PendingRegistration]

	// InitiateExternalRegistration starts external user registration (without email initially)
	InitiateExternalRegistration(ctx context.Context, whatsapp, identityNumber string) Result[*PendingRegistration]

	// InitiateRegistrationWithSelection creates pending for user type selection (when API fails)
	InitiateRegistrationWithSelection(ctx context.Context, whatsapp, identityNumber string) Result[*PendingRegistration]

	// InitiatePendingForCedula creates pending registration waiting for cedula input
	InitiatePendingForCedula(ctx context.Context, whatsapp string) Result[*PendingRegistration]

	// UpdatePendingWithRole updates pending registration with user's selected role
	UpdatePendingWithRole(ctx context.Context, whatsapp, role string) Result[*PendingRegistration]

	// CompleteExternalRegistration completes external registration with name and email
	CompleteExternalRegistration(ctx context.Context, whatsapp, name, email string) Result[*PendingRegistration]

	// VerifyAndRegister verifies OTP and completes user registration
	VerifyAndRegister(ctx context.Context, whatsapp, otpCode string) Result[*WhatsAppUser]

	// GetPendingRegistration retrieves pending registration by WhatsApp
	GetPendingRegistration(ctx context.Context, whatsapp string) Result[*PendingRegistration]

	// ResendOTP generates and sends a new OTP code
	ResendOTP(ctx context.Context, whatsapp string) Result[*PendingRegistration]

	// CancelPendingRegistration cancels/deletes a pending registration
	CancelPendingRegistration(ctx context.Context, whatsapp string) Result[bool]
}

// OTPMailer interface for sending OTP emails
type OTPMailer interface {
	SendOTPEmail(ctx context.Context, email, name, otpCode, userType string) error
}
