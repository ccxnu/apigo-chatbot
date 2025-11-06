package domain

import (
	"time"

	"api-chatbot/api/dal"
	"context"
)

// AdminUser represents an administrative user
type AdminUser struct {
	ID             int        `json:"id" db:"adm_id"`
	Username       string     `json:"username" db:"adm_username"`
	Email          string     `json:"email" db:"adm_email"`
	PasswordHash   string     `json:"-" db:"adm_password_hash"` // Never expose in JSON
	Name           string     `json:"name" db:"adm_name"`
	Role           string     `json:"role" db:"adm_role"`
	Permissions    []string   `json:"permissions" db:"adm_permissions"`
	Claims         Data       `json:"claims" db:"adm_claims"`
	IsActive       bool       `json:"isActive" db:"adm_is_active"`
	IsLocked       bool       `json:"isLocked" db:"adm_is_locked"`
	FailedAttempts int        `json:"failedAttempts" db:"adm_failed_attempts"`
	LastLogin      *time.Time `json:"lastLogin,omitempty" db:"adm_last_login"`
	LastLoginIP    *string    `json:"lastLoginIp,omitempty" db:"adm_last_login_ip"`
	CreatedAt      time.Time  `json:"createdAt" db:"adm_created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"adm_updated_at"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID            int        `json:"id" db:"rft_id"`
	AdminID       int        `json:"adminId" db:"rft_admin_id"`
	Token         string     `json:"token" db:"rft_token"`
	TokenFamily   string     `json:"tokenFamily" db:"rft_token_family"`
	UserAgent     *string    `json:"userAgent,omitempty" db:"rft_user_agent"`
	IPAddress     *string    `json:"ipAddress,omitempty" db:"rft_ip_address"`
	ExpiresAt     time.Time  `json:"expiresAt" db:"rft_expires_at"`
	IsRevoked     bool       `json:"isRevoked" db:"rft_is_revoked"`
	RevokedAt     *time.Time `json:"revokedAt,omitempty" db:"rft_revoked_at"`
	RevokedReason *string    `json:"revokedReason,omitempty" db:"rft_revoked_reason"`
	CreatedAt     time.Time  `json:"createdAt" db:"rft_created_at"`
}

// RefreshTokenWithAdmin combines refresh token with admin user info
type RefreshTokenWithAdmin struct {
	RefreshToken
	Username    string   `json:"username" db:"adm_username"`
	Email       string   `json:"email" db:"adm_email"`
	Name        string   `json:"name" db:"adm_name"`
	Role        string   `json:"role" db:"adm_role"`
	Permissions []string `json:"permissions" db:"adm_permissions"`
	Claims      Data     `json:"claims" db:"adm_claims"`
	AdminActive bool     `json:"adminActive" db:"adm_is_active"`
}

// APIKey represents an API key for external integrations
type APIKey struct {
	ID          int        `json:"id" db:"key_id"`
	Name        string     `json:"name" db:"key_name"`
	Value       string     `json:"value" db:"key_value"`
	Type        string     `json:"type" db:"key_type"`
	Claims      Data       `json:"claims" db:"key_claims"`
	RateLimit   int        `json:"rateLimit" db:"key_rate_limit"`
	AllowedIPs  []string   `json:"allowedIps" db:"key_allowed_ips"`
	Permissions []string   `json:"permissions" db:"key_permissions"`
	IsActive    bool       `json:"isActive" db:"key_is_active"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" db:"key_expires_at"`
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty" db:"key_last_used_at"`
	CreatedBy   *int       `json:"createdBy,omitempty" db:"key_created_by"`
	CreatedAt   time.Time  `json:"createdAt" db:"key_created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"key_updated_at"`
}

// AuthLog represents an authentication event log
type AuthLog struct {
	ID        int       `json:"id" db:"log_id"`
	UserID    *int      `json:"userId,omitempty" db:"log_user_id"`
	Username  string    `json:"username" db:"log_username"`
	Action    string    `json:"action" db:"log_action"`
	Status    string    `json:"status" db:"log_status"`
	IPAddress *string   `json:"ipAddress,omitempty" db:"log_ip_address"`
	UserAgent *string   `json:"userAgent,omitempty" db:"log_user_agent"`
	Details   Data      `json:"details" db:"log_details"`
	CreatedAt time.Time `json:"createdAt" db:"log_created_at"`
}

// CreateAdminUserParams parameters for creating an admin user
type CreateAdminUserParams struct {
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"-"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Permissions  []string `json:"permissions,omitempty"`
	Claims       Data     `json:"claims,omitempty"`
}

