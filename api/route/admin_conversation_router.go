package route

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

type GetConversationsResponse struct {
	Body d.Result[[]d.AdminConversationListItem]
}

type GetConversationMessagesResponse struct {
	Body d.Result[[]d.AdminConversationMessage]
}

type SendAdminMessageResponse struct {
	Body d.Result[d.Data]
}

type MarkMessagesReadResponse struct {
	Body d.Result[d.Data]
}

type BlockUserResponse struct {
	Body d.Result[d.Data]
}

type DeleteConversationResponse struct {
	Body d.Result[d.Data]
}

type SetTemporaryResponse struct {
	Body d.Result[d.Data]
}

func SetupAdminConversationRoutes(humaAPI huma.API, adminConvUC d.AdminConversationUseCase) {

	// GET /api/v1/admin/conversations - List all conversations
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-admin-conversations",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/get-all",
		Summary:     "Get all conversations for admin panel",
		Description: "Retrieves paginated list of conversations with filters (WhatsApp-like view)",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.GetConversationsRequest
	}) (*GetConversationsResponse, error) {
		filter := input.Body.Filter
		if filter == "" {
			filter = "all"
		}

		limit := input.Body.Limit
		if limit == 0 {
			limit = 50
		}

		offset := input.Body.Offset

		result := adminConvUC.GetAllConversations(ctx, filter, limit, offset)
		return &GetConversationsResponse{Body: result}, nil
	})

	// GET /api/v1/admin/conversations/:id/messages - Get conversation messages
	huma.Register(humaAPI, huma.Operation{
		OperationID: "get-conversation-messages",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/get-messages",
		Summary:     "Get conversation message history",
		Description: "Retrieves all messages for a conversation",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.GetConversationMessagesRequest
	}) (*GetConversationMessagesResponse, error) {
		limit := input.Body.Limit
		if limit == 0 {
			limit = 100
		}

		result := adminConvUC.GetConversationMessages(ctx, input.Body.ConversationID, limit)
		return &GetConversationMessagesResponse{Body: result}, nil
	})

	// POST /api/v1/admin/conversations/:id/send - Admin sends a message
	huma.Register(humaAPI, huma.Operation{
		OperationID: "send-admin-message",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/send",
		Summary:     "Admin sends a message",
		Description: "Admin sends a message in a conversation (also sends via WhatsApp)",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.SendAdminMessageRequest
	}) (*SendAdminMessageResponse, error) {
		// TODO: Extract admin ID from JWT token
		adminID := 1 // Placeholder

		params := d.SendAdminMessageParams{
			ConversationID: input.Body.ConversationID,
			AdminID:        adminID,
			MessageID:      generateMessageID(),
			Body:           input.Body.Message,
		}

		result := adminConvUC.SendAdminMessage(ctx, params)
		return &SendAdminMessageResponse{Body: result}, nil
	})

	// POST /api/v1/admin/conversations/:id/mark-read - Mark messages as read
	huma.Register(humaAPI, huma.Operation{
		OperationID: "mark-messages-read",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/mark-read",
		Summary:     "Mark messages as read",
		Description: "Marks all messages in a conversation as read by admin",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.MarkMessagesReadRequest
	}) (*MarkMessagesReadResponse, error) {
		result := adminConvUC.MarkMessagesAsRead(ctx, input.Body.ConversationID)
		return &MarkMessagesReadResponse{Body: result}, nil
	})

	// POST /api/v1/admin/users/:id/block - Block a user
	huma.Register(humaAPI, huma.Operation{
		OperationID: "block-user",
		Method:      "POST",
		Path:        "/api/v1/admin/users/block",
		Summary:     "Block or unblock a user",
		Description: "Blocks or unblocks a user from using the chatbot",
		Tags:        []string{"Admin - Users"},
	}, func(ctx context.Context, input *struct {
		Body request.BlockUserRequest
	}) (*BlockUserResponse, error) {
		// TODO: Extract admin ID from JWT token
		adminID := 1 // Placeholder

		params := d.BlockUserParams{
			UserID:  input.Body.UserID,
			Blocked: input.Body.Blocked,
			AdminID: adminID,
			Reason:  input.Body.Reason,
		}

		result := adminConvUC.BlockUser(ctx, params)
		return &BlockUserResponse{Body: result}, nil
	})

	// DELETE /api/v1/admin/conversations/:id - Delete a conversation
	huma.Register(humaAPI, huma.Operation{
		OperationID: "delete-conversation",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/delete",
		Summary:     "Delete a conversation",
		Description: "Soft deletes a conversation (marks as inactive)",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.DeleteConversationRequest
	}) (*DeleteConversationResponse, error) {
		result := adminConvUC.DeleteConversation(ctx, input.Body.ConversationID)
		return &DeleteConversationResponse{Body: result}, nil
	})

	// POST /api/v1/admin/conversations/:id/temporary - Set conversation as temporary
	huma.Register(humaAPI, huma.Operation{
		OperationID: "set-conversation-temporary",
		Method:      "POST",
		Path:        "/api/v1/admin/conversations/set-temporary",
		Summary:     "Enable/disable temporary conversation",
		Description: "Sets conversation to auto-delete after specified hours",
		Tags:        []string{"Admin - Conversations"},
	}, func(ctx context.Context, input *struct {
		Body request.SetTemporaryRequest
	}) (*SetTemporaryResponse, error) {
		hoursUntilExpiry := input.Body.HoursUntilExpiry
		if hoursUntilExpiry == 0 {
			hoursUntilExpiry = 24 // Default 24 hours
		}

		params := d.SetTemporaryParams{
			ConversationID:   input.Body.ConversationID,
			Temporary:        input.Body.Temporary,
			HoursUntilExpiry: hoursUntilExpiry,
		}

		result := adminConvUC.SetConversationTemporary(ctx, params)
		return &SetTemporaryResponse{Body: result}, nil
	})
}

// generateMessageID generates a unique message ID for admin messages
func generateMessageID() string {
	return fmt.Sprintf("admin-%d", time.Now().UnixNano())
}
