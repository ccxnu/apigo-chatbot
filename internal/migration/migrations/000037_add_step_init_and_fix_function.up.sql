-- =====================================================
-- Migration 000037: Add STEP_INIT and fix fn_get_pending_registration_by_whatsapp
-- =====================================================

-- =====================================================
-- Add STEP_INIT parameter (Step 0 - Registration initiated)
-- =====================================================
do $$
begin
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_INIT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_INIT',
            '{"step": 0, "code": "STEP_INIT", "description": "Registro iniciado - esperando número de cédula"}'::jsonb,
            'Initial step - user started registration, waiting for cedula'
        );
    end if;
end $$;

-- =====================================================
-- Update fn_get_pending_registration_by_whatsapp to include registration_step
-- =====================================================
-- Drop existing function first (to change return type)
DROP FUNCTION IF EXISTS fn_get_pending_registration_by_whatsapp;

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
    registration_step VARCHAR(50),
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
        pnd_registration_step,
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

COMMENT ON FUNCTION fn_get_pending_registration_by_whatsapp IS 'Retrieves pending registration by WhatsApp number including registration step';