// CreateAdminUserResult result from creating admin user
type CreateAdminUserResult struct {
	dal.DbResult
	AdminID *int `db:"admin_id"`
}

// StoreRefreshTokenParams parameters for storing refresh token
type StoreRefreshTokenParams struct {
	AdminID     int
	Token       string
	TokenFamily string
	UserAgent   *string
	IPAddress   *string
	ExpiresAt   time.Time
}

// StoreRefreshTokenResult result from storing refresh token
type StoreRefreshTokenResult struct {
	dal.DbResult
	TokenID *int `db:"token_id"`
}

// RevokeTokenResult result from revoking token
type RevokeTokenResult struct {
	dal.DbResult
}

// RevokeTokenFamilyResult result from revoking token family
type RevokeTokenFamilyResult struct {
	dal.DbResult
	RevokedCount *int `db:"revoked_count"`
}

// UpdateLoginResult result from updating login info
type UpdateLoginResult struct {
	dal.DbResult
}

// IncrementFailedAttemptsResult result from incrementing failed attempts
type IncrementFailedAttemptsResult struct {
	dal.DbResult
	IsLocked *bool `db:"is_locked"`
}

// LogAuthEventResult result from logging auth event
type LogAuthEventResult struct {
	dal.DbResult
}

// TokenPairResponse represents the response with token pair
type TokenPairResponse struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	TokenType    string         `json:"tokenType"`
	ExpiresIn    int            `json:"expiresIn"`
	ExpiresAt    time.Time      `json:"expiresAt"`
	User         *AdminUserInfo `json:"user"`
}

// AdminUserInfo represents safe admin user info for responses
type AdminUserInfo struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	Claims      Data     `json:"claims,omitempty"`
}

