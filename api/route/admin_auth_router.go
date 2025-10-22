package route

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"api-chatbot/api/request"
	d "api-chatbot/domain"
)

// Huma response types for admin authentication
type LoginResponse struct {
	Body d.Result[*d.TokenPairResponse]
}

type RefreshTokenResponse struct {
	Body d.Result[*d.TokenPairResponse]
}

type LogoutResponse struct {
	Body d.Result[d.Data]
}

type CreateAdminResponse struct {
	Body d.Result[*d.AdminUser]
}

func NewAdminAuthRouter(adminUseCase d.AdminUseCase, mux *http.ServeMux, humaAPI huma.API) {
	// Login endpoint
	huma.Register(humaAPI, huma.Operation{
		OperationID: "admin-login",
		Method:      "POST",
		Path:        "/admin/auth/login",
		Summary:     "Admin login",
		Description: "Authenticate admin user with username and password. Returns access token and refresh token.",
		Tags:        []string{"Admin Authentication"},
	}, func(ctx context.Context, input *struct {
		Body request.LoginRequest
	}) (*LoginResponse, error) {
		result := adminUseCase.Login(ctx, input.Body.Username, input.Body.Password, input.Body.DeviceAddress, input.Body.IdDevice)
		return &LoginResponse{Body: result}, nil
	})

	// Refresh token endpoint
	huma.Register(humaAPI, huma.Operation{
		OperationID: "admin-refresh-token",
		Method:      "POST",
		Path:        "/admin/auth/refresh",
		Summary:     "Refresh access token",
		Description: "Generate new access token and refresh token using valid refresh token. Implements token rotation for security.",
		Tags:        []string{"Admin Authentication"},
	}, func(ctx context.Context, input *struct {
		Body request.RefreshTokenRequest
	}) (*RefreshTokenResponse, error) {
		result := adminUseCase.RefreshToken(ctx, input.Body.RefreshToken, input.Body.DeviceAddress, input.Body.IdDevice)
		return &RefreshTokenResponse{Body: result}, nil
	})

	// Logout endpoint
	huma.Register(humaAPI, huma.Operation{
		OperationID: "admin-logout",
		Method:      "POST",
		Path:        "/admin/auth/logout",
		Summary:     "Admin logout",
		Description: "Revoke refresh token and log out admin user",
		Tags:        []string{"Admin Authentication"},
	}, func(ctx context.Context, input *struct {
		Body request.LogoutRequest
	}) (*LogoutResponse, error) {
		result := adminUseCase.Logout(ctx, input.Body.RefreshToken)
		return &LogoutResponse{Body: result}, nil
	})

	// Create admin user endpoint (protected - should require admin JWT later)
	huma.Register(humaAPI, huma.Operation{
		OperationID: "create-admin-user",
		Method:      "POST",
		Path:        "/admin/users/create",
		Summary:     "Create admin user",
		Description: "Create a new admin user with specified role and permissions. Requires super admin privileges.",
		Tags:        []string{"Admin Management"},
	}, func(ctx context.Context, input *struct {
		Body request.CreateAdminRequest
	}) (*CreateAdminResponse, error) {
		params := d.CreateAdminUserParams{
			Username:    input.Body.Username,
			Email:       input.Body.Email,
			Name:        input.Body.Name,
			Role:        input.Body.Role,
			Permissions: input.Body.Permissions,
			Claims:      input.Body.Claims,
		}
		result := adminUseCase.CreateAdmin(ctx, params, input.Body.Password)
		return &CreateAdminResponse{Body: result}, nil
	})
}
