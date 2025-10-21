package usecase

import (
	"context"
	"time"

	d "api-chatbot/domain"
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
		return d.Error[*d.WhatsAppSession](u.paramCache, "ERR_INTERNAL_DB")
	}

	if session == nil {
		return d.Error[*d.WhatsAppSession](u.paramCache, "ERR_WHATSAPP_SESSION_NOT_FOUND")
	}

	return d.Success(session)
}

func (u *whatsAppSessionUseCase) GetQRCode(c context.Context, sessionName string) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	session, err := u.sessionRepo.GetBySessionName(ctx, sessionName)
	if err != nil {
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if session == nil {
		return d.Error[d.Data](u.paramCache, "ERR_WHATSAPP_SESSION_NOT_FOUND")
	}

	return d.Success(d.Data{
		"qrCode":    session.QRCode,
		"connected": session.Connected,
	})
}

func (u *whatsAppSessionUseCase) UpdateConnectionStatus(c context.Context, params d.UpdateSessionStatusParams) d.Result[d.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.sessionRepo.UpdateStatus(ctx, params)
	if err != nil || result == nil {
		return d.Error[d.Data](u.paramCache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		return d.Error[d.Data](u.paramCache, result.Code)
	}

	return d.Success(d.Data{})
}
