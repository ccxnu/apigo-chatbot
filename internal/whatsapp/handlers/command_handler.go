package handlers

import (
	"context"
	"strings"

	"api-chatbot/domain"
)

type CommandHandler struct {
	client     WhatsAppClient
	paramCache domain.ParameterCache
	priority   int
}

type WhatsAppClient interface {
	SendText(chatID, message string) error
}

func NewCommandHandler(client WhatsAppClient, paramCache domain.ParameterCache, priority int) *CommandHandler {
	return &CommandHandler{
		client:     client,
		paramCache: paramCache,
		priority:   priority,
	}
}

func (h *CommandHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	if msg.FromMe || msg.IsGroup || msg.ChatID == "status@broadcast" {
		return false
	}
	return len(msg.Body) > 0 && msg.Body[0] == '/'
}

func (h *CommandHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	cmd := strings.ToLower(strings.TrimPrefix(msg.Body, "/"))
	cmd = strings.Fields(cmd)[0]

	switch cmd {
	case "help", "ayuda":
		return h.handleHelp(ctx, msg)
	case "horarios", "schedule":
		return h.handleSchedules(ctx, msg)
	case "commands", "comandos":
		return h.handleCommands(ctx, msg)
	case "start", "inicio":
		return h.handleStart(ctx, msg)
	default:
		return h.handleUnknownCommand(ctx, msg)
	}
}

func (h *CommandHandler) Priority() int {
	return h.priority
}

func (h *CommandHandler) handleHelp(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_HELP", "👋 *Bienvenido al Asistente del Instituto*\n\nEscribe tu pregunta.")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleSchedules(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_SCHEDULES", "📅 *Consulta de Horarios*\n\n¿Qué horario necesitas?")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleCommands(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_COMMANDS", "⚡ *Comandos Disponibles*\n\n/help - Ayuda\n/horarios - Horarios")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleStart(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_START", "👋 ¡Hola! Soy el asistente virtual del Instituto.")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleUnknownCommand(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_UNKNOWN_COMMAND", "❓ Comando no reconocido.")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) sendMessage(chatID, message string) error {
	return h.client.SendText(chatID, message)
}

func (h *CommandHandler) getParam(code, defaultValue string) string {
	param, exists := h.paramCache.Get(code)
	if !exists {
		return defaultValue
	}
	data, err := param.GetDataAsMap()
	if err != nil {
		return defaultValue
	}
	if msg, ok := data["message"].(string); ok {
		return msg
	}
	return defaultValue
}
