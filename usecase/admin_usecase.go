package usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	d "api-chatbot/domain"
	"api-chatbot/internal/jwttoken"
	"api-chatbot/internal/logger"
)

type adminUseCase struct {
	adminRepo    d.AdminRepository
	tokenService *jwttoken.TokenService
	paramCache   d.ParameterCache
}

// NewAdminUseCase creates a new admin use case
func NewAdminUseCase(
	adminRepo d.AdminRepository,
	tokenService *jwttoken.TokenService,
	paramCache d.ParameterCache,
) d.AdminUseCase {
	return &adminUseCase{
		adminRepo:    adminRepo,
		tokenService: tokenService,
		paramCache:   paramCache,
	}
}

// Login authenticates admin user and returns token pair
func (uc *adminUseCase) Login(ctx context.Context, username, password, ipAddress, userAgent string) d.Result[*d.TokenPairResponse] {
	// Get admin user
	admin, err := uc.adminRepo.GetAdminByUsername(ctx, username)
	if err != nil {
		// Log failed attempt with username
		uc.adminRepo.LogAuthEvent(ctx, nil, username, "login", "failure", &ipAddress, &userAgent, d.Data{
			"reason": "user_not_found",
		})
		logger.LogInfo(ctx, "Admin user not found", "username", username)
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_USER_NOT_FOUND")
	}

	// Check if account is locked
	if admin.IsLocked {
		uc.adminRepo.LogAuthEvent(ctx, &admin.ID, username, "login", "failure", &ipAddress, &userAgent, d.Data{
			"reason": "account_locked",
		})
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_ACCOUNT_LOCKED")
	}

	// Check if account is active
	if !admin.IsActive {
		uc.adminRepo.LogAuthEvent(ctx, &admin.ID, username, "login", "failure", &ipAddress, &userAgent, d.Data{
			"reason": "account_inactive",
		})
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_ACCOUNT_INACTIVE")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		// Increment failed attempts
		uc.adminRepo.IncrementFailedAttempts(ctx, username)
		uc.adminRepo.LogAuthEvent(ctx, &admin.ID, username, "login", "failure", &ipAddress, &userAgent, d.Data{
			"reason": "invalid_password",
		})
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_INVALID_CREDENTIALS")
	}

	// Create token pair
	tokenPair, err := uc.tokenService.CreateTokenPair(jwttoken.TokenMetadata{
		UserID:      admin.ID,
		Username:    admin.Username,
		Email:       admin.Email,
		Name:        admin.Name,
		Role:        admin.Role,
		Permissions: admin.Permissions,
		Claims:      admin.Claims,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	})

	if err != nil {
		logger.LogError(ctx, "Failed to create token pair", err, "username", username)
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	refreshClaims, err := uc.tokenService.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		// Esto no debería pasar si CreateTokenPair fue exitoso, pero es buena práctica.
		logger.LogError(ctx, "Failed to validate/get claims from new refresh token", err, "username", username)
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	// Store refresh token in database
	expiresAt := time.Now().Add(uc.tokenService.GetRefreshTokenExpiry())
	storeResult, err := uc.adminRepo.StoreRefreshToken(ctx, d.StoreRefreshTokenParams{
		AdminID:     admin.ID,
		Token:       tokenPair.RefreshToken,
		TokenFamily: refreshClaims.TokenFamily,
		UserAgent:   &userAgent,
		IPAddress:   &ipAddress,
		ExpiresAt:   expiresAt,
	})

	var tokenID *int
	if err != nil {
		logger.LogWarn(ctx, "Failed to store refresh token", "adminID", admin.ID, "error", err)
		// Continue anyway - token is still valid
	} else if !storeResult.Success {
		logger.LogWarn(ctx, "Failed to store refresh token", "adminID", admin.ID, "code", storeResult.Code)
	} else {
		tokenID = storeResult.TokenID
	}

	// Update last login
	_, err = uc.adminRepo.UpdateAdminLogin(ctx, admin.ID, ipAddress, true)
	if err != nil {
		logger.LogWarn(ctx, "Failed to update admin login", "adminID", admin.ID, "error", err)
	}

	// Log successful login
	uc.adminRepo.LogAuthEvent(ctx, &admin.ID, username, "login", "success", &ipAddress, &userAgent, d.Data{
		"token_family": tokenID,
	})

	response := &d.TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: &d.AdminUserInfo{
			ID:          admin.ID,
			Username:    admin.Username,
			Email:       admin.Email,
			Name:        admin.Name,
			Role:        admin.Role,
			Permissions: admin.Permissions,
			Claims:      admin.Claims,
		},
	}

	return d.Success(response)
}

