package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetWhatsAppSession       = "fn_get_whatsapp_session"
	fnGetActiveWhatsAppSession = "fn_get_active_whatsapp_session"
	// Stored Procedures (Writes)
	spUpdateWhatsAppSessionStatus = "sp_update_whatsapp_session_status"
	spUpdateWhatsAppQRCode        = "sp_update_whatsapp_qr_code"
)

type whatsAppSessionRepository struct {
	dal *dal.DAL
}

func NewWhatsAppSessionRepository(dal *dal.DAL) d.WhatsAppSessionRepository {
	return &whatsAppSessionRepository{
		dal: dal,
	}
}

// GetBySessionName retrieves a WhatsApp session by name
func (r *whatsAppSessionRepository) GetBySessionName(ctx context.Context, sessionName string) (*d.WhatsAppSession, error) {
	sessions, err := dal.QueryRows[d.WhatsAppSession](r.dal, ctx, fnGetWhatsAppSession, sessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get WhatsApp session via %s: %w", fnGetWhatsAppSession, err)
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}

// GetActiveSession retrieves the currently active WhatsApp session
func (r *whatsAppSessionRepository) GetActiveSession(ctx context.Context) (*d.WhatsAppSession, error) {
	sessions, err := dal.QueryRows[d.WhatsAppSession](r.dal, ctx, fnGetActiveWhatsAppSession)
	if err != nil {
		return nil, fmt.Errorf("failed to get active WhatsApp session via %s: %w", fnGetActiveWhatsAppSession, err)
	}

	if len(sessions) == 0 {
		return nil, nil
	}

	return &sessions[0], nil
}

// UpdateStatus updates WhatsApp session connection status
func (r *whatsAppSessionRepository) UpdateStatus(ctx context.Context, params d.UpdateSessionStatusParams) (*d.UpdateSessionStatusResult, error) {
	result, err := dal.ExecProc[d.UpdateSessionStatusResult](
		r.dal,
		ctx,
		spUpdateWhatsAppSessionStatus,
		params.SessionName,
		params.PhoneNumber,
		params.DeviceName,
		params.Platform,
		params.Connected,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spUpdateWhatsAppSessionStatus, err)
	}

	return result, nil
}

// UpdateQRCode updates the QR code for a WhatsApp session
func (r *whatsAppSessionRepository) UpdateQRCode(ctx context.Context, sessionName, qrCode string) error {
	result, err := dal.ExecProc[d.UpdateSessionStatusResult](
		r.dal,
		ctx,
		spUpdateWhatsAppQRCode,
		sessionName,
		qrCode,
	)

	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", spUpdateWhatsAppQRCode, err)
	}

	if !result.Success {
		return fmt.Errorf("failed to update QR code: %s", result.Code)
	}

	return nil
}
