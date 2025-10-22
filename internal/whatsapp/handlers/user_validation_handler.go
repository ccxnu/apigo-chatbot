package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"api-chatbot/domain"
)

type UserValidationHandler struct {
	userUseCase domain.WhatsAppUserUseCase
	convUseCase domain.ConversationUseCase
	client      WhatsAppClient
	paramCache  domain.ParameterCache
	priority    int
}

func NewUserValidationHandler(
	userUseCase domain.WhatsAppUserUseCase,
	convUseCase domain.ConversationUseCase,
	client WhatsAppClient,
	paramCache domain.ParameterCache,
	priority int,
) *UserValidationHandler {
	return &UserValidationHandler{
		userUseCase: userUseCase,
		convUseCase: convUseCase,
		client:      client,
		paramCache:  paramCache,
		priority:    priority,
	}
}

func (h *UserValidationHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	if msg.FromMe || msg.IsGroup || msg.ChatID == "status@broadcast" {
		return false
	}

	result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
	return !result.Success || result.Data == nil
}

func (h *UserValidationHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	slog.Info("User not registered, starting registration flow",
		"whatsapp", msg.From,
		"chatID", msg.ChatID,
	)

	cedula := h.extractCedula(msg.Body)

	if cedula == "" {
		return h.requestCedula(msg.ChatID)
	}

	return h.registerUser(ctx, msg, cedula)
}

func (h *UserValidationHandler) Priority() int {
	return h.priority
}

func (h *UserValidationHandler) extractCedula(text string) string {
	re := regexp.MustCompile(`\b\d{10}\b`)
	match := re.FindString(text)
	return match
}

func (h *UserValidationHandler) requestCedula(chatID string) error {
	message := `👋 ¡Hola! Bienvenido al asistente virtual del Instituto.

Para poder ayudarte, necesito que te registres primero.

Por favor, envíame tu número de cédula (10 dígitos).

Ejemplo: 1234567890`

	return h.client.SendText(chatID, message)
}

func (h *UserValidationHandler) registerUser(ctx context.Context, msg *domain.IncomingMessage, cedula string) error {
	slog.Info("Validating user with AcademicOK API", "cedula", cedula)

	validationResult := h.userUseCase.ValidateWithInstituteAPI(ctx, cedula)

	if !validationResult.Success {
		if validationResult.Code == "ERR_EXTERNAL_USER_INFO_REQUIRED" {
			return h.handleExternalUser(msg.ChatID)
		}

		slog.Error("Failed to validate user",
			"cedula", cedula,
			"code", validationResult.Code,
		)
		return h.client.SendText(msg.ChatID,
			"❌ No pude validar tu cédula. Por favor verifica que sea correcta e intenta nuevamente.")
	}

	instituteData := validationResult.Data

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
	}

	welcomeMessage := h.buildWelcomeMessage(instituteData, user)

	err := h.client.SendText(msg.ChatID, welcomeMessage)
	if err != nil {
		return err
	}

	helpMessage := h.getParam("MESSAGE_HELP", "👋 *Bienvenido al Asistente del Instituto*\n\nEscribe /help para más información.")

	slog.Info("User registered successfully",
		"cedula", cedula,
		"name", user.Name,
		"role", user.Role,
		"whatsapp", msg.From,
	)

	return h.client.SendText(msg.ChatID, helpMessage)
}

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

Has sido registrado exitosamente como %s.`, roleEmoji, user.Name, roleText)
}

func (h *UserValidationHandler) getParam(code, defaultValue string) string {
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