// RefreshToken generates new token pair using refresh token
func (uc *adminUseCase) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) d.Result[*d.TokenPairResponse] {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.LogInfo(ctx, "Invalid refresh token", "error", err.Error())
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_INVALID_TOKEN")
	}

	// Get refresh token from database
	tokenData, err := uc.adminRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.LogInfo(ctx, "Refresh token not found", "error", err.Error())
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_TOKEN_NOT_FOUND")
	}

	// Check if token is revoked
	if tokenData.IsRevoked {
		// Security breach detected - revoke entire token family
		uc.adminRepo.RevokeTokenFamily(ctx, tokenData.TokenFamily, "security_breach_token_reuse")
		uc.adminRepo.LogAuthEvent(ctx, &tokenData.AdminID, tokenData.Username, "refresh_token", "failure", &ipAddress, &userAgent, d.Data{
			"reason":       "token_reuse_detected",
			"token_family": tokenData.TokenFamily,
		})
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_TOKEN_REVOKED")
	}

	// Check if token is expired
	if time.Now().After(tokenData.ExpiresAt) {
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_TOKEN_EXPIRED")
	}

	// Check if admin is still active
	if !tokenData.AdminActive {
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_ACCOUNT_INACTIVE")
	}

	// Revoke old refresh token
	uc.adminRepo.RevokeRefreshToken(ctx, refreshToken, "token_rotated")

	// Create new token pair with same family
	tokenPair, err := uc.tokenService.CreateTokenPair(jwttoken.TokenMetadata{
		UserID:      tokenData.AdminID,
		Username:    tokenData.Username,
		Email:       tokenData.Email,
		Name:        tokenData.Name,
		Role:        tokenData.Role,
		Permissions: tokenData.Permissions,
		Claims:      tokenData.Claims,
		TokenFamily: tokenData.TokenFamily, // Keep same family for rotation tracking
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	})

	if err != nil {
		logger.LogError(ctx, "Failed to create new token pair", err, "username", claims.Username)
		return d.Error[*d.TokenPairResponse](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	// Store new refresh token
	expiresAt := time.Now().Add(uc.tokenService.GetRefreshTokenExpiry())
	storeResult, err := uc.adminRepo.StoreRefreshToken(ctx, d.StoreRefreshTokenParams{
		AdminID:     tokenData.AdminID,
		Token:       tokenPair.RefreshToken,
		TokenFamily: tokenData.TokenFamily,
		UserAgent:   &userAgent,
		IPAddress:   &ipAddress,
		ExpiresAt:   expiresAt,
	})

	if err != nil {
		logger.LogWarn(ctx, "Failed to store new refresh token", "adminID", tokenData.AdminID, "error", err)
	} else if !storeResult.Success {
		logger.LogWarn(ctx, "Failed to store new refresh token", "adminID", tokenData.AdminID, "code", storeResult.Code)
	}

	// Log successful refresh
	uc.adminRepo.LogAuthEvent(ctx, &tokenData.AdminID, tokenData.Username, "refresh_token", "success", &ipAddress, &userAgent, d.Data{
		"token_family": tokenData.TokenFamily,
	})

	response := &d.TokenPairResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		ExpiresAt:    tokenPair.ExpiresAt,
		User: &d.AdminUserInfo{
			ID:          tokenData.AdminID,
			Username:    tokenData.Username,
			Email:       tokenData.Email,
			Name:        tokenData.Name,
			Role:        tokenData.Role,
			Permissions: tokenData.Permissions,
			Claims:      tokenData.Claims,
		},
	}

	return d.Success(response)
}

