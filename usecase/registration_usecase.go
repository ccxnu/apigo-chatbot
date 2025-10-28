package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type registrationUseCase struct {
	regRepo        d.RegistrationRepository
	userRepo       d.WhatsAppUserRepository
	userUseCase    d.WhatsAppUserUseCase
	mailer         d.OTPMailer
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewRegistrationUseCase(
	regRepo d.RegistrationRepository,
	userRepo d.WhatsAppUserRepository,
	userUseCase d.WhatsAppUserUseCase,
	mailer d.OTPMailer,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.RegistrationUseCase {
	return &registrationUseCase{
		regRepo:        regRepo,
		userRepo:       userRepo,
		userUseCase:    userUseCase,
		mailer:         mailer,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// InitiateRegistration starts the registration process with OTP
func (uc *registrationUseCase) InitiateRegistration(
	c context.Context,
	whatsapp, identityNumber string,
) d.Result[*d.PendingRegistration] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to check existing user", err,
			"operation", "InitiateRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if existingUser != nil {
		logger.LogWarn(ctx, "User already registered",
			"operation", "InitiateRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_USER_ALREADY_EXISTS")
	}

	// Check if cedula is already registered with different WhatsApp
	existingByCedula, err := uc.userRepo.GetByIdentity(ctx, identityNumber)
	if err != nil {
		logger.LogError(ctx, "Failed to check existing user by cedula", err,
			"operation", "InitiateRegistration",
			"identityNumber", identityNumber,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if existingByCedula != nil && existingByCedula.WhatsApp != whatsapp {
		logger.LogWarn(ctx, "Identity number already registered with different WhatsApp",
			"operation", "InitiateRegistration",
			"identityNumber", identityNumber,
			"existingWhatsApp", existingByCedula.WhatsApp,
			"newWhatsApp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_IDENTITY_ALREADY_REGISTERED")
	}

	// Validate with institute API
	validationResult := uc.userUseCase.ValidateWithInstituteAPI(ctx, identityNumber)

	var name, email, phone, role, userType string
	var details d.Data

	if validationResult.Success {
		// Institute user (student or professor)
		instituteData := validationResult.Data
		name = instituteData.Name
		email = instituteData.Email
		phone = instituteData.Phone
		role = instituteData.Role
		userType = "institute"
		details = d.Data{}
	} else if validationResult.Code == "ERR_EXTERNAL_USER_INFO_REQUIRED" {
		// External user - will need to provide name and email via WhatsApp later
		// For now, create pending with minimal info
		userType = "external"
		role = "ROLE_EXTERNAL"
		details = d.Data{}

		logger.LogInfo(ctx, "External user registration initiated",
			"operation", "InitiateRegistration",
			"identityNumber", identityNumber,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_EXTERNAL_USER_INFO_REQUIRED")
	} else {
		// Validation failed for other reason
		logger.LogError(ctx, "Identity validation failed", nil,
			"operation", "InitiateRegistration",
			"code", validationResult.Code,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, validationResult.Code)
	}

	// Generate OTP code
	otpCode, err := uc.generateOTP()
	if err != nil {
		logger.LogError(ctx, "Failed to generate OTP", err,
			"operation", "InitiateRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL")
	}

	// Get OTP expiration duration from parameters (default 10 minutes)
	otpExpirationMinutes := 10
	if param, exists := uc.paramCache.Get("OTP_EXPIRATION_MINUTES"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if minutes, ok := data["minutes"].(float64); ok {
				otpExpirationMinutes = int(minutes)
			}
		}
	}

	otpExpiresAt := time.Now().Add(time.Duration(otpExpirationMinutes) * time.Minute)

	// Create pending registration
	createParams := d.CreatePendingRegistrationParams{
		IdentityNumber: identityNumber,
		WhatsApp:       whatsapp,
		Name:           name,
		Email:          email,
		Phone:          phone,
		Role:           role,
		UserType:       userType,
		Details:        details,
		OTPCode:        otpCode,
		OTPExpiresAt:   otpExpiresAt,
	}

	createResult, err := uc.regRepo.CreatePendingRegistration(ctx, createParams)
	if err != nil {
		logger.LogError(ctx, "Failed to create pending registration", err,
			"operation", "InitiateRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !createResult.Success {
		logger.LogWarn(ctx, "Pending registration creation failed with business logic error",
			"operation", "InitiateRegistration",
			"code", createResult.Code,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, createResult.Code)
	}

	// Send OTP email
	err = uc.mailer.SendOTPEmail(ctx, email, name, otpCode, userType)
	if err != nil {
		logger.LogError(ctx, "Failed to send OTP email", err,
			"operation", "InitiateRegistration",
			"email", email,
		)
		// Don't fail the registration, but log the error
		// User can still request a new OTP
	}

	// Retrieve the pending registration
	pending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil || pending == nil {
		logger.LogError(ctx, "Failed to retrieve pending registration", err,
			"operation", "InitiateRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(ctx, "Registration initiated successfully",
		"operation", "InitiateRegistration",
		"whatsapp", whatsapp,
		"email", email,
		"userType", userType,
	)

	return d.Success(pending)
}

// VerifyAndRegister verifies OTP and completes user registration
func (uc *registrationUseCase) VerifyAndRegister(
	c context.Context,
	whatsapp, otpCode string,
) d.Result[*d.WhatsAppUser] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Verify OTP
	verifyParams := d.VerifyOTPParams{
		WhatsApp:  whatsapp,
		OTPCode:   otpCode,
		IPAddress: nil, // TODO: Extract from context if available
	}

	verifyResult, err := uc.regRepo.VerifyOTP(ctx, verifyParams)
	if err != nil {
		logger.LogError(ctx, "Failed to verify OTP", err,
			"operation", "VerifyAndRegister",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !verifyResult.Success {
		logger.LogWarn(ctx, "OTP verification failed",
			"operation", "VerifyAndRegister",
			"whatsapp", whatsapp,
			"code", verifyResult.Code,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, verifyResult.Code)
	}

	// Create the user
	createUserParams := d.CreateUserParams{
		IdentityNumber: verifyResult.IdentityNumber,
		Name:           verifyResult.Name,
		Email:          verifyResult.Email,
		Phone:          verifyResult.Phone,
		Role:           verifyResult.Role,
		WhatsApp:       whatsapp,
		Details:        verifyResult.Details,
	}

	createResult, err := uc.userRepo.Create(ctx, createUserParams)
	if err != nil {
		logger.LogError(ctx, "Failed to create user", err,
			"operation", "VerifyAndRegister",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !createResult.Success {
		logger.LogWarn(ctx, "User creation failed with business logic error",
			"operation", "VerifyAndRegister",
			"code", createResult.Code,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, createResult.Code)
	}

	// Delete pending registration
	if verifyResult.PendingID != nil {
		err = uc.regRepo.DeletePendingRegistration(ctx, *verifyResult.PendingID)
		if err != nil {
			logger.LogWarn(ctx, "Failed to delete pending registration after user creation",
				"operation", "VerifyAndRegister",
				"pendingID", *verifyResult.PendingID,
				"error", err,
			)
			// Don't fail - user is already created
		}
	}

	// Retrieve the newly created user
	user, err := uc.userRepo.GetByWhatsApp(ctx, whatsapp)
	if err != nil || user == nil {
		logger.LogError(ctx, "Failed to retrieve newly created user", err,
			"operation", "VerifyAndRegister",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(ctx, "User registered successfully",
		"operation", "VerifyAndRegister",
		"whatsapp", whatsapp,
		"name", user.Name,
		"role", user.Role,
	)

	return d.Success(user)
}

// GetPendingRegistration retrieves pending registration by WhatsApp
func (uc *registrationUseCase) GetPendingRegistration(
	c context.Context,
	whatsapp string,
) d.Result[*d.PendingRegistration] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	pending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to get pending registration", err,
			"operation", "GetPendingRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if pending == nil {
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_NO_PENDING_REGISTRATION")
	}

	return d.Success(pending)
}

// ResendOTP generates and sends a new OTP code
func (uc *registrationUseCase) ResendOTP(
	c context.Context,
	whatsapp string,
) d.Result[*d.PendingRegistration] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Get existing pending registration
	pending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to get pending registration for resend", err,
			"operation", "ResendOTP",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if pending == nil {
		logger.LogWarn(ctx, "No pending registration found for OTP resend",
			"operation", "ResendOTP",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_NO_PENDING_REGISTRATION")
	}

	// Generate new OTP
	newOTPCode, err := uc.generateOTP()
	if err != nil {
		logger.LogError(ctx, "Failed to generate new OTP", err,
			"operation", "ResendOTP",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL")
	}

	// Get OTP expiration duration
	otpExpirationMinutes := 10
	if param, exists := uc.paramCache.Get("OTP_EXPIRATION_MINUTES"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if minutes, ok := data["minutes"].(float64); ok {
				otpExpirationMinutes = int(minutes)
			}
		}
	}

	newExpiresAt := time.Now().Add(time.Duration(otpExpirationMinutes) * time.Minute)

	// Update pending registration with new OTP
	updateParams := d.CreatePendingRegistrationParams{
		IdentityNumber: pending.IdentityNumber,
		WhatsApp:       whatsapp,
		Name:           pending.Name,
		Email:          pending.Email,
		Phone:          pending.Phone,
		Role:           pending.Role,
		UserType:       pending.UserType,
		Details:        pending.Details,
		OTPCode:        newOTPCode,
		OTPExpiresAt:   newExpiresAt,
	}

	updateResult, err := uc.regRepo.CreatePendingRegistration(ctx, updateParams)
	if err != nil {
		logger.LogError(ctx, "Failed to update pending registration with new OTP", err,
			"operation", "ResendOTP",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !updateResult.Success {
		logger.LogWarn(ctx, "OTP resend failed with business logic error",
			"operation", "ResendOTP",
			"code", updateResult.Code,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, updateResult.Code)
	}

	// Send new OTP email
	err = uc.mailer.SendOTPEmail(ctx, pending.Email, pending.Name, newOTPCode, pending.UserType)
	if err != nil {
		logger.LogError(ctx, "Failed to send new OTP email", err,
			"operation", "ResendOTP",
			"email", pending.Email,
		)
		// Don't fail - OTP is already updated in DB
	}

	// Retrieve updated pending registration
	updatedPending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil || updatedPending == nil {
		logger.LogError(ctx, "Failed to retrieve updated pending registration", err,
			"operation", "ResendOTP",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(ctx, "OTP resent successfully",
		"operation", "ResendOTP",
		"whatsapp", whatsapp,
		"email", pending.Email,
	)

	return d.Success(updatedPending)
}

// InitiateExternalRegistration creates pending registration for external user (without email)
func (uc *registrationUseCase) InitiateExternalRegistration(
	c context.Context,
	whatsapp, identityNumber string,
) d.Result[*d.PendingRegistration] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to check existing user", err,
			"operation", "InitiateExternalRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if existingUser != nil {
		logger.LogWarn(ctx, "User already registered",
			"operation", "InitiateExternalRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_USER_ALREADY_EXISTS")
	}

	// Create pending registration WITHOUT email (will be collected in next step)
	// No OTP yet - will be generated after email is provided
	createParams := d.CreatePendingRegistrationParams{
		IdentityNumber: identityNumber,
		WhatsApp:       whatsapp,
		Name:           "", // Will be provided by user
		Email:          "", // Will be provided by user
		Phone:          "",
		Role:           "ROLE_EXTERNAL",
		UserType:       "external",
		Details:        d.Data{},
		OTPCode:        "", // No OTP yet
		OTPExpiresAt:   time.Time{}, // No expiration yet
	}

	createResult, err := uc.regRepo.CreatePendingRegistration(ctx, createParams)
	if err != nil {
		logger.LogError(ctx, "Failed to create external pending registration", err,
			"operation", "InitiateExternalRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !createResult.Success {
		logger.LogWarn(ctx, "External pending registration creation failed",
			"operation", "InitiateExternalRegistration",
			"code", createResult.Code,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, createResult.Code)
	}

	// Retrieve the pending registration
	pending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil || pending == nil {
		logger.LogError(ctx, "Failed to retrieve external pending registration", err,
			"operation", "InitiateExternalRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(ctx, "External pending registration created",
		"operation", "InitiateExternalRegistration",
		"whatsapp", whatsapp,
	)

	return d.Success(pending)
}

// CompleteExternalRegistration updates pending registration with name/email and generates OTP
func (uc *registrationUseCase) CompleteExternalRegistration(
	c context.Context,
	whatsapp, name, email string,
) d.Result[*d.PendingRegistration] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// Get existing pending registration
	pending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to get pending registration", err,
			"operation", "CompleteExternalRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if pending == nil {
		logger.LogWarn(ctx, "No pending registration found",
			"operation", "CompleteExternalRegistration",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_NO_PENDING_REGISTRATION")
	}

	// Generate OTP code
	otpCode, err := uc.generateOTP()
	if err != nil {
		logger.LogError(ctx, "Failed to generate OTP", err,
			"operation", "CompleteExternalRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL")
	}

	// Get OTP expiration duration
	otpExpirationMinutes := 10
	if param, exists := uc.paramCache.Get("OTP_EXPIRATION_MINUTES"); exists {
		if data, err := param.GetDataAsMap(); err == nil {
			if minutes, ok := data["minutes"].(float64); ok {
				otpExpirationMinutes = int(minutes)
			}
		}
	}

	otpExpiresAt := time.Now().Add(time.Duration(otpExpirationMinutes) * time.Minute)

	// Update pending registration with name, email, and OTP
	updateParams := d.CreatePendingRegistrationParams{
		IdentityNumber: pending.IdentityNumber,
		WhatsApp:       whatsapp,
		Name:           name,
		Email:          email,
		Phone:          "",
		Role:           "ROLE_EXTERNAL",
		UserType:       "external",
		Details:        d.Data{},
		OTPCode:        otpCode,
		OTPExpiresAt:   otpExpiresAt,
	}

	updateResult, err := uc.regRepo.CreatePendingRegistration(ctx, updateParams)
	if err != nil {
		logger.LogError(ctx, "Failed to update external pending registration", err,
			"operation", "CompleteExternalRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !updateResult.Success {
		logger.LogWarn(ctx, "External registration completion failed",
			"operation", "CompleteExternalRegistration",
			"code", updateResult.Code,
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, updateResult.Code)
	}

	// Send OTP email
	err = uc.mailer.SendOTPEmail(ctx, email, name, otpCode, "external")
	if err != nil {
		logger.LogError(ctx, "Failed to send OTP email to external user", err,
			"operation", "CompleteExternalRegistration",
			"email", email,
		)
		// Don't fail - OTP is already saved in DB
	}

	// Retrieve updated pending registration
	updatedPending, err := uc.regRepo.GetPendingByWhatsApp(ctx, whatsapp)
	if err != nil || updatedPending == nil {
		logger.LogError(ctx, "Failed to retrieve updated external pending registration", err,
			"operation", "CompleteExternalRegistration",
		)
		return d.Error[*d.PendingRegistration](uc.paramCache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(ctx, "External registration completed with OTP sent",
		"operation", "CompleteExternalRegistration",
		"whatsapp", whatsapp,
		"email", email,
	)

	return d.Success(updatedPending)
}

// generateOTP generates a 6-digit OTP code
func (uc *registrationUseCase) generateOTP() (string, error) {
	// Generate a random 6-digit number
	max := big.NewInt(1000000) // 0-999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format as 6-digit string with leading zeros
	return fmt.Sprintf("%06d", n.Int64()), nil
}
