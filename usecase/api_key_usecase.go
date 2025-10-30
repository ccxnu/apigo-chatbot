package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type apiKeyUseCase struct {
	apiKeyRepo d.APIKeyRepository
	cache      d.ParameterCache
	timeout    time.Duration
}

func NewAPIKeyUseCase(
	apiKeyRepo d.APIKeyRepository,
	cache d.ParameterCache,
	timeout time.Duration,
) d.APIKeyUseCase {
	return &apiKeyUseCase{
		apiKeyRepo: apiKeyRepo,
		cache:      cache,
		timeout:    timeout,
	}
}

// GenerateAPIKey generates a secure random API key with the format: sk_live_xxxxx
func GenerateAPIKey() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64 and add prefix
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return fmt.Sprintf("sk_live_%s", encoded), nil
}

// HashAPIKey hashes an API key using bcrypt (similar to password hashing)
func HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyAPIKey verifies an API key against its hash
func VerifyAPIKey(apiKey, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(apiKey))
}

func (u *apiKeyUseCase) CreateAPIKey(ctx context.Context, params d.CreateAPIKeyParams) d.Result[*d.APIKey] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Creating new API key",
		"operation", "CreateAPIKey",
		"name", params.Name,
		"type", params.Type,
		"rateLimit", params.RateLimit,
	)

	// Generate API key value if not provided
	if params.Value == "" {
		generatedKey, err := GenerateAPIKey()
		if err != nil {
			logger.LogError(c, "Failed to generate API key", err,
				"operation", "CreateAPIKey",
			)
			return d.Error[*d.APIKey](u.cache, "ERR_GENERATE_API_KEY")
		}
		params.Value = generatedKey
	}

	// Store the original key value to return to the user
	originalKey := params.Value

	// Hash the API key before storing (security best practice)
	hashedKey, err := HashAPIKey(params.Value)
	if err != nil {
		logger.LogError(c, "Failed to hash API key", err,
			"operation", "CreateAPIKey",
		)
		return d.Error[*d.APIKey](u.cache, "ERR_HASH_API_KEY")
	}
	params.Value = hashedKey

	// Create in database
	result, err := u.apiKeyRepo.Create(c, params)
	if err != nil || result == nil {
		logger.LogError(c, "Failed to create API key in database", err,
			"operation", "CreateAPIKey",
			"name", params.Name,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(c, "API key creation failed with business logic error",
			"operation", "CreateAPIKey",
			"code", result.Code,
			"name", params.Name,
		)
		return d.Error[*d.APIKey](u.cache, result.Code)
	}

	// Retrieve the created API key
	apiKey, err := u.apiKeyRepo.GetByID(c, *result.KeyID)
	if err != nil || apiKey == nil {
		logger.LogError(c, "Failed to retrieve created API key", err,
			"operation", "CreateAPIKey",
			"keyID", *result.KeyID,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_INTERNAL_DB")
	}

	// Replace hashed value with original for response (only time the key is shown in plain text)
	apiKey.Value = originalKey

	logger.LogInfo(c, "API key created successfully",
		"operation", "CreateAPIKey",
		"keyID", apiKey.ID,
		"name", apiKey.Name,
	)

	return d.Success(apiKey)
}

func (u *apiKeyUseCase) ValidateAPIKey(ctx context.Context, keyValue, ipAddress, endpoint string) d.Result[*d.APIKey] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Validating API key",
		"operation", "ValidateAPIKey",
		"endpoint", endpoint,
		"ipAddress", ipAddress,
	)

	// Note: Since we hash API keys, we can't do a direct lookup by value
	// We need to retrieve all active keys and check each hash
	// For better performance, consider caching active API keys
	apiKeys, err := u.apiKeyRepo.GetAll(c)
	if err != nil {
		logger.LogError(c, "Failed to retrieve API keys for validation", err,
			"operation", "ValidateAPIKey",
		)
		return d.Error[*d.APIKey](u.cache, "ERR_INTERNAL_DB")
	}

	var matchedKey *d.APIKey
	for i := range apiKeys {
		if err := VerifyAPIKey(keyValue, apiKeys[i].Value); err == nil {
			matchedKey = &apiKeys[i]
			break
		}
	}

	if matchedKey == nil {
		logger.LogWarn(c, "Invalid API key",
			"operation", "ValidateAPIKey",
			"endpoint", endpoint,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_INVALID_API_KEY")
	}

	// Check if key is active
	if !matchedKey.IsActive {
		logger.LogWarn(c, "API key is not active",
			"operation", "ValidateAPIKey",
			"keyID", matchedKey.ID,
			"name", matchedKey.Name,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_API_KEY_INACTIVE")
	}

	// Check if key has expired
	if matchedKey.ExpiresAt != nil && time.Now().After(*matchedKey.ExpiresAt) {
		logger.LogWarn(c, "API key has expired",
			"operation", "ValidateAPIKey",
			"keyID", matchedKey.ID,
			"name", matchedKey.Name,
			"expiresAt", *matchedKey.ExpiresAt,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_API_KEY_EXPIRED")
	}

	// Check IP whitelist (if configured)
	if len(matchedKey.AllowedIPs) > 0 {
		ipAllowed := false
		for _, allowedIP := range matchedKey.AllowedIPs {
			if allowedIP == ipAddress {
				ipAllowed = true
				break
			}
		}

		if !ipAllowed {
			logger.LogWarn(c, "IP address not allowed for API key",
				"operation", "ValidateAPIKey",
				"keyID", matchedKey.ID,
				"name", matchedKey.Name,
				"ipAddress", ipAddress,
			)
			return d.Error[*d.APIKey](u.cache, "ERR_IP_NOT_ALLOWED")
		}
	}

	// Check permissions (if configured)
	if len(matchedKey.Permissions) > 0 {
		permissionAllowed := false
		for _, permission := range matchedKey.Permissions {
			// Simple prefix matching for permissions
			// e.g., permission "/v1/chat" allows "/v1/chat/completions"
			if len(endpoint) >= len(permission) && endpoint[:len(permission)] == permission {
				permissionAllowed = true
				break
			}
		}

		if !permissionAllowed {
			logger.LogWarn(c, "Endpoint not allowed for API key",
				"operation", "ValidateAPIKey",
				"keyID", matchedKey.ID,
				"name", matchedKey.Name,
				"endpoint", endpoint,
			)
			return d.Error[*d.APIKey](u.cache, "ERR_ENDPOINT_NOT_ALLOWED")
		}
	}

	// Update last used timestamp asynchronously
	go func() {
		asyncCtx, asyncCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer asyncCancel()
		_, _ = u.apiKeyRepo.UpdateLastUsed(asyncCtx, matchedKey.ID)
	}()

	logger.LogInfo(c, "API key validated successfully",
		"operation", "ValidateAPIKey",
		"keyID", matchedKey.ID,
		"name", matchedKey.Name,
	)

	return d.Success(matchedKey)
}

func (u *apiKeyUseCase) UpdateAPIKey(ctx context.Context, params d.UpdateAPIKeyParams) d.Result[d.Data] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Updating API key",
		"operation", "UpdateAPIKey",
		"keyID", params.KeyID,
	)

	result, err := u.apiKeyRepo.Update(c, params)
	if err != nil || result == nil {
		logger.LogError(c, "Failed to update API key in database", err,
			"operation", "UpdateAPIKey",
			"keyID", params.KeyID,
		)
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(c, "API key update failed with business logic error",
			"operation", "UpdateAPIKey",
			"code", result.Code,
			"keyID", params.KeyID,
		)
		return d.Error[d.Data](u.cache, result.Code)
	}

	logger.LogInfo(c, "API key updated successfully",
		"operation", "UpdateAPIKey",
		"keyID", params.KeyID,
	)

	return d.Success(d.Data{})
}

