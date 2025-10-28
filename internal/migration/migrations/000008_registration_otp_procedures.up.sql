-- =====================================================
-- Stored Procedures for OTP Registration Management
-- =====================================================

-- =====================================================
-- Create or Update Pending Registration
-- =====================================================
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
            pnd_otp_code,
            pnd_otp_generated_at,
            pnd_otp_expires_at,
            pnd_otp_attempts,
            pnd_verified
        ) VALUES (
            p_identity_number,
            p_whatsapp,
            p_name,
            p_email,
            p_phone,
            p_role,
            p_user_type,
            p_details,
            p_otp_code,
            CURRENT_TIMESTAMP,
            p_otp_expires_at,
            0,
            FALSE
        ) RETURNING pnd_id INTO o_pending_id;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_CREATE_PENDING_REG';
        o_pending_id := NULL;
        RAISE NOTICE 'Error creating pending registration: %', SQLERRM;
END;
$$;

-- =====================================================
-- Verify OTP Code
-- =====================================================
CREATE OR REPLACE FUNCTION fn_verify_otp_code(
    i_whatsapp VARCHAR(50),
    i_otp_code VARCHAR(6),
    i_ip_address INET DEFAULT NULL
)
RETURNS TABLE (
    success BOOLEAN,
    code VARCHAR(10),
    message TEXT,
    pending_id INT,
    identity_number VARCHAR(20),
    name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(50),
    user_type VARCHAR(20),
    details JSONB
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_pending_record RECORD;
    v_max_attempts INT := 5;
    v_verified BOOLEAN;
BEGIN
    -- Find pending registration by WhatsApp
    SELECT *
    INTO v_pending_record
    FROM cht_pending_registrations
    WHERE pnd_whatsapp = i_whatsapp
      AND pnd_verified = FALSE
    ORDER BY pnd_created_at DESC
    LIMIT 1;

    -- No pending registration found
    IF NOT FOUND THEN
        RETURN QUERY SELECT
            FALSE,
            'ERR_NO_PENDING_REG'::VARCHAR,
            'No pending registration found for this WhatsApp number'::TEXT,
            NULL::INT, NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR,
            NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR, NULL::JSONB;
        RETURN;
    END IF;

    -- Check if OTP has expired
    IF v_pending_record.pnd_otp_expires_at < CURRENT_TIMESTAMP THEN
        -- Log failed attempt
        INSERT INTO cht_otp_verification_log (ovl_fk_pending, ovl_code_attempted, ovl_success, ovl_ip_address)
        VALUES (v_pending_record.pnd_id, i_otp_code, FALSE, i_ip_address);

        RETURN QUERY SELECT
            FALSE,
            'ERR_OTP_EXPIRED'::VARCHAR,
            'OTP code has expired. Please request a new code'::TEXT,
            NULL::INT, NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR,
            NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR, NULL::JSONB;
        RETURN;
    END IF;

    -- Check if max attempts exceeded
    IF v_pending_record.pnd_otp_attempts >= v_max_attempts THEN
        -- Log failed attempt
        INSERT INTO cht_otp_verification_log (ovl_fk_pending, ovl_code_attempted, ovl_success, ovl_ip_address)
        VALUES (v_pending_record.pnd_id, i_otp_code, FALSE, i_ip_address);

        RETURN QUERY SELECT
            FALSE,
            'ERR_MAX_ATTEMPTS'::VARCHAR,
            'Maximum verification attempts exceeded. Please request a new code'::TEXT,
            NULL::INT, NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR,
            NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR, NULL::JSONB;
        RETURN;
    END IF;

    -- Verify OTP code
    v_verified := (v_pending_record.pnd_otp_code = i_otp_code);

    IF v_verified THEN
        -- Mark as verified
        UPDATE cht_pending_registrations
        SET pnd_verified = TRUE,
            pnd_updated_at = CURRENT_TIMESTAMP
        WHERE pnd_id = v_pending_record.pnd_id;

        -- Log successful attempt
        INSERT INTO cht_otp_verification_log (ovl_fk_pending, ovl_code_attempted, ovl_success, ovl_ip_address)
        VALUES (v_pending_record.pnd_id, i_otp_code, TRUE, i_ip_address);

        RETURN QUERY SELECT
            TRUE,
            'COD_OK'::VARCHAR,
            'OTP verified successfully'::TEXT,
            v_pending_record.pnd_id,
            v_pending_record.pnd_identity_number,
            v_pending_record.pnd_name,
            v_pending_record.pnd_email,
            v_pending_record.pnd_phone,
            v_pending_record.pnd_role,
            v_pending_record.pnd_user_type,
            v_pending_record.pnd_details;
    ELSE
        -- Increment failed attempts
        UPDATE cht_pending_registrations
        SET pnd_otp_attempts = pnd_otp_attempts + 1,
            pnd_updated_at = CURRENT_TIMESTAMP
        WHERE pnd_id = v_pending_record.pnd_id;

        -- Log failed attempt
        INSERT INTO cht_otp_verification_log (ovl_fk_pending, ovl_code_attempted, ovl_success, ovl_ip_address)
        VALUES (v_pending_record.pnd_id, i_otp_code, FALSE, i_ip_address);

        RETURN QUERY SELECT
            FALSE,
            'ERR_INVALID_OTP'::VARCHAR,
            FORMAT('Invalid OTP code. %s attempts remaining', v_max_attempts - v_pending_record.pnd_otp_attempts - 1)::TEXT,
            NULL::INT, NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR,
            NULL::VARCHAR, NULL::VARCHAR, NULL::VARCHAR, NULL::JSONB;
    END IF;

END;
$$;

-- =====================================================
-- Get Pending Registration by WhatsApp
-- =====================================================
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

-- =====================================================
-- Delete Pending Registration After Successful Registration
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_delete_pending_registration(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_pending_id INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := TRUE;
    code := 'OK';

    DELETE FROM cht_pending_registrations
    WHERE pnd_id = p_pending_id;

    IF NOT FOUND THEN
        success := FALSE;
        code := 'ERR_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_DELETE_PENDING';
        RAISE NOTICE 'Error deleting pending registration: %', SQLERRM;
END;
$$;

-- =====================================================
-- Clean up expired pending registrations (run periodically)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_cleanup_expired_pending_registrations()
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    v_deleted_count INT;
BEGIN
    -- Delete pending registrations older than 24 hours and not verified
    DELETE FROM cht_pending_registrations
    WHERE pnd_verified = FALSE
      AND pnd_created_at < CURRENT_TIMESTAMP - INTERVAL '24 hours';

    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;

    RETURN v_deleted_count;
END;
$$;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON PROCEDURE sp_create_pending_registration IS 'Create or update pending registration with OTP code';
COMMENT ON FUNCTION fn_verify_otp_code IS 'Verify OTP code and return user data if successful';
COMMENT ON FUNCTION fn_get_pending_registration_by_whatsapp IS 'Get pending registration by WhatsApp number';
COMMENT ON PROCEDURE sp_delete_pending_registration IS 'Delete pending registration after successful user creation';
COMMENT ON FUNCTION fn_cleanup_expired_pending_registrations IS 'Clean up expired pending registrations older than 24 hours';
