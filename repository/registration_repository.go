package repository

import (
	"context"
	"fmt"

	"api-chatbot/api/dal"
	d "api-chatbot/domain"
)

const (
	// Functions (Read-only)
	fnGetPendingRegistrationByWhatsApp = "fn_get_pending_registration_by_whatsapp"
	fnVerifyOTPCode                    = "fn_verify_otp_code"
	fnCleanupExpiredPendingRegs        = "fn_cleanup_expired_pending_registrations"
	// Stored Procedures (Writes)
	spCreatePendingRegistration = "sp_create_pending_registration"
	spDeletePendingRegistration = "sp_delete_pending_registration"
)

type registrationRepository struct {
	dal *dal.DAL
}

func NewRegistrationRepository(dal *dal.DAL) d.RegistrationRepository {
	return &registrationRepository{
		dal: dal,
	}
}

// CreatePendingRegistration creates or updates a pending registration with OTP
func (r *registrationRepository) CreatePendingRegistration(
	ctx context.Context,
	params d.CreatePendingRegistrationParams,
) (*d.CreatePendingRegistrationResult, error) {
	result, err := dal.ExecProc[d.CreatePendingRegistrationResult](
		r.dal,
		ctx,
		spCreatePendingRegistration,
		params.IdentityNumber,
		params.WhatsApp,
		params.Name,
		params.Email,
		params.Phone,
		params.Role,
		params.UserType,
		params.Details,
		params.OTPCode,
		params.OTPExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", spCreatePendingRegistration, err)
	}

	return result, nil
}

// VerifyOTP verifies the OTP code for a WhatsApp number
func (r *registrationRepository) VerifyOTP(
	ctx context.Context,
	params d.VerifyOTPParams,
) (*d.OTPVerificationResult, error) {
	results, err := dal.QueryRows[d.OTPVerificationResult](
		r.dal,
		ctx,
		fnVerifyOTPCode,
		params.WhatsApp,
		params.OTPCode,
		params.IPAddress,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute %s: %w", fnVerifyOTPCode, err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no result returned from %s", fnVerifyOTPCode)
	}

	return &results[0], nil
}

// GetPendingByWhatsApp retrieves pending registration by WhatsApp number
func (r *registrationRepository) GetPendingByWhatsApp(
	ctx context.Context,
	whatsapp string,
) (*d.PendingRegistration, error) {
	results, err := dal.QueryRows[d.PendingRegistration](
		r.dal,
		ctx,
		fnGetPendingRegistrationByWhatsApp,
		whatsapp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get pending registration via %s: %w", fnGetPendingRegistrationByWhatsApp, err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

// DeletePendingRegistration deletes a pending registration after successful user creation
func (r *registrationRepository) DeletePendingRegistration(
	ctx context.Context,
	pendingID int,
) error {
	result, err := dal.ExecProc[d.DeletePendingRegistrationResult](
		r.dal,
		ctx,
		spDeletePendingRegistration,
		pendingID,
	)

	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", spDeletePendingRegistration, err)
	}

	if !result.Success {
		return fmt.Errorf("failed to delete pending registration: code=%s", result.Code)
	}

	return nil
}

// CleanupExpiredRegistrations removes expired pending registrations
func (r *registrationRepository) CleanupExpiredRegistrations(ctx context.Context) (int, error) {
	type CleanupResult struct {
		Count int `db:"fn_cleanup_expired_pending_registrations"`
	}

	results, err := dal.QueryRows[CleanupResult](
		r.dal,
		ctx,
		fnCleanupExpiredPendingRegs,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to execute %s: %w", fnCleanupExpiredPendingRegs, err)
	}

	if len(results) == 0 {
		return 0, nil
	}

	return results[0].Count, nil
}
