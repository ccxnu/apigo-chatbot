package route

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	"api-chatbot/domain"
)

// Huma response types for WhatsApp admin
type GetWhatsAppStatusResponse struct {
	Body domain.Result[*domain.WhatsAppSession]
}

type GetWhatsAppQRCodeResponse struct {
	Body domain.Result[domain.Data]
}

type UpdateWhatsAppStatusResponse struct {
	Body domain.Result[domain.Data]
}

type GetConversationHistoryResponse struct {
	Body domain.Result[[]domain.ConversationMessage]
}

func NewWhatsAppAdminRouter(
	sessionUseCase domain.WhatsAppSessionUseCase,
	convUseCase domain.ConversationUseCase,
	mux *http.ServeMux,
	humaAPI huma.API,
) {
	// Huma documented routes with /admin/whatsapp/ prefix
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-whatsapp-status",
		Method:      "POST",
		Path:        "/admin/whatsapp/status",
		Summary:     "Get WhatsApp connection status",
		Description: "Retrieves the current WhatsApp session status including connection state and device info",
		Tags:        []string{"Admin", "WhatsApp"},
	}, func(ctx context.Context, input *struct {
		Body request.GetWhatsAppStatusRequest
	}) (*GetWhatsAppStatusResponse, error) {
		result := sessionUseCase.GetSessionStatus(ctx, input.Body.SessionName)
		return &GetWhatsAppStatusResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-whatsapp-qr-code",
		Method:      "POST",
		Path:        "/admin/whatsapp/qr-code",
		Summary:     "Get WhatsApp QR code",
		Description: "Retrieves the QR code for WhatsApp authentication. Admin scans this to connect the bot.",
		Tags:        []string{"Admin", "WhatsApp"},
	}, func(ctx context.Context, input *struct {
		Body request.GetWhatsAppQRCodeRequest
	}) (*GetWhatsAppQRCodeResponse, error) {
		result := sessionUseCase.GetQRCode(ctx, input.Body.SessionName)
		return &GetWhatsAppQRCodeResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "update-whatsapp-status",
		Method:      "POST",
		Path:        "/admin/whatsapp/update-status",
		Summary:     "Update WhatsApp connection status",
		Description: "Updates WhatsApp session connection status (internal use)",
		Tags:        []string{"Admin", "WhatsApp"},
	}, func(ctx context.Context, input *struct {
		Body request.UpdateWhatsAppStatusRequest
	}) (*UpdateWhatsAppStatusResponse, error) {
		params := domain.UpdateSessionStatusParams{
			SessionName: input.Body.SessionName,
		}
		// Handle optional fields
		if input.Body.PhoneNumber != nil {
			params.PhoneNumber = *input.Body.PhoneNumber
		}
		if input.Body.DeviceName != nil {
			params.DeviceName = *input.Body.DeviceName
		}
		if input.Body.Platform != nil {
			params.Platform = *input.Body.Platform
		}
		if input.Body.Connected != nil {
			params.Connected = *input.Body.Connected
		}
		result := sessionUseCase.UpdateConnectionStatus(ctx, params)
		return &UpdateWhatsAppStatusResponse{Body: result}, nil
	})

	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-conversation-history",
		Method:      "POST",
		Path:        "/admin/whatsapp/conversation-history",
		Summary:     "Get conversation message history",
		Description: "Retrieves message history for a specific WhatsApp conversation",
		Tags:        []string{"Admin", "WhatsApp", "Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.GetConversationHistoryRequest
	}) (*GetConversationHistoryResponse, error) {
		result := convUseCase.GetConversationHistory(ctx, input.Body.ChatID, input.Body.Limit)
		return &GetConversationHistoryResponse{Body: result}, nil
	})
}
