-- =====================================================
-- Migration 000036: Update sp_create_pending_registration
-- Add registration_step parameter to the stored procedure
-- =====================================================

-- Drop existing procedure first (to change signature)
DROP PROCEDURE IF EXISTS sp_create_pending_registration;

CREATE OR REPLACE PROCEDURE sp_create_pending_registration(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT o_pending_id INT,
    IN p_identity_number VARCHAR(20),
    IN p_whatsapp VARCHAR(50),
    IN p_name VARCHAR(100),
    IN p_email VARCHAR(100),
    IN p_phone VARCHAR(20) DEFAULT NULL,
    IN p_role VARCHAR(50) DEFAULT 'ROLE_STUDENT',
    IN p_user_type VARCHAR(20) DEFAULT 'institute',
    IN p_details JSONB DEFAULT '{}'::JSONB,
    IN p_registration_step VARCHAR(50) DEFAULT 'STEP_REQUEST_CEDULA',
    IN p_otp_code VARCHAR(6) DEFAULT NULL,
    IN p_otp_expires_at TIMESTAMP DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_existing_id INT;
BEGIN
    success := TRUE;
    code := 'OK';
    o_pending_id := NULL;

    -- Check if pending registration already exists
    SELECT pnd_id INTO v_existing_id
    FROM cht_pending_registrations
    WHERE pnd_whatsapp = p_whatsapp
      AND pnd_identity_number = p_identity_number
      AND pnd_verified = FALSE;

    IF v_existing_id IS NOT NULL THEN
        -- Update existing pending registration with new OTP
        UPDATE cht_pending_registrations
        SET pnd_name = COALESCE(p_name, pnd_name),
            pnd_email = COALESCE(p_email, pnd_email),
            pnd_phone = COALESCE(p_phone, pnd_phone),
            pnd_role = COALESCE(p_role, pnd_role),
            pnd_user_type = p_user_type,
            pnd_details = COALESCE(p_details, pnd_details),
            pnd_registration_step = p_registration_step,
            pnd_otp_code = p_otp_code,
            pnd_otp_generated_at = CURRENT_TIMESTAMP,
            pnd_otp_expires_at = p_otp_expires_at,
            pnd_otp_attempts = 0,
            pnd_updated_at = CURRENT_TIMESTAMP
        WHERE pnd_id = v_existing_id;

        o_pending_id := v_existing_id;
    ELSE
        -- Create new pending registration
        INSERT INTO cht_pending_registrations (
            pnd_identity_number,
            pnd_whatsapp,
            pnd_name,
            pnd_email,
            pnd_phone,
            pnd_role,
            pnd_user_type,
            pnd_details,
            pnd_registration_step,
            pnd_otp_code,
            pnd_otp_generated_at,
            pnd_otp_expires_at,
            pnd_otp_attempts,
            pnd_verified,
            pnd_created_at,
            pnd_updated_at
        )
        VALUES (
            p_identity_number,
            p_whatsapp,
            p_name,
            p_email,
            p_phone,
            p_role,
            p_user_type,
            p_details,
            p_registration_step,
            p_otp_code,
            CURRENT_TIMESTAMP,
            p_otp_expires_at,
            0,
            FALSE,
            CURRENT_TIMESTAMP,
            CURRENT_TIMESTAMP
        )
        RETURNING pnd_id INTO o_pending_id;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_CREATE_PENDING_REGISTRATION';
        o_pending_id := NULL;
        RAISE NOTICE 'Error creating pending registration: %', SQLERRM;
END;
$$;

COMMENT ON PROCEDURE sp_create_pending_registration IS 'Creates or updates pending registration with registration step tracking';
