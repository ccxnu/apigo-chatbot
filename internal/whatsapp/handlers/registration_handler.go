package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"api-chatbot/domain"
)

// RegistrationHandler handles the OTP-based registration flow
type RegistrationHandler struct {
	regUseCase  domain.RegistrationUseCase
	userUseCase domain.WhatsAppUserUseCase
	convUseCase domain.ConversationUseCase
	client      WhatsAppClient
	paramCache  domain.ParameterCache
	priority    int
}

func NewRegistrationHandler(
	regUseCase domain.RegistrationUseCase,
	userUseCase domain.WhatsAppUserUseCase,
	convUseCase domain.ConversationUseCase,
	client WhatsAppClient,
	paramCache domain.ParameterCache,
	priority int,
) *RegistrationHandler {
	return &RegistrationHandler{
		regUseCase:  regUseCase,
		userUseCase: userUseCase,
		convUseCase: convUseCase,
		client:      client,
		paramCache:  paramCache,
		priority:    priority,
	}
}

func (h *RegistrationHandler) Match(ctx context.Context, msg *domain.IncomingMessage) bool {
	// Skip bot messages, groups, and broadcasts
	if msg.FromMe || msg.IsGroup || msg.ChatID == "status@broadcast" {
		return false
	}

	// Check if user is already registered
	result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
	if result.Success && result.Data != nil {
		return false // User already registered, skip this handler
	}

	// This handler matches for unregistered users
	return true
}

func (h *RegistrationHandler) Handle(ctx context.Context, msg *domain.IncomingMessage) error {
	slog.Info("User not registered, starting OTP registration flow",
		"whatsapp", msg.From,
		"chatID", msg.ChatID,
	)

	// Check if user has a pending registration
	pendingResult := h.regUseCase.GetPendingRegistration(ctx, msg.From)

	if pendingResult.Success && pendingResult.Data != nil {
		pending := pendingResult.Data

		// Check if user is external and hasn't provided email yet
		if pending.UserType == "external" && pending.Email == "" {
			// Extract email from message
			email := h.extractEmail(msg.Body)
			name := h.extractName(msg.Body)

			if email != "" && name != "" {
				// User provided name and email, initiate external registration
				return h.initiateExternalRegistration(ctx, msg, pending.IdentityNumber, name, email)
			}

			// Remind them to provide name and email
			return h.requestExternalUserInfo(msg.ChatID)
		}

		// User has pending registration with email - check if they're sending an OTP code
		otpCode := h.extractOTPCode(msg.Body)
		if otpCode != "" {
			return h.verifyOTPAndRegister(ctx, msg, otpCode)
		}

		// Not an OTP code - check if they sent "reenviar" or "resend"
		if h.isResendRequest(msg.Body) {
			return h.resendOTP(ctx, msg)
		}

		// Remind them to enter OTP
		return h.sendOTPReminder(msg.ChatID, pending)
	}

	// No pending registration - start new registration flow
	cedula := h.extractCedula(msg.Body)

	if cedula == "" {
		return h.requestCedula(msg.ChatID)
	}

	return h.initiateRegistration(ctx, msg, cedula)
}

func (h *RegistrationHandler) Priority() int {
	return h.priority
}

// Extract cedula (10 digits)
func (h *RegistrationHandler) extractCedula(text string) string {
	re := regexp.MustCompile(`\b\d{10}\b`)
	match := re.FindString(text)
	return match
}

// Extract OTP code (6 digits)
func (h *RegistrationHandler) extractOTPCode(text string) string {
	re := regexp.MustCompile(`\b\d{6}\b`)
	match := re.FindString(text)
	return match
}

// Check if message is a resend request
func (h *RegistrationHandler) isResendRequest(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	return lower == "reenviar" || lower == "resend" || lower == "nuevo código" || lower == "nuevo codigo"
}

// Request cedula from user
func (h *RegistrationHandler) requestCedula(chatID string) error {
	message := h.getParam("MESSAGE_REQUEST_CEDULA", `👋 ¡Hola! Bienvenido al asistente virtual del Instituto.

Para poder ayudarte, necesito que te registres primero.

Por favor, envíame tu número de cédula (10 dígitos).

Ejemplo: 1234567890`)

	return h.client.SendText(chatID, message)
}

