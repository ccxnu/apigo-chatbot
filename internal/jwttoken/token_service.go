package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenService handles JWT token creation and validation
type TokenService struct {
	accessSecret       string
	refreshSecret      string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(accessSecret, refreshSecret string, accessExpiryHours, refreshExpiryHours int) *TokenService {
	return &TokenService{
		accessSecret:       accessSecret,
		refreshSecret:      refreshSecret,
		accessTokenExpiry:  time.Duration(accessExpiryHours) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshExpiryHours) * time.Hour,
	}
}

// CreateTokenPair creates both access and refresh tokens
func (ts *TokenService) CreateTokenPair(metadata TokenMetadata) (*TokenPair, error) {
	// Create access token
	accessToken, expiresAt, err := ts.CreateAccessToken(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := ts.CreateRefreshToken(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(ts.accessTokenExpiry.Seconds()),
		ExpiresAt:    expiresAt,
	}, nil
}

// CreateAccessToken creates a new access token with custom claims
func (ts *TokenService) CreateAccessToken(metadata TokenMetadata) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ts.accessTokenExpiry)

	claims := CustomClaims{
		UserID:      metadata.UserID,
		Username:    metadata.Username,
		Email:       metadata.Email,
		Name:        metadata.Name,
		Role:        metadata.Role,
		Permissions: metadata.Permissions,
		Claims:      metadata.Claims,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "chatbot-admin",
			Subject:   metadata.Username,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.accessSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// CreateRefreshToken creates a new refresh token
func (ts *TokenService) CreateRefreshToken(metadata TokenMetadata) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ts.refreshTokenExpiry)

	// Generate token family if not provided
	tokenFamily := metadata.TokenFamily
	if tokenFamily == "" {
		tokenFamily = uuid.New().String()
	}

	claims := RefreshTokenClaims{
		UserID:      metadata.UserID,
		Username:    metadata.Username,
		TokenFamily: tokenFamily,
		TokenType:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "chatbot-admin",
			Subject:   metadata.Username,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.refreshSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates and parses an access token
func (ts *TokenService) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ts.accessSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify token type
	if claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid token type: expected access, got %s", claims.TokenType)
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (ts *TokenService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ts.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verify token type
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type: expected refresh, got %s", claims.TokenType)
	}

	return claims, nil
}

// ExtractClaimsWithoutValidation extracts claims without validating signature (use cautiously)
func (ts *TokenService) ExtractClaimsWithoutValidation(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	return claims, nil
}

// IsTokenExpired checks if a token is expired without full validation
func (ts *TokenService) IsTokenExpired(tokenString string) (bool, error) {
	claims, err := ts.ExtractClaimsWithoutValidation(tokenString)
	if err != nil {
		return false, err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid expiration claim")
	}

	return time.Now().Unix() > int64(exp), nil
}

// GetTokenExpiry returns the expiry duration for access tokens
func (ts *TokenService) GetTokenExpiry() time.Duration {
	return ts.accessTokenExpiry
}

// GetRefreshTokenExpiry returns the expiry duration for refresh tokens
func (ts *TokenService) GetRefreshTokenExpiry() time.Duration {
	return ts.refreshTokenExpiry
}
