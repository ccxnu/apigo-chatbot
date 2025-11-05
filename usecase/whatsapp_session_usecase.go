package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
	"api-chatbot/internal/whatsapp"
)

type whatsAppSessionUseCase struct {
	sessionRepo    d.WhatsAppSessionRepository
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewWhatsAppSessionUseCase(
	sessionRepo d.WhatsAppSessionRepository,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.WhatsAppSessionUseCase {
	return &whatsAppSessionUseCase{
		sessionRepo:    sessionRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

func (u *whatsAppSessionUseCase) GetSessionStatus(c context.Context, sessionName string) d.Result[*d.WhatsAppSession] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	session, err := u.sessionRepo.GetBySessionName(ctx, sessionName)
	if err != nil {
		logger.LogError(ctx, "Failed to fetch WhatsApp session from database", err,
			"operation", "GetSessionStatus",
			"sessionName", sessionName,
		)
		return d.Error[*d.WhatsAppSession](u.paramCache, "ERR_INTERNAL_DB")
	}

	if session == nil {
		logger.LogWarn(ctx, "WhatsApp session not found",
			"operation", "GetSessionStatus",
			"sessionName", sessionName,
		)
		return d.Error[*d.WhatsAppSession](u.paramCache, "ERR_WHATSAPP_SESSION_NOT_FOUND")
	}

	return d.Success(session)
}

func (u *whatsAppSessionUseCase) GetQRCode(c context.Context, sessionName string) d.Result[d.Data] {
	// Get QR code from in-memory service manager instead of database
	manager := whatsapp.GetManager()
	qrCode := manager.GetCurrentQR()
	connected := manager.IsConnected()

	logger.LogInfo(c, "Retrieved QR code from in-memory WhatsApp service",
		"operation", "GetQRCode",
		"sessionName", sessionName,
		"hasQRCode", qrCode != "",
		"qrCodeLength", len(qrCode),
		"connected", connected,
	)

	if qrCode != "" {
		logger.LogInfo(c, "Returning QR code to client",
			"qr_preview", qrCode[:min(50, len(qrCode))],
		)
	} else {
		logger.LogWarn(c, "No QR code available in memory",
			"connected", connected,
			"sessionName", sessionName,
		)
	}

	return d.Success(d.Data{
		"qrCode":    qrCode,
		"connected": connected,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (u *whatsAppSessionUseCase) UpdateConnectionStatus(c context.Context, params d.UpdateSessionStatusParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.sessionRepo.UpdateStatus(ctx, params)
	if err != nil || result == nil {
		logger.LogError(ctx, "Failed to update WhatsApp session status in database", err,
			"operation", "UpdateConnectionStatus",
			"sessionName", params.SessionName,
		)
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(ctx, "WhatsApp session status update failed with business logic error",
			"operation", "UpdateConnectionStatus",
			"code", result.Code,
			"sessionName", params.SessionName,
		)
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}

func (u *whatsAppSessionUseCase) UpdateQRCode(c context.Context, sessionName, qrCode string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err := u.sessionRepo.UpdateQRCode(ctx, sessionName, qrCode)
	if err != nil {
		logger.LogError(ctx, "Failed to update QR code in database", err,
			"operation", "UpdateQRCode",
			"sessionName", sessionName,
		)
		return err
	}

	logger.LogInfo(ctx, "QR code saved to database",
		"operation", "UpdateQRCode",
		"sessionName", sessionName,
	)

	return nil
}

func (u *whatsAppSessionUseCase) Logout(c context.Context, sessionName string) d.Result[d.Data] {
	manager := whatsapp.GetManager()

	if err := manager.Logout(c); err != nil {
		logger.LogError(c, "Failed to logout from WhatsApp", err,
			"operation", "Logout",
			"sessionName", sessionName,
		)
		return d.Error[d.Data](u.paramCache, "ERR_WHATSAPP_LOGOUT_FAILED")
	}

	logger.LogInfo(c, "WhatsApp logged out successfully",
		"operation", "Logout",
		"sessionName", sessionName,
	)

	return d.Success(d.Data{
		"message": "Logged out successfully. Device pairing cleared.",
	})
}

func (u *whatsAppSessionUseCase) Reconnect(c context.Context, sessionName string) d.Result[d.Data] {
	manager := whatsapp.GetManager()

	if err := manager.Reconnect(c); err != nil {
		logger.LogError(c, "Failed to reconnect WhatsApp", err,
			"operation", "Reconnect",
			"sessionName", sessionName,
		)
		return d.Error[d.Data](u.paramCache, "ERR_WHATSAPP_RECONNECT_FAILED")
	}

	logger.LogInfo(c, "WhatsApp reconnected successfully",
		"operation", "Reconnect",
		"sessionName", sessionName,
	)

	return d.Success(d.Data{
		"message": "Reconnecting... New QR code will be generated.",
	})
}