// TokenClaims represents validated token claims for use in middleware
type TokenClaims struct {
	UserID      int      `json:"userId"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	Claims      Data     `json:"claims"`
}

type CleanupResult struct {
	dal.DbResult
	DeletedCount int `db:"deleted_count"`
}

type AdminRepository interface {
	// User operations
	CreateAdminUser(ctx context.Context, params CreateAdminUserParams) (*CreateAdminUserResult, error)
	GetAdminByUsername(ctx context.Context, username string) (*AdminUser, error)
	GetAdminByID(ctx context.Context, id int) (*AdminUser, error)
	UpdateAdminLogin(ctx context.Context, adminID int, ipAddress string, resetFailedAttempts bool) (*UpdateLoginResult, error)
	IncrementFailedAttempts(ctx context.Context, username string) (*IncrementFailedAttemptsResult, error)

	// Refresh token operations
	StoreRefreshToken(ctx context.Context, params StoreRefreshTokenParams) (*StoreRefreshTokenResult, error)
	GetRefreshToken(ctx context.Context, token string) (*RefreshTokenWithAdmin, error)
	RevokeRefreshToken(ctx context.Context, token string, reason string) (*RevokeTokenResult, error)
	RevokeTokenFamily(ctx context.Context, tokenFamily string, reason string) (*RevokeTokenFamilyResult, error)
	CleanupExpiredTokens(ctx context.Context) (*CleanupResult, error)

	// Auth logging
	LogAuthEvent(ctx context.Context, userID *int, username, action, status string, ipAddress, userAgent *string, details Data) (*LogAuthEventResult, error)
}

// AdminUseCase defines business logic for admin authentication
type AdminUseCase interface {
	// Authentication
	Login(ctx context.Context, username, password, ipAddress, userAgent string) Result[*TokenPairResponse]
	RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) Result[*TokenPairResponse]
	Logout(ctx context.Context, refreshToken string) Result[Data]

	// User management
	CreateAdmin(ctx context.Context, params CreateAdminUserParams, password string) Result[*AdminUser]
	GetAdminByUsername(ctx context.Context, username string) Result[*AdminUser]
	GetAdminByID(ctx context.Context, id int) Result[*AdminUser]

	// Token validation
	ValidateAccessToken(ctx context.Context, token string) Result[*TokenClaims]
}

// APIUsage represents a tracked API usage record
type APIUsage struct {
	ID            int       `json:"id" db:"usg_id"`
	APIKeyID      int       `json:"apiKeyId" db:"usg_api_key_id"`
	Endpoint      string    `json:"endpoint" db:"usg_endpoint"`
	Method        string    `json:"method" db:"usg_method"`
	StatusCode    int       `json:"statusCode" db:"usg_status_code"`
	TokensUsed    int       `json:"tokensUsed" db:"usg_tokens_used"`
	RequestTimeMs int       `json:"requestTimeMs" db:"usg_request_time_ms"`
	IPAddress     *string   `json:"ipAddress,omitempty" db:"usg_ip_address"`
	UserAgent     *string   `json:"userAgent,omitempty" db:"usg_user_agent"`
	ErrorMessage  *string   `json:"errorMessage,omitempty" db:"usg_error_message"`
	CreatedAt     time.Time `json:"createdAt" db:"usg_created_at"`
}

// APIUsageStats represents aggregated usage statistics
type APIUsageStats struct {
	TotalRequests      int64   `json:"totalRequests" db:"total_requests"`
	TotalTokens        int64   `json:"totalTokens" db:"total_tokens"`
	AvgResponseTime    float64 `json:"avgResponseTime" db:"avg_response_time"`
	SuccessRate        float64 `json:"successRate" db:"success_rate"`
	RequestsByEndpoint Data    `json:"requestsByEndpoint" db:"requests_by_endpoint"`
	RequestsByStatus   Data    `json:"requestsByStatus" db:"requests_by_status"`
}

// CreateAPIKeyParams parameters for creating an API key
type CreateAPIKeyParams struct {
	Name        string
	Value       string
	Type        string
	Claims      Data
	RateLimit   int
	AllowedIPs  []string
	Permissions []string
	ExpiresAt   *time.Time
	CreatedBy   *int
}

// CreateAPIKeyResult result from creating API key
type CreateAPIKeyResult struct {
	dal.DbResult
	KeyID *int `db:"key_id"`
}

// UpdateAPIKeyParams parameters for updating an API key
type UpdateAPIKeyParams struct {
	KeyID       int
	Name        *string
	RateLimit   *int
	AllowedIPs  *[]string
	Permissions *[]string
	IsActive    *bool
	ExpiresAt   *time.Time
}

// UpdateAPIKeyResult result from updating API key
type UpdateAPIKeyResult struct {
	dal.DbResult
}

// DeleteAPIKeyResult result from deleting API key
type DeleteAPIKeyResult struct {
	dal.DbResult
}

// UpdateAPIKeyLastUsedResult result from updating last used timestamp
type UpdateAPIKeyLastUsedResult struct {
	dal.DbResult
}

// TrackAPIUsageParams parameters for tracking API usage
type TrackAPIUsageParams struct {
	APIKeyID      int
	Endpoint      string
	Method        string
	StatusCode    int
	TokensUsed    int
	RequestTimeMs int
	IPAddress     *string
	UserAgent     *string
	ErrorMessage  *string
}

// TrackAPIUsageResult result from tracking API usage
type TrackAPIUsageResult struct {
	dal.DbResult
	UsageID *int `db:"usage_id"`
}

// APIKeyRepository defines database operations for API keys
type APIKeyRepository interface {
	// API key operations
	Create(ctx context.Context, params CreateAPIKeyParams) (*CreateAPIKeyResult, error)
	GetByValue(ctx context.Context, keyValue string) (*APIKey, error)
	GetByID(ctx context.Context, keyID int) (*APIKey, error)
	GetAll(ctx context.Context) ([]APIKey, error)
	Update(ctx context.Context, params UpdateAPIKeyParams) (*UpdateAPIKeyResult, error)
	UpdateLastUsed(ctx context.Context, keyID int) (*UpdateAPIKeyLastUsedResult, error)
	Delete(ctx context.Context, keyID int) (*DeleteAPIKeyResult, error)
}

// APIUsageRepository defines database operations for API usage tracking
type APIUsageRepository interface {
	// Usage tracking
	Track(ctx context.Context, params TrackAPIUsageParams) (*TrackAPIUsageResult, error)
	GetStats(ctx context.Context, keyID int, from, to *time.Time) (*APIUsageStats, error)
}

// APIKeyUseCase defines business logic for API key management
type APIKeyUseCase interface {
	// API key management
	CreateAPIKey(ctx context.Context, params CreateAPIKeyParams) Result[*APIKey]
	ValidateAPIKey(ctx context.Context, keyValue, ipAddress, endpoint string) Result[*APIKey]
	UpdateAPIKey(ctx context.Context, params UpdateAPIKeyParams) Result[Data]
	RevokeAPIKey(ctx context.Context, keyID int) Result[Data]
	ListAPIKeys(ctx context.Context) Result[[]APIKey]
	GetAPIKeyByID(ctx context.Context, keyID int) Result[*APIKey]
}
