package usecase

import (
	"context"
	"fmt"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type whatsappUserUseCase struct {
	repo           d.WhatsAppUserRepository
	httpClient     d.HTTPClient
	paramCache     d.ParameterCache
	contextTimeout time.Duration
}

func NewWhatsAppUserUseCase(
	repo d.WhatsAppUserRepository,
	httpClient d.HTTPClient,
	paramCache d.ParameterCache,
	timeout time.Duration,
) d.WhatsAppUserUseCase {
	return &whatsappUserUseCase{
		repo:           repo,
		httpClient:     httpClient,
		paramCache:     paramCache,
		contextTimeout: timeout,
	}
}

// GetOrRegisterUser retrieves a user by WhatsApp or registers them after validation
func (uc *whatsappUserUseCase) GetOrRegisterUser(c context.Context, whatsapp string, identityNumber string) d.Result[*d.WhatsAppUser] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	// First, try to get user by WhatsApp
	user, err := uc.repo.GetByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to get user by WhatsApp from database", err,
			"operation", "GetOrRegisterUser",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if user != nil {
		return d.Success(user)
	}

	// User doesn't exist, validate with institute API
	validationResult := uc.ValidateWithInstituteAPI(ctx, identityNumber)
	if !validationResult.Success {
		return d.Error[*d.WhatsAppUser](uc.paramCache, validationResult.Code)
	}

	instituteData := validationResult.Data
	if !instituteData.IsValid {
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INVALID_IDENTITY")
	}

	// Create new user
	createParams := d.CreateUserParams{
		IdentityNumber: instituteData.IdentityNumber,
		Name:           instituteData.Name,
		Email:          instituteData.Email,
		Phone:          instituteData.Phone,
		Role:           instituteData.Role,
		WhatsApp:       whatsapp,
		Details:        d.Data{},
	}

	createResult, err := uc.repo.Create(ctx, createParams)
	if err != nil || createResult == nil {
		logger.LogError(ctx, "Failed to create user in database", err,
			"operation", "GetOrRegisterUser",
			"identityNumber", identityNumber,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if !createResult.Success {
		logger.LogWarn(ctx, "User creation failed with business logic error",
			"operation", "GetOrRegisterUser",
			"code", createResult.Code,
			"identityNumber", identityNumber,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, createResult.Code)
	}

	// Retrieve the newly created user
	user, err = uc.repo.GetByWhatsApp(ctx, whatsapp)
	if err != nil || user == nil {
		logger.LogError(ctx, "Failed to retrieve newly created user from database", err,
			"operation", "GetOrRegisterUser",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	return d.Success(user)
}

// ValidateWithInstituteAPI validates a user's identity with the AcademicOK API
func (uc *whatsappUserUseCase) ValidateWithInstituteAPI(c context.Context, identityNumber string) d.Result[*d.InstituteUserData] {
	ctx, cancel := context.WithTimeout(c, 300*uc.contextTimeout)
	defer cancel()

	// Get AcademicOK API configuration from parameters
	param, exists := uc.paramCache.Get("ACADEMICOK_CONFIGURATION")
	if !exists {
		logger.LogWarn(ctx, "ACADEMICOK_CONFIGURATION parameter not found - using mock data",
			"operation", "ValidateWithInstituteAPI",
			"identityNumber", identityNumber,
		)
		// For development, return mock data
		return uc.getMockInstituteData(identityNumber)
	}

	config, err := param.GetDataAsMap()
	if err != nil {
		logger.LogError(ctx, "Failed to parse ACADEMICOK_CONFIGURATION", err,
			"operation", "ValidateWithInstituteAPI",
		)
		return uc.getMockInstituteData(identityNumber)
	}

	personaURL, _ := config["personaURL"].(string)
	docenteURL, _ := config["docenteURL"].(string)
	personaKey, _ := config["personaKey"].(string)
	docenteKey, _ := config["docenteKey"].(string)

	if personaURL == "" || docenteURL == "" {
		logger.LogWarn(ctx, "AcademicOK API URLs not configured - using mock data",
			"operation", "ValidateWithInstituteAPI",
		)
		return uc.getMockInstituteData(identityNumber)
	}

	// Step 1: Call apidatospersona first
	personaData, err := uc.callAcademicOKPersonaAPI(ctx, personaURL, identityNumber, personaKey)
	if err != nil {
		logger.LogError(ctx, "Failed to call apidatospersona API", err,
			"operation", "ValidateWithInstituteAPI",
			"identityNumber", identityNumber,
		)
		return d.Error[*d.InstituteUserData](uc.paramCache, "ERR_IDENTITY_NOT_FOUND")
	}

	// Step 2: Check if user is a student (has careras with length > 0)
	if personaData != nil && len(personaData.Data.Careras) > 0 {
		// It's a student
		return d.Success(&d.InstituteUserData{
			IdentityNumber: personaData.Data.Cedula,
			Name:           personaData.Data.Nombre,
			Email:          personaData.Data.Email,
			Phone:          "",
			Role:           "ROLE_STUDENT",
			IsValid:        true,
		})
	}

	// Step 3: If has name but no careras, check if it's a professor
	if personaData != nil && personaData.Data.Nombre != "" {
		docenteData, err := uc.callAcademicOKDocenteAPI(ctx, docenteURL, identityNumber, docenteKey)
		if err == nil && docenteData != nil && docenteData.Data.Nombre != "" {
			// It's a professor
			return d.Success(&d.InstituteUserData{
				IdentityNumber: docenteData.Data.Cedula,
				Name:           docenteData.Data.Nombre,
				Email:          docenteData.Data.Email,
				Phone:          "",
				Role:           "ROLE_PROFESSOR",
				IsValid:        true,
			})
		}
	}

	// Step 4: Not found in either API - it's an external user
	return d.Error[*d.InstituteUserData](uc.paramCache, "ERR_EXTERNAL_USER_INFO_REQUIRED")
}

// AcademicOKPersonaResponse represents the response from apidatospersona
type AcademicOKPersonaResponse struct {
	Data struct {
		Nombre    string `json:"nombre"`
		Cedula    string `json:"cedula"`
		Pasaporte string `json:"pasaporte"`
		Email     string `json:"email"`
		Foto      string `json:"foto"`
		Careras   []struct {
			Estado   bool   `json:"estado"`
			Egresado bool   `json:"egresado"`
			Graduado bool   `json:"graduado"`
			Carrera  string `json:"carrera"`
			Sesion   string `json:"sesion"`
			Periodo  string `json:"periodo"`
		} `json:"careras"`
	} `json:"data"`
	Result string `json:"result"`
}

// AcademicOKDocenteResponse represents the response from apidatosdocente
type AcademicOKDocenteResponse struct {
	Data struct {
		Nombre    string `json:"nombre"`
		Cedula    string `json:"cedula"`
		Pasaporte string `json:"pasaporte"`
		Email     string `json:"email"`
		Foto      string `json:"foto"`
	} `json:"data"`
	Result string `json:"result"`
}

// callAcademicOKPersonaAPI calls the apidatospersona endpoint
func (uc *whatsappUserUseCase) callAcademicOKPersonaAPI(ctx context.Context, baseURL, identityNumber, apiKey string) (*AcademicOKPersonaResponse, error) {
	url := fmt.Sprintf("%s&key=%s&identificacion=%s", baseURL, apiKey, identityNumber)

	var response AcademicOKPersonaResponse

	req := d.HTTPRequest{
		URL:    url,
		Method: "GET",
	}

	err := uc.httpClient.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != "ok" {
		return nil, fmt.Errorf("API returned error result: %s", response.Result)
	}

	return &response, nil
}

// callAcademicOKDocenteAPI calls the apidatosdocente endpoint
func (uc *whatsappUserUseCase) callAcademicOKDocenteAPI(ctx context.Context, baseURL, identityNumber, apiKey string) (*AcademicOKDocenteResponse, error) {
	url := fmt.Sprintf("%s&key=%s&identificacion=%s", baseURL, apiKey, identityNumber)

	var response AcademicOKDocenteResponse

	req := d.HTTPRequest{
		URL:    url,
		Method: "GET",
	}

	err := uc.httpClient.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	if response.Result != "ok" {
		return nil, fmt.Errorf("API returned error result: %s", response.Result)
	}

	return &response, nil
}

// getMockInstituteData returns mock data for development
func (uc *whatsappUserUseCase) getMockInstituteData(identityNumber string) d.Result[*d.InstituteUserData] {
	// For development: return mock data
	mockData := &d.InstituteUserData{
		IdentityNumber: identityNumber,
		Name:           "Test User",
		Email:          fmt.Sprintf("%s@institute.edu", identityNumber),
		Phone:          "1234567890",
		Role:           "ROLE_STUDENT",
		IsValid:        true,
	}

	return d.Success(mockData)
}

// GetUserByWhatsApp retrieves a user by WhatsApp number
func (uc *whatsappUserUseCase) GetUserByWhatsApp(c context.Context, whatsapp string) d.Result[*d.WhatsAppUser] {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	user, err := uc.repo.GetByWhatsApp(ctx, whatsapp)
	if err != nil {
		logger.LogError(ctx, "Failed to get user by WhatsApp from database", err,
			"operation", "GetUserByWhatsApp",
			"whatsapp", whatsapp,
		)
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_INTERNAL_DB")
	}

	if user == nil {
		return d.Error[*d.WhatsAppUser](uc.paramCache, "ERR_USER_NOT_FOUND")
	}

	return d.Success(user)
}