// Logout revokes refresh token
func (uc *adminUseCase) Logout(ctx context.Context, refreshToken string) d.Result[d.Data] {
	// Validate token format (don't need to check expiry for logout)
	claims, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		// Even if invalid, try to revoke it
		logger.LogInfo(ctx, "Logout with invalid token", "error", err.Error())
	}

	// Revoke token
	revokeResult, err := uc.adminRepo.RevokeRefreshToken(ctx, refreshToken, "user_logout")
	if err != nil {
		logger.LogWarn(ctx, "Failed to revoke refresh token during logout", "error", err)
		// Continue anyway - logout should succeed
	} else if !revokeResult.Success && revokeResult.Code != "ERR_TOKEN_NOT_FOUND" {
		return d.Error[d.Data](uc.paramCache, revokeResult.Code)
	}

	// Log logout
	if claims != nil {
		uc.adminRepo.LogAuthEvent(ctx, &claims.UserID, claims.Username, "logout", "success", nil, nil, d.Data{})
	}

	return d.Success(d.Data{
		"message": "Logout exitoso",
	})
}

// CreateAdmin creates a new admin user
func (uc *adminUseCase) CreateAdmin(ctx context.Context, params d.CreateAdminUserParams, password string) d.Result[*d.AdminUser] {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError(ctx, "Failed to hash password", err)
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	params.PasswordHash = string(hashedPassword)

	// Create admin
	createResult, err := uc.adminRepo.CreateAdminUser(ctx, params)
	if err != nil {
		logger.LogError(ctx, "Failed to create admin user", err, "username", params.Username)
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	if !createResult.Success {
		return d.Error[*d.AdminUser](uc.paramCache, createResult.Code)
	}

	if createResult.AdminID == nil {
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_CREATE_ADMIN")
	}

	// Get created admin
	admin, err := uc.adminRepo.GetAdminByID(ctx, *createResult.AdminID)
	if err != nil {
		logger.LogError(ctx, "Failed to get created admin", err, "adminID", *createResult.AdminID)
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_INTERNAL_SERVER")
	}

	return d.Success(admin)
}

// GetAdminByUsername retrieves admin by username
func (uc *adminUseCase) GetAdminByUsername(ctx context.Context, username string) d.Result[*d.AdminUser] {
	admin, err := uc.adminRepo.GetAdminByUsername(ctx, username)
	if err != nil {
		logger.LogInfo(ctx, "Admin user not found", "username", username)
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_USER_NOT_FOUND")
	}
	return d.Success(admin)
}

// GetAdminByID retrieves admin by ID
func (uc *adminUseCase) GetAdminByID(ctx context.Context, id int) d.Result[*d.AdminUser] {
	admin, err := uc.adminRepo.GetAdminByID(ctx, id)
	if err != nil {
		logger.LogInfo(ctx, "Admin user not found", "adminID", id)
		return d.Error[*d.AdminUser](uc.paramCache, "ERR_USER_NOT_FOUND")
	}
	return d.Success(admin)
}

// ValidateAccessToken validates and extracts claims from access token
func (uc *adminUseCase) ValidateAccessToken(ctx context.Context, token string) d.Result[*d.TokenClaims] {
	claims, err := uc.tokenService.ValidateAccessToken(token)
	if err != nil {
		logger.LogInfo(ctx, "Invalid access token", "error", err.Error())
		return d.Error[*d.TokenClaims](uc.paramCache, "ERR_INVALID_TOKEN")
	}

	tokenClaims := &d.TokenClaims{
		UserID:      claims.UserID,
		Username:    claims.Username,
		Email:       claims.Email,
		Name:        claims.Name,
		Role:        claims.Role,
		Permissions: claims.Permissions,
		Claims:      claims.Claims,
	}

	return d.Success(tokenClaims)
}