// Initiate registration with cedula
func (h *RegistrationHandler) initiateRegistration(ctx context.Context, msg *domain.IncomingMessage, cedula string) error {
	slog.Info("Initiating registration with cedula", "cedula", cedula, "whatsapp", msg.From)

	// Start registration process (validates with AcademicOK and creates pending registration)
	result := h.regUseCase.InitiateRegistration(ctx, msg.From, cedula)

	if !result.Success {
		slog.Error("Registration initiation failed",
			"cedula", cedula,
			"whatsapp", msg.From,
			"code", result.Code,
		)

		// Handle specific error codes
		if result.Code == "ERR_USER_ALREADY_EXISTS" {
			return h.client.SendText(msg.ChatID,
				"✅ Ya estás registrado en el sistema. Puedes empezar a chatear conmigo.")
		}

		if result.Code == "ERR_IDENTITY_ALREADY_REGISTERED" {
			return h.client.SendText(msg.ChatID,
				"❌ Esta cédula ya está registrada con otro número de WhatsApp. Si crees que esto es un error, contacta al administrador.")
		}

		if result.Code == "ERR_EXTERNAL_USER_INFO_REQUIRED" {
			return h.handleExternalUser(ctx, msg.From, cedula, msg.ChatID)
		}

		if result.Code == "ERR_IDENTITY_NOT_FOUND" {
			return h.client.SendText(msg.ChatID,
				"❌ No pude validar tu cédula. Por favor verifica que sea correcta e intenta nuevamente.")
		}

		return h.client.SendText(msg.ChatID,
			"❌ Ocurrió un error al iniciar tu registro. Por favor intenta nuevamente.")
	}

	pending := result.Data

	// Send success message with OTP instructions
	message := fmt.Sprintf(`✅ ¡Hola %s!

He enviado un código de verificación de 6 dígitos a tu correo electrónico:
📧 %s

Por favor, revisa tu bandeja de entrada (y también la carpeta de spam) y envíame el código para completar tu registro.

El código expirará en 10 minutos.

Si no recibes el correo, escribe "reenviar" para generar un nuevo código.`,
		pending.Name,
		maskEmail(pending.Email))

	return h.client.SendText(msg.ChatID, message)
}

// Verify OTP and complete registration
func (h *RegistrationHandler) verifyOTPAndRegister(ctx context.Context, msg *domain.IncomingMessage, otpCode string) error {
	slog.Info("Verifying OTP code", "whatsapp", msg.From, "otpLength", len(otpCode))

	result := h.regUseCase.VerifyAndRegister(ctx, msg.From, otpCode)

	if !result.Success {
		slog.Warn("OTP verification failed",
			"whatsapp", msg.From,
			"code", result.Code,
		)

		// Handle specific error codes
		if result.Code == "ERR_INVALID_OTP" {
			message, _ := h.paramCache.Get("ERR_INVALID_OTP")
			if message != nil {
				if data, err := message.GetDataAsMap(); err == nil {
					if msg, ok := data["message"].(string); ok {
						return h.client.SendText(msg.ChatID, msg)
					}
				}
			}
			return h.client.SendText(msg.ChatID,
				"❌ Código incorrecto. Por favor verifica e intenta nuevamente.\n\nSi no tienes el código, escribe 'reenviar'.")
		}

		if result.Code == "ERR_OTP_EXPIRED" {
			return h.client.SendText(msg.ChatID,
				"⏰ El código ha expirado. Escribe 'reenviar' para generar un nuevo código.")
		}

		if result.Code == "ERR_MAX_ATTEMPTS" {
			return h.client.SendText(msg.ChatID,
				"🚫 Has excedido el número máximo de intentos. Escribe 'reenviar' para generar un nuevo código.")
		}

		if result.Code == "ERR_NO_PENDING_REG" {
			return h.client.SendText(msg.ChatID,
				"❌ No tienes un registro pendiente. Por favor envía tu cédula para iniciar el registro.")
		}

		return h.client.SendText(msg.ChatID,
			"❌ Ocurrió un error al verificar tu código. Por favor intenta nuevamente.")
	}

	user := result.Data

	// Create conversation
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

	// Send welcome message
	welcomeMessage := h.buildWelcomeMessage(user)
	err := h.client.SendText(msg.ChatID, welcomeMessage)
	if err != nil {
		return err
	}

	// Send help message
	helpMessage := h.getParam("MESSAGE_HELP", "👋 *Bienvenido al Asistente del Instituto*\n\nEscribe /help para más información.")

	slog.Info("User registered successfully via OTP",
		"whatsapp", msg.From,
		"name", user.Name,
		"role", user.Role,
	)

	return h.client.SendText(msg.ChatID, helpMessage)
}

// Resend OTP code
func (h *RegistrationHandler) resendOTP(ctx context.Context, msg *domain.IncomingMessage) error {
	slog.Info("Resending OTP code", "whatsapp", msg.From)

	result := h.regUseCase.ResendOTP(ctx, msg.From)

	if !result.Success {
		slog.Error("OTP resend failed",
			"whatsapp", msg.From,
			"code", result.Code,
		)

		if result.Code == "ERR_NO_PENDING_REGISTRATION" {
			return h.client.SendText(msg.ChatID,
				"❌ No tienes un registro pendiente. Por favor envía tu cédula para iniciar el registro.")
		}

		return h.client.SendText(msg.ChatID,
			"❌ Ocurrió un error al reenviar el código. Por favor intenta nuevamente.")
	}

	pending := result.Data

	message := fmt.Sprintf(`✅ He enviado un nuevo código de verificación a tu correo:
📧 %s

El código anterior ya no es válido. Por favor envíame el nuevo código de 6 dígitos.

El código expirará en 10 minutos.`,
		maskEmail(pending.Email))

	return h.client.SendText(msg.ChatID, message)
}

