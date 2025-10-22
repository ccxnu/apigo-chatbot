package jwttoken

import (
	d "api-chatbot/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims represents JWT claims with custom fields
type CustomClaims struct {
	UserID      int                    `json:"userId"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Name        string                 `json:"name"`
	Role        string                 `json:"role"`
	Permissions []string               `json:"permissions,omitempty"`
	Claims      d.Data `json:"claims,omitempty"` // Custom extensible claims
	TokenType   string                 `json:"tokenType"`        // "access" or "refresh"
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents minimal claims for refresh tokens
type RefreshTokenClaims struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	TokenFamily string `json:"tokenFamily"`
	TokenType   string `json:"tokenType"` // Always "refresh"
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens together
type TokenPair struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    string    `json:"tokenType"` // "Bearer"
	ExpiresIn    int       `json:"expiresIn"` // Seconds until access token expires
	ExpiresAt    time.Time `json:"expiresAt"`
}

// TokenMetadata contains information about token generation
type TokenMetadata struct {
	UserID      int
	Username    string
	Email       string
	Name        string
	Role        string
	Permissions []string
	Claims      d.Data
	TokenFamily string // For refresh token rotation
	IPAddress   string
	UserAgent   string
}
