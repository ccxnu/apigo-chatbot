package handlers

//
// import (
// 	"context"
// 	"log/slog"
// 	"time"
//
// 	"go.mau.fi/whatsmeow"
// 	"go.mau.fi/whatsmeow/types"
// 	"go.mau.fi/whatsmeow/proto/waE2E" // Se mantiene la importación para el mensaje de ayuda
// 	"go.mau.fi/whatsmeow/binary/proto" // Se mantiene la importación para las MessageKey en el pin
//
// 	"api-chatbot/domain"
// )
//
// type PinMessageHandler struct {
// // ... (resto de la estructura y métodos NewPinMessageHandler, Match, Priority no cambian)
// 	userUseCase domain.WhatsAppUserUseCase
// 	client      *whatsmeow.Client
// 	paramCache  domain.ParameterCache
// 	priority    int
// }
//
// func NewPinMessageHandler(
// 	userUseCase domain.WhatsAppUserUseCase,
// 	client *whatsmeow.Client,
// 	paramCache domain.ParameterCache,
// 	priority int,
// ) *PinMessageHandler {
// 	return &PinMessageHandler{
// 		userUseCase: userUseCase,
// 		client:      client,
// 		paramCache:  paramCache,
// 		priority:    priority,
// 	}
// }
//
// func (h *PinMessageHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
// 	if msg.FromMe || msg.IsGroup || msg.ChatID == "status@broadcast" {
// 		return false
// 	}
//
// 	result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
// 	if !result.Success || result.Data == nil {
// 		return false
// 	}
//
// 	user := result.Data
// 	detailsMap, ok := user.Details.(map[string]interface{})
// 	if !ok {
// 		return true
// 	}
//
// 	pinned, exists := detailsMap["help_pinned"].(bool)
// 	return !exists || !pinned
// }
//
// func (h *PinMessageHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
// 	slog.Info("Pinning help message for user",
// 		"whatsapp", msg.From,
// 		"chatID", msg.ChatID,
// 	)
//
// 	helpMessage := h.getParam("MESSAGE_HELP", "👋 *Bienvenido al Asistente del Instituto*\n\nEscribe /help para más información.")
//
// 	jid, err := types.ParseJID(msg.ChatID)
// 	if err != nil {
// 		slog.Error("Failed to parse JID", "error", err, "chatID", msg.ChatID)
// 		return nil
// 	}
//
// 	// Envía el mensaje de ayuda usando la estructura waE2E.Message
// 	sentMsg, err := h.client.SendMessage(ctx, jid, &waE2E.Message{ Conversation: &helpMessage })
// 	if err != nil {
// 		slog.Error("Failed to send help message for pinning", "error", err)
// 		return nil
// 	}
//
// 	time.Sleep(500 * time.Millisecond)
//
// 	// **Llama al método pinMessage refactorizado**
// 	err = h.pinMessage(jid, sentMsg.ID, 30*24*time.Hour)
// 	if err != nil {
// 		slog.Error("Failed to pin message", "error", err, "messageID", sentMsg.ID)
// 		return nil
// 	}
//
// 	slog.Info("Help message pinned successfully",
// 		"whatsapp", msg.From,
// 		"messageID", sentMsg.ID,
// 	)
//
// 	return nil
// }
//
// func (h *PinMessageHandler) Priority() int {
// 	return h.priority
// }
//
// func (h *PinMessageHandler) pinMessage(chatJID types.JID, messageID string, duration time.Duration) error {
// 	var pinDuration uint32
// 	switch duration {
// 	case 24 * time.Hour:
// 		pinDuration = 86400
// 	case 7 * 24 * time.Hour:
// 		pinDuration = 604800
// 	case 30 * time.Hour: // Nota: 30 * 24 * time.Hour es más claro para 30 días
// 		pinDuration = 2592000
// 	default:
// 		pinDuration = 604800
// 	}
//
// 	msgKeyToPin := h.client.BuildMessageKey(chatJID, , types.MessageID(messageID))
//
// 	msgForSend := &waE2E.Message{
// 		PinInChatMessage: &waE2E.PinInChatMessage{
// 			Key:               g.BuildMessageKey(chat, sender, messageId),
// 			Key: msgKeyToPin,
//
// 			// Type: El tipo de pin (fijar para todos).
// 			Type: waE2E.PinInChatMessage_PIN_FOR_ALL.Enum(),
//
// 			// Timestamp de cuándo se realiza la acción
// 			SenderTimestampMS: Int64(time.Now().UnixMilli()),
// 		},
// 	}
//
// 	_, err := h.client.SendMessage(context.Background(), chatJID, msgForSend)
//
// 	return err
// }
//
// func (h *PinMessageHandler) getParam(code, defaultValue string) string {
// 	param, exists := h.paramCache.Get(code)
// 	if !exists {
// 		return defaultValue
// 	}
// 	data, err := param.GetDataAsMap()
// 	if err != nil {
// 		return defaultValue
// 	}
// 	if msg, ok := data["message"].(string); ok {
// 		return msg
// 	}
// 	return defaultValue
// }
