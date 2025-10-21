package request

import (
	"api-chatbot/domain"
)

// WhatsApp Admin Requests

type GetWhatsAppStatusRequest struct {
	domain.Base
	SessionName string `json:"sessionName" validate:"required,min=1" doc:"WhatsApp session name to check status for"`
}

type GetWhatsAppQRCodeRequest struct {
	domain.Base
	SessionName string `json:"sessionName" validate:"required,min=1" doc:"WhatsApp session name to get QR code for"`
}

type UpdateWhatsAppStatusRequest struct {
	domain.Base
	SessionName string  `json:"sessionName" validate:"required,min=1" doc:"WhatsApp session name"`
	PhoneNumber *string `json:"phoneNumber" validate:"omitempty" doc:"Phone number connected to session"`
	DeviceName  *string `json:"deviceName" validate:"omitempty" doc:"Device name"`
	Platform    *string `json:"platform" validate:"omitempty" doc:"Platform (android/ios/web)"`
	Connected   *bool   `json:"connected" validate:"omitempty" doc:"Connection status"`
}

type GetConversationHistoryRequest struct {
	domain.Base
	ChatID string `json:"chatId" validate:"required,min=1" doc:"WhatsApp chat ID to retrieve history for"`
	Limit  int    `json:"limit" validate:"omitempty,gte=1,lte=1000" doc:"Maximum number of messages to retrieve (default: 50)"`
}
