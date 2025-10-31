-- =====================================================
-- Rollback OTP Registration Parameters
-- =====================================================

-- Remove OTP-related error codes
DELETE FROM cht_parameters WHERE prm_code = 'ERR_INVALID_OTP';
DELETE FROM cht_parameters WHERE prm_code = 'ERR_OTP_EXPIRED';
DELETE FROM cht_parameters WHERE prm_code = 'ERR_MAX_ATTEMPTS';
DELETE FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REG';
DELETE FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REGISTRATION';
DELETE FROM cht_parameters WHERE prm_code = 'ERR_IDENTITY_ALREADY_REGISTERED';

-- Remove OTP configuration
DELETE FROM cht_parameters WHERE prm_code = 'OTP_EXPIRATION_MINUTES';

-- Revert EMAIL_CONFIG to basic configuration
UPDATE cht_parameters
SET prm_data = '{"sender": "noreply@example.com"}'::jsonb,
    prm_description = 'Email sender configuration'
WHERE prm_code = 'EMAIL_CONFIG';
