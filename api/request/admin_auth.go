package request

import "api-chatbot/domain"

// LoginRequest represents admin login credentials
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" doc:"Admin username"`
	Password string `json:"password" validate:"required,min=8" doc:"Admin password"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" doc:"Valid refresh token"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required" doc:"Refresh token to revoke"`
}

// CreateAdminRequest represents admin user creation request
type CreateAdminRequest struct {
	domain.Base
	Username    string              `json:"username" validate:"required,min=3,max=50" doc:"Admin username"`
	Email       string              `json:"email" validate:"required,email,max=100" doc:"Admin email address"`
	Password    string              `json:"password" validate:"required,min=8" doc:"Admin password"`
	Name        string              `json:"name" validate:"required,min=3,max=100" doc:"Full name"`
	Role        string              `json:"role" validate:"required,min=2,max=50" doc:"Admin role (e.g., super_admin, admin, moderator)"`
	Permissions []string            `json:"permissions,omitempty" doc:"Array of permission strings"`
	Claims      map[string]any      `json:"claims,omitempty" doc:"Custom claims as key-value pairs"`
}
