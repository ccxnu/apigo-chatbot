package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"api-chatbot/domain"
)

// UserValidationHandler ensures user is registered before processing messages
// This has the HIGHEST priority - it runs before any other handler
type UserValidationHandler struct {
	userUseCase domain.WhatsAppUserUseCase
	convUseCase domain.ConversationUseCase
	client      WhatsAppClient
	priority    int
}

// NewUserValidationHandler creates a new user validation handler
func NewUserValidationHandler(
	userUseCase domain.WhatsAppUserUseCase,
	convUseCase domain.ConversationUseCase,
	client WhatsAppClient,
	priority int,
) *UserValidationHandler {
	return &UserValidationHandler{
		userUseCase: userUseCase,
		convUseCase: convUseCase,
		client:      client,
		priority:    priority,
	}
}

// Match - Only match if user is NOT registered (to handle registration)
func (h *UserValidationHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	// Skip own messages
	if msg.FromMe {
		return false
	}

	// ONLY respond to personal/direct messages - skip groups
	if msg.IsGroup {
		return false
	}

	// Skip WhatsApp status broadcasts
	if msg.ChatID == "status@broadcast" {
		return false
	}

	// Check if user exists
	result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)

	// Match only if user does NOT exist (needs registration)
	return !result.Success || result.Data == nil
}

// Handle - Process user registration flow
func (h *UserValidationHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	slog.Info("User not registered, starting registration flow",
		"whatsapp", msg.From,
		"chatID", msg.ChatID,
	)

	// Check if message contains a cédula (10 digits)
	cedula := h.extractCedula(msg.Body)

	if cedula == "" {
		// Ask for cédula
		return h.requestCedula(msg.ChatID)
	}

	// Validate with AcademicOK API
	return h.registerUser(ctx, msg, cedula)
}

// Priority - HIGHEST priority (runs first)
func (h *UserValidationHandler) Priority() int {
	return h.priority
}

// extractCedula extracts a 10-digit cédula from message
func (h *UserValidationHandler) extractCedula(text string) string {
	// Match 10 consecutive digits
	re := regexp.MustCompile(`\b\d{10}\b`)
	match := re.FindString(text)
	return match
}

// requestCedula asks user to provide their cédula
func (h *UserValidationHandler) requestCedula(chatID string) error {
	message := `👋 ¡Hola! Bienvenido al asistente virtual del Instituto.

Para poder ayudarte, necesito que te registres primero.

Por favor, envíame tu número de cédula (10 dígitos).

Ejemplo: 1234567890`

	return h.client.SendText(chatID, message)
}

// registerUser validates and registers the user
func (h *UserValidationHandler) registerUser(ctx context.Context, msg *domain.IncomingMessage, cedula string) error {
	slog.Info("Validating user with AcademicOK API", "cedula", cedula)

	// Validate with institute API
	validationResult := h.userUseCase.ValidateWithInstituteAPI(ctx, cedula)

	if !validationResult.Success {
		// Check if it's an external user
		if validationResult.Code == "ERR_EXTERNAL_USER_INFO_REQUIRED" {
			return h.handleExternalUser(msg.ChatID)
		}

		// Other validation error
		slog.Error("Failed to validate user",
			"cedula", cedula,
			"code", validationResult.Code,
		)
		return h.client.SendText(msg.ChatID,
			"❌ No pude validar tu cédula. Por favor verifica que sea correcta e intenta nuevamente.")
	}

	instituteData := validationResult.Data

	// Register user
	registrationResult := h.userUseCase.GetOrRegisterUser(ctx, msg.From, cedula)

	if !registrationResult.Success {
		slog.Error("Failed to register user",
			"cedula", cedula,
			"whatsapp", msg.From,
			"code", registrationResult.Code,
		)
		return h.client.SendText(msg.ChatID,
			"❌ Ocurrió un error al registrarte. Por favor intenta nuevamente.")
	}

	user := registrationResult.Data

	// Create or get conversation
	contactName := user.Name
	var groupName *string
	if msg.GroupName != "" {
		groupName = &msg.GroupName
	}

	convResult := h.convUseCase.GetOrCreateConversation(
		ctx,
		msg.ChatID,
		msg.From,
		&contactName,
		msg.IsGroup,
		groupName,
	)
	if !convResult.Success {
		slog.Error("Failed to create conversation",
			"chatID", msg.ChatID,
			"code", convResult.Code,
		)
		// Don't fail registration, just log
	}

	// Welcome message based on role
	welcomeMessage := h.buildWelcomeMessage(instituteData, user)

	slog.Info("User registered successfully",
		"cedula", cedula,
		"name", user.Name,
		"role", user.Role,
		"whatsapp", msg.From,
	)

	return h.client.SendText(msg.ChatID, welcomeMessage)
}

// handleExternalUser handles users not found in institute database
func (h *UserValidationHandler) handleExternalUser(chatID string) error {
	message := `👤 No encontré tu cédula en nuestra base de datos.

Si eres un visitante externo, por favor proporciona:
1. Tu nombre completo
2. Tu correo electrónico

Ejemplo:
Juan Pérez
juan.perez@email.com

O si eres estudiante/docente, verifica que tu cédula sea correcta.`

	return h.client.SendText(chatID, message)
}

// buildWelcomeMessage creates a personalized welcome message
func (h *UserValidationHandler) buildWelcomeMessage(instituteData *domain.InstituteUserData, user *domain.WhatsAppUser) string {
	var roleEmoji string
	var roleText string

	switch user.Role {
	case "ROLE_STUDENT":
		roleEmoji = "🎓"
		roleText = "estudiante"
	case "ROLE_PROFESSOR":
		roleEmoji = "👨‍🏫"
		roleText = "docente"
	default:
		roleEmoji = "👤"
		roleText = "usuario"
	}

	return fmt.Sprintf(`%s ¡Bienvenido, %s!

Has sido registrado exitosamente como %s.

Ahora puedes:
• Hacer preguntas sobre el instituto
• Consultar horarios con /horarios
• Ver ayuda con /help

¿En qué puedo ayudarte hoy?`, roleEmoji, user.Name, roleText)
}