// Send reminder to enter OTP
func (h *RegistrationHandler) sendOTPReminder(chatID string, pending *domain.PendingRegistration) error {
	message := fmt.Sprintf(`⏳ Estás en proceso de registro.

He enviado un código de verificación a tu correo:
📧 %s

Por favor, envíame el código de 6 dígitos que recibiste.

Si no lo has recibido, escribe "reenviar" para generar un nuevo código.`,
		maskEmail(pending.Email))

	return h.client.SendText(chatID, message)
}

// Handle external users (not in AcademicOK) - create pending registration
func (h *RegistrationHandler) handleExternalUser(ctx context.Context, whatsapp, cedula, chatID string) error {
	slog.Info("External user detected, creating pending registration",
		"whatsapp", whatsapp,
		"cedula", cedula,
	)

	// Create pending registration without email (will be collected next)
	result := h.regUseCase.InitiateExternalRegistration(ctx, whatsapp, cedula)

	if !result.Success {
		slog.Error("Failed to create external pending registration",
			"code", result.Code,
		)
		return h.client.SendText(chatID,
			"❌ Ocurrió un error al iniciar tu registro. Por favor intenta nuevamente.")
	}

	// Ask for name and email
	return h.requestExternalUserInfo(chatID)
}

// Request name and email from external user
func (h *RegistrationHandler) requestExternalUserInfo(chatID string) error {
	message := `👤 No encontré tu cédula en nuestra base de datos institucional.

Sin embargo, puedes registrarte como usuario externo.

Por favor envíame tu información en el siguiente formato:

*Nombre Completo / correo@email.com*

Ejemplo:
Juan Pérez / juan.perez@gmail.com

Recibirás un código de verificación en ese correo para completar tu registro.`

	return h.client.SendText(chatID, message)
}

// Initiate external user registration with provided email
func (h *RegistrationHandler) initiateExternalRegistration(ctx context.Context, msg *domain.IncomingMessage, cedula, name, email string) error {
	slog.Info("Initiating external registration with email",
		"whatsapp", msg.From,
		"name", name,
		"email", email,
	)

	// Complete the external registration with name and email
	result := h.regUseCase.CompleteExternalRegistration(ctx, msg.From, name, email)

	if !result.Success {
		slog.Error("Failed to complete external registration",
			"code", result.Code,
		)
		return h.client.SendText(msg.ChatID,
			"❌ Ocurrió un error al registrar tu información. Por favor intenta nuevamente.")
	}

	pending := result.Data

	// Send success message with OTP instructions
	message := fmt.Sprintf(`✅ ¡Hola %s!

He enviado un código de verificación de 6 dígitos a tu correo electrónico:
📧 %s

Por favor, revisa tu bandeja de entrada (y también la carpeta de spam) y envíame el código para completar tu registro.

El código expirará en 10 minutos.

Si no recibes el correo, escribe "reenviar" para generar un nuevo código.`,
		pending.Name,
		maskEmail(pending.Email))

	return h.client.SendText(msg.ChatID, message)
}

// Extract email from message (format: "Name / email@domain.com" or just "email@domain.com")
func (h *RegistrationHandler) extractEmail(text string) string {
	// Email regex
	re := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	match := re.FindString(text)
	return strings.TrimSpace(match)
}

// Extract name from message (format: "Name / email@domain.com")
func (h *RegistrationHandler) extractName(text string) string {
	// Split by "/" to get name part
	parts := strings.Split(text, "/")
	if len(parts) >= 2 {
		// Name is before the "/"
		name := strings.TrimSpace(parts[0])
		// Validate name has at least 2 words (first and last name)
		words := strings.Fields(name)
		if len(words) >= 2 {
			return name
		}
	}
	return ""
}

// Build welcome message based on user role
func (h *RegistrationHandler) buildWelcomeMessage(user *domain.WhatsAppUser) string {
	var roleEmoji string
	var roleText string

	switch user.Role {
	case "ROLE_STUDENT":
		roleEmoji = "🎓"
		roleText = "estudiante"
	case "ROLE_PROFESSOR":
		roleEmoji = "👨‍🏫"
		roleText = "docente"
	case "ROLE_EXTERNAL":
		roleEmoji = "👤"
		roleText = "usuario externo"
	default:
		roleEmoji = "👤"
		roleText = "usuario"
	}

	return fmt.Sprintf(`%s ¡Registro completado, %s!

Has sido registrado exitosamente como %s.

Ahora puedes hacer preguntas sobre el instituto y recibir asistencia.`, roleEmoji, user.Name, roleText)
}

// Get parameter value with fallback
func (h *RegistrationHandler) getParam(code, defaultValue string) string {
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

// Mask email for privacy (show only first 2 chars and domain)
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) <= 2 {
		return fmt.Sprintf("%s***@%s", localPart, domain)
	}

	return fmt.Sprintf("%s***@%s", localPart[:2], domain)
}
