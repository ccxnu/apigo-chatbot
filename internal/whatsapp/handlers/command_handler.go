package handlers

import (
	"context"

	"api-chatbot/domain"
	"api-chatbot/internal/whatsapp"
)

// CommandHandler handles bot commands like /help, /horarios, etc.
type CommandHandler struct {
	whatsapp.BaseHandler
	filter *whatsapp.MessageFilter
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(
	client *whatsapp.Client,
	convUseCase domain.ConversationUseCase,
) *CommandHandler {
	return &CommandHandler{
		BaseHandler: whatsapp.BaseHandler{
			Client:      client,
			ConvUseCase: convUseCase,
		},
		filter: whatsapp.NewMessageFilter(),
	}
}

// Match checks if message is a command
func (h *CommandHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	// Skip own messages
	if h.filter.IsFromMe(msg) {
		return false
	}

	// Check if it's a command
	return len(msg.Body) > 0 && msg.Body[0] == '/'
}

// Handle processes the command
func (h *CommandHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	switch {
	case h.filter.IsCommand(msg, "help"):
		return h.handleHelp(ctx, msg)
	case h.filter.IsCommand(msg, "horarios"):
		return h.handleSchedules(ctx, msg)
	case h.filter.IsCommand(msg, "commands"):
	case h.filter.IsCommand(msg, "comandos"):
		return h.handleCommands(ctx, msg)
	case h.filter.IsCommand(msg, "start"):
		return h.handleStart(ctx, msg)
	default:
		return h.handleUnknownCommand(ctx, msg)
	}

	return nil
}

// Priority - commands have higher priority than RAG
func (h *CommandHandler) Priority() int {
	return 100 // High priority
}

// handleHelp shows help information
func (h *CommandHandler) handleHelp(ctx context.Context, msg *domain.IncomingMessage) error {
	helpText := `👋 *Bienvenido al Asistente del Instituto*

Soy tu asistente virtual y puedo ayudarte con:

🎓 *Información Académica*
   • Programas y carreras
   • Requisitos de admisión
   • Proceso de matrícula
   • Calendario académico

📚 *Consultas Generales*
   Solo escribe tu pregunta y te ayudaré a encontrar la información que necesitas.

⚡ *Comandos Disponibles*
   /help - Muestra esta ayuda
   /horarios - Consulta horarios de clases
   /comandos - Lista todos los comandos

💬 También puedes hacer preguntas directamente, por ejemplo:
   "¿Cuál es el proceso de matrícula?"
   "¿Qué carreras ofrecen?"

¿En qué puedo ayudarte hoy?`

	return h.sendMessage(msg.ChatID, helpText)
}

// handleSchedules shows schedules information
func (h *CommandHandler) handleSchedules(ctx context.Context, msg *domain.IncomingMessage) error {
	scheduleText := `📅 *Consulta de Horarios*

Para consultar horarios, por favor proporciona:
   • Nombre de la carrera o programa
   • Semestre o nivel
   • (Opcional) Materia específica

Ejemplo: "Horario de Ingeniería en Sistemas, tercer semestre"

También puedo ayudarte con horarios de:
   🏫 Horarios de atención administrativa
   📖 Horarios de biblioteca
   🏃 Horarios de actividades extracurriculares

¿Qué horario necesitas consultar?`

	return h.sendMessage(msg.ChatID, scheduleText)
}

// handleCommands lists all available commands
func (h *CommandHandler) handleCommands(ctx context.Context, msg *domain.IncomingMessage) error {
	commandsText := `⚡ *Comandos Disponibles*

/help - Muestra ayuda general del bot
/horarios - Consulta horarios de clases
/comandos - Muestra esta lista de comandos
/start - Reinicia la conversación

💡 *Tip*: No necesitas usar comandos para hacer preguntas. ¡Solo escribe tu consulta!`

	return h.sendMessage(msg.ChatID, commandsText)
}

// handleStart welcomes the user
func (h *CommandHandler) handleStart(ctx context.Context, msg *domain.IncomingMessage) error {
	welcomeText := `👋 ¡Hola! Soy el asistente virtual del Instituto.

Estoy aquí para ayudarte con información sobre:
   • Programas académicos
   • Admisiones y matrículas
   • Horarios y calendarios
   • Y mucho más...

Escribe /help para ver todo lo que puedo hacer, o simplemente hazme una pregunta.

¿En qué puedo ayudarte?`

	return h.sendMessage(msg.ChatID, welcomeText)
}

// handleUnknownCommand responds to unknown commands
func (h *CommandHandler) handleUnknownCommand(ctx context.Context, msg *domain.IncomingMessage) error {
	unknownText := `❓ Comando no reconocido.

Escribe /help para ver los comandos disponibles, o simplemente hazme tu pregunta directamente.`

	return h.sendMessage(msg.ChatID, unknownText)
}

// sendMessage sends a text message
func (h *CommandHandler) sendMessage(chatID, message string) error {
	// TODO: Implement proper message sending
	// For now, placeholder
	return nil
}
