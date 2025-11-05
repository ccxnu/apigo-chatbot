-- =====================================================
-- Rollback Migration 000037
-- =====================================================

-- Restore original function without registration_step
CREATE OR REPLACE FUNCTION fn_get_pending_registration_by_whatsapp(
    i_whatsapp VARCHAR(50)
)
RETURNS TABLE (
    pending_id INT,
    identity_number VARCHAR(20),
    whatsapp VARCHAR(50),
    name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(50),
    user_type VARCHAR(20),
    details JSONB,
    otp_expires_at TIMESTAMP,
    otp_attempts INT,
    verified BOOLEAN,
    created_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        pnd_id,
        pnd_identity_number,
        pnd_whatsapp,
        pnd_name,
        pnd_email,
        pnd_phone,
        pnd_role,
        pnd_user_type,
        pnd_details,
        pnd_otp_expires_at,
        pnd_otp_attempts,
        pnd_verified,
        pnd_created_at
    FROM cht_pending_registrations
    WHERE pnd_whatsapp = i_whatsapp
      AND pnd_verified = FALSE
    ORDER BY pnd_created_at DESC
    LIMIT 1;
END;
$$;

-- Remove STEP_INIT parameter
DELETE FROM cht_parameters WHERE prm_code = 'REG_STEP_INIT';
