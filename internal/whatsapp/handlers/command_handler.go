package handlers

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow/types"

	"api-chatbot/domain"
)

type CommandHandler struct {
	client      WhatsAppClient
	paramCache  domain.ParameterCache
	regUseCase  domain.RegistrationUseCase
	userUseCase domain.WhatsAppUserUseCase
	convUseCase domain.ConversationUseCase
	priority    int
}

type WhatsAppClient interface {
	SendText(chatID, message string) error
	SendChatPresence(chatID string, state types.ChatPresence, media types.ChatPresenceMedia) error
}

func NewCommandHandler(
	client WhatsAppClient,
	paramCache domain.ParameterCache,
	regUseCase domain.RegistrationUseCase,
	userUseCase domain.WhatsAppUserUseCase,
	convUseCase domain.ConversationUseCase,
	priority int,
) *CommandHandler {
	return &CommandHandler{
		client:      client,
		paramCache:  paramCache,
		regUseCase:  regUseCase,
		userUseCase: userUseCase,
		convUseCase: convUseCase,
		priority:    priority,
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
	case "ayuda", "help":
		return h.handleHelp(ctx, msg)
	case "horarios", "schedule":
		return h.handleSchedules(ctx, msg)
	case "comandos", "commands":
		return h.handleCommands(ctx, msg)
	case "inicio", "start":
		return h.handleStart(ctx, msg)
	case "registrar", "registro", "register":
		return h.handleRegister(ctx, msg)
	case "cancelar", "reset":
		return h.handleReset(ctx, msg)
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
	message := h.getParam("MESSAGE_COMMANDS", `⚡ *Comandos Disponibles*

/ayuda - Ayuda general del bot
/inicio - Mensaje de bienvenida
/horarios - Consulta horarios disponibles
/registrar - Registrarse en el sistema
/cancelar - Cancelar registro en curso
/comandos - Ver esta lista de comandos

💡 _También puedes escribir tus preguntas directamente sin usar comandos._`)
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleStart(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_START", "👋 ¡Hola! Soy el asistente virtual del Instituto.")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleUnknownCommand(ctx context.Context, msg *domain.IncomingMessage) error {
	message := h.getParam("MESSAGE_UNKNOWN_COMMAND", "❓ Comando no reconocido.\n\nUsa /comandos para ver los comandos disponibles.")
	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) sendMessage(chatID, message string) error {
	return h.client.SendText(chatID, message)
}

func (h *CommandHandler) handleRegister(ctx context.Context, msg *domain.IncomingMessage) error {
	// Check if user is already registered
	result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
	if result.Success && result.Data != nil {
		return h.sendMessage(msg.ChatID,
			"✅ Ya estás registrado en el sistema.\n\nPuedes usar /ayuda para ver lo que puedo hacer por ti.")
	}

	// Check if user already has a pending registration
	pendingResult := h.regUseCase.GetPendingRegistration(ctx, msg.From)
	if pendingResult.Success && pendingResult.Data != nil {
		pending := pendingResult.Data
		// User has pending registration - inform them about current step
		var stepMessage string
		switch pending.RegistrationStep {
		case "STEP_INIT", "STEP_REQUEST_CEDULA":
			stepMessage = "Estás en proceso de registro. Por favor envía tu cédula de 10 dígitos."
		case "STEP_SELECT_USER_TYPE":
			stepMessage = "Estás en proceso de registro. Por favor selecciona tu tipo:\n\n*1* - 🎓 Estudiante\n*2* - 👨‍🏫 Docente\n*3* - 👤 Usuario externo"
		case "STEP_REQUEST_EMAIL_NAME":
			stepMessage = "Estás en proceso de registro. Por favor envía tu información:\n\n*Nombre Completo / correo@email.com*"
		case "STEP_VERIFY_OTP":
			stepMessage = "Estás en proceso de registro. Por favor ingresa el código de verificación que te enviamos por correo.\n\nSi no lo recibiste, escribe 'reenviar'."
		default:
			stepMessage = "Ya tienes un registro en curso."
		}

		return h.sendMessage(msg.ChatID,
			"ℹ️ "+stepMessage+"\n\nSi quieres cancelar y empezar de nuevo, usa /cancelar")
	}

	// Create a pending registration with STEP_REQUEST_CEDULA
	// so RegistrationHandler can match subsequent messages
	createResult := h.regUseCase.InitiatePendingForCedula(ctx, msg.From)
	if !createResult.Success {
		return h.sendMessage(msg.ChatID,
			"❌ Ocurrió un error al iniciar el registro. Por favor intenta nuevamente.")
	}

	// Show the cedula request message
	message := h.getParam("MESSAGE_REQUEST_CEDULA", `👋 ¡Hola! Vamos a registrarte en el sistema.

Para comenzar, envíame tu número de cédula (10 dígitos).

Ejemplo: 1234567890`)

	return h.sendMessage(msg.ChatID, message)
}

func (h *CommandHandler) handleReset(ctx context.Context, msg *domain.IncomingMessage) error {
	// Check if user has a pending registration
	pendingResult := h.regUseCase.GetPendingRegistration(ctx, msg.From)
	if !pendingResult.Success || pendingResult.Data == nil {
		return h.sendMessage(msg.ChatID,
			"ℹ️ No tienes un registro en curso.\n\nUsa /registrar para iniciar el registro.")
	}

	// Cancel the pending registration
	cancelResult := h.regUseCase.CancelPendingRegistration(ctx, msg.From)
	if !cancelResult.Success {
		return h.sendMessage(msg.ChatID,
			"❌ Ocurrió un error al cancelar tu registro. Por favor intenta nuevamente.")
	}

	return h.sendMessage(msg.ChatID,
		`✅ Tu registro ha sido cancelado.

Ahora puedes:
• Usar /registrar para iniciar un nuevo registro
• Chatear conmigo si eres usuario externo (límite de 10 mensajes por día)
• Usar /ayuda para ver la ayuda`)
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
