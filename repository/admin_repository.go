package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
	"api-chatbot/internal/logger"
)

type adminRepository struct {
	dal *dal.DAL
}

const (
	spCreateAdminUser = "sp_create_admin_user"
)

// NewAdminRepository creates a new admin repository
func NewAdminRepository(d *dal.DAL) d.AdminRepository {
	return &adminRepository{
		dal: d,
	}
}

func (r *adminRepository) CreateAdminUser(ctx context.Context, params d.CreateAdminUserParams) (*d.CreateAdminUserResult, error) {
	permissionsJSON, err := json.Marshal(params.Permissions)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal permissions: %w", err)
	}

	claimsJSON, err := json.Marshal(params.Claims)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal claims: %w", err)
	}

	result, err := dal.ExecProc[d.CreateAdminUserResult](
		r.dal,
		ctx,
		spCreateAdminUser,
		params.Username,
		params.Email,
		params.PasswordHash,
		params.Name,
		params.Role,
		permissionsJSON,
		claimsJSON,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create admin user via %s: %w", spCreateAdminUser, err)
	}

	return result, nil
}

func (r *adminRepository) GetAdminByUsername(ctx context.Context, username string) (*d.AdminUser, error) {
	rows, err := dal.QueryRows[d.AdminUser](r.dal, ctx, "fn_get_admin_by_username", username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get admin by username: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("Admin user not found: %s", username)
	}

	return &rows[0], nil
}

func (r *adminRepository) GetAdminByID(ctx context.Context, id int) (*d.AdminUser, error) {
	rows, err := dal.QueryRows[d.AdminUser](r.dal, ctx, "fn_get_admin_by_id", id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get admin by ID: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("Admin user not found with ID: %d", id)
	}

	return &rows[0], nil
}

func (r *adminRepository) UpdateAdminLogin(ctx context.Context, adminID int, ipAddress string, resetFailedAttempts bool) (*d.UpdateLoginResult, error) {
	result, err := dal.ExecProc[d.UpdateLoginResult](
		r.dal,
		ctx,
		"sp_update_admin_login",
		adminID,
		ipAddress,
		resetFailedAttempts,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to update admin login: %w", err)
	}

	return result, nil
}

func (r *adminRepository) IncrementFailedAttempts(ctx context.Context, username string) (*d.IncrementFailedAttemptsResult, error) {
	result, err := dal.ExecProc[d.IncrementFailedAttemptsResult](
		r.dal,
		ctx,
		"sp_increment_failed_attempts",
		username,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to increment failed attempts: %w", err)
	}

	return result, nil
}

func (r *adminRepository) StoreRefreshToken(ctx context.Context, params d.StoreRefreshTokenParams) (*d.StoreRefreshTokenResult, error) {
	result, err := dal.ExecProc[d.StoreRefreshTokenResult](
		r.dal,
		ctx,
		"sp_store_refresh_token",
		params.AdminID,
		params.Token,
		params.TokenFamily,
		params.UserAgent,
		params.IPAddress,
		params.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to store refresh token: %w", err)
	}

	return result, nil
}

func (r *adminRepository) GetRefreshToken(ctx context.Context, token string) (*d.RefreshTokenWithAdmin, error) {
	rows, err := dal.QueryRows[d.RefreshTokenWithAdmin](r.dal, ctx, "fn_get_refresh_token", token)
	if err != nil {
		return nil, fmt.Errorf("Failed to get refresh token: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("Refresh token not found")
	}

	return &rows[0], nil
}

func (r *adminRepository) RevokeRefreshToken(ctx context.Context, token string, reason string) (*d.RevokeTokenResult, error) {
	result, err := dal.ExecProc[d.RevokeTokenResult](
		r.dal,
		ctx,
		"sp_revoke_refresh_token",
		token,
		reason,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to revoke refresh token: %w", err)
	}

	return result, nil
}

func (r *adminRepository) RevokeTokenFamily(ctx context.Context, tokenFamily string, reason string) (*d.RevokeTokenFamilyResult, error) {
	result, err := dal.ExecProc[d.RevokeTokenFamilyResult](
		r.dal,
		ctx,
		"sp_revoke_token_family",
		tokenFamily,
		reason,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to revoke token family: %w", err)
	}

	return result, nil
}

func (r *adminRepository) CleanupExpiredTokens(ctx context.Context) (*d.CleanupResult, error) {

	result, err := dal.ExecProc[d.CleanupResult](r.dal, ctx, "sp_cleanup_expired_tokens")

	if err != nil {
		return nil, fmt.Errorf("Failed to cleanup expired tokens: %w", err)
	}

	return result, nil
}

func (r *adminRepository) LogAuthEvent(ctx context.Context, userID *int, username, action, status string, ipAddress, userAgent *string, details d.Data) (*d.LogAuthEventResult, error) {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal auth event details: %w", err)
	}

	result, err := dal.ExecProc[d.LogAuthEventResult](
		r.dal,
		ctx,
		"sp_log_auth_event",
		userID,
		username,
		action,
		status,
		ipAddress,
		userAgent,
		detailsJSON,
	)

	if err != nil {
		logger.LogWarn(ctx, "Failed to log auth event", "username", username, "action", action, "error", err)
		// Don't fail the main operation if logging fails - return empty result with no error
		return nil, nil
	}

	return result, nil
}
