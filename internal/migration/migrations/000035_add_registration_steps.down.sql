-- =====================================================
-- Rollback Migration 000035: Remove Registration Steps
-- =====================================================

-- Drop stored procedure
DROP PROCEDURE IF EXISTS sp_update_registration_step;

-- Remove registration step column
ALTER TABLE public.cht_pending_registrations
DROP COLUMN IF EXISTS pnd_registration_step;

-- Drop index
DROP INDEX IF EXISTS idx_pending_reg_step;

-- Remove registration step parameters
DELETE FROM cht_parameters WHERE prm_code IN (
    'REG_STEP_REQUEST_CEDULA',
    'REG_STEP_VALIDATE_ACADEMICOK',
    'REG_STEP_SELECT_USER_TYPE',
    'REG_STEP_REQUEST_EMAIL_NAME',
    'REG_STEP_SEND_OTP',
    'REG_STEP_VERIFY_OTP',
    'REG_STEP_CREATE_USER',
    'REG_STEP_COMPLETED',
    'REG_MSG_WAITING_CEDULA',
    'REG_MSG_WAITING_USER_TYPE',
    'REG_MSG_WAITING_EMAIL_NAME',
    'REG_MSG_WAITING_OTP'
);