func (u *apiKeyUseCase) RevokeAPIKey(ctx context.Context, keyID int) d.Result[d.Data] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Revoking API key",
		"operation", "RevokeAPIKey",
		"keyID", keyID,
	)

	result, err := u.apiKeyRepo.Delete(c, keyID)
	if err != nil || result == nil {
		logger.LogError(c, "Failed to revoke API key in database", err,
			"operation", "RevokeAPIKey",
			"keyID", keyID,
		)
		return d.Error[d.Data](u.cache, "ERR_INTERNAL_DB")
	}

	if !result.Success {
		logger.LogWarn(c, "API key revocation failed with business logic error",
			"operation", "RevokeAPIKey",
			"code", result.Code,
			"keyID", keyID,
		)
		return d.Error[d.Data](u.cache, result.Code)
	}

	logger.LogInfo(c, "API key revoked successfully",
		"operation", "RevokeAPIKey",
		"keyID", keyID,
	)

	return d.Success(d.Data{})
}

func (u *apiKeyUseCase) ListAPIKeys(ctx context.Context) d.Result[[]d.APIKey] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Listing all API keys",
		"operation", "ListAPIKeys",
	)

	apiKeys, err := u.apiKeyRepo.GetAll(c)
	if err != nil {
		logger.LogError(c, "Failed to retrieve API keys from database", err,
			"operation", "ListAPIKeys",
		)
		return d.Error[[]d.APIKey](u.cache, "ERR_INTERNAL_DB")
	}

	logger.LogInfo(c, "API keys retrieved successfully",
		"operation", "ListAPIKeys",
		"count", len(apiKeys),
	)

	return d.Success(apiKeys)
}

func (u *apiKeyUseCase) GetAPIKeyByID(ctx context.Context, keyID int) d.Result[*d.APIKey] {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	logger.LogInfo(c, "Getting API key by ID",
		"operation", "GetAPIKeyByID",
		"keyID", keyID,
	)

	apiKey, err := u.apiKeyRepo.GetByID(c, keyID)
	if err != nil {
		logger.LogError(c, "Failed to retrieve API key from database", err,
			"operation", "GetAPIKeyByID",
			"keyID", keyID,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_INTERNAL_DB")
	}

	if apiKey == nil {
		logger.LogWarn(c, "API key not found",
			"operation", "GetAPIKeyByID",
			"keyID", keyID,
		)
		return d.Error[*d.APIKey](u.cache, "ERR_API_KEY_NOT_FOUND")
	}

	logger.LogInfo(c, "API key retrieved successfully",
		"operation", "GetAPIKeyByID",
		"keyID", keyID,
		"name", apiKey.Name,
	)

	return d.Success(apiKey)
}
