package usecase

import (
	"context"
	"time"

	"api-chatbot/domain"
)

type whatsAppSessionUseCase struct {
	sessionRepo    domain.WhatsAppSessionRepository
	paramCache     domain.ParameterCache
	contextTimeout time.Duration
}

func NewWhatsAppSessionUseCase(
	sessionRepo domain.WhatsAppSessionRepository,
	paramCache domain.ParameterCache,
	timeout time.Duration,
) domain.WhatsAppSessionUseCase {
	return &whatsAppSessionUseCase{
		sessionRepo:    sessionRepo,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// getErrorMessage retrieves error message from parameter cache
func (u *whatsAppSessionUseCase) getErrorMessage(errorCode string) string {
	if param, exists := u.paramCache.Get(errorCode); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if message, ok := data["message"].(string); ok {
				return message
			}
		}
	}
	return "Ha ocurrido un error"
}

func (u *whatsAppSessionUseCase) GetSessionStatus(c context.Context, sessionName string) domain.Result[*domain.WhatsAppSession] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	session, err := u.sessionRepo.GetBySessionName(ctx, sessionName)
	if err != nil {
		return domain.Result[*domain.WhatsAppSession]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if session == nil {
		return domain.Result[*domain.WhatsAppSession]{
			Success: false,
			Code:    "ERR_WHATSAPP_SESSION_NOT_FOUND",
			Info:    u.getErrorMessage("ERR_WHATSAPP_SESSION_NOT_FOUND"),
			Data:    nil,
		}
	}

	return domain.Result[*domain.WhatsAppSession]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    session,
	}
}

func (u *whatsAppSessionUseCase) GetQRCode(c context.Context, sessionName string) domain.Result[domain.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	session, err := u.sessionRepo.GetBySessionName(ctx, sessionName)
	if err != nil {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if session == nil {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    "ERR_WHATSAPP_SESSION_NOT_FOUND",
			Info:    u.getErrorMessage("ERR_WHATSAPP_SESSION_NOT_FOUND"),
			Data:    nil,
		}
	}

	return domain.Result[domain.Data]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data: domain.Data{
			"qrCode":    session.QRCode,
			"connected": session.Connected,
		},
	}
}

func (u *whatsAppSessionUseCase) UpdateConnectionStatus(c context.Context, params domain.UpdateSessionStatusParams) domain.Result[domain.Data] {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	result, err := u.sessionRepo.UpdateStatus(ctx, params)
	if err != nil || result == nil {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    "ERR_INTERNAL_DB",
			Info:    u.getErrorMessage("ERR_INTERNAL_DB"),
			Data:    nil,
		}
	}

	if !result.Success {
		return domain.Result[domain.Data]{
			Success: false,
			Code:    result.Code,
			Info:    u.getErrorMessage(result.Code),
			Data:    nil,
		}
	}

	return domain.Result[domain.Data]{
		Success: true,
		Code:    "OK",
		Info:    u.getErrorMessage("OK"),
		Data:    nil,
	}
}
