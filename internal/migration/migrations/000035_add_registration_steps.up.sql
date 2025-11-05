-- =====================================================
-- Migration 000035: Add Registration Step Tracking
-- =====================================================
-- This migration adds:
-- 1. Registration step parameters to track the registration flow
-- 2. A column to cht_pending_registrations to track current step

-- =====================================================
-- Add registration step column to pending registrations
-- =====================================================
ALTER TABLE public.cht_pending_registrations
ADD COLUMN IF NOT EXISTS pnd_registration_step VARCHAR(50) DEFAULT 'STEP_REQUEST_CEDULA';

-- Add index for step queries
CREATE INDEX IF NOT EXISTS idx_pending_reg_step ON cht_pending_registrations(pnd_registration_step);

-- =====================================================
-- Registration Step Parameters
-- =====================================================
do $$
begin
    -- =====================================================
    -- Registration Steps
    -- =====================================================

    -- Step 1: Request Cedula
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_REQUEST_CEDULA') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_REQUEST_CEDULA',
            '{"step": 1, "code": "STEP_REQUEST_CEDULA", "description": "Solicitar número de cédula al usuario"}'::jsonb,
            'Initial step - request identity number from user'
        );
    end if;

    -- Step 2: Validate with AcademicOK
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_VALIDATE_ACADEMICOK') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_VALIDATE_ACADEMICOK',
            '{"step": 2, "code": "STEP_VALIDATE_ACADEMICOK", "description": "Validar cédula con sistema AcademicOK"}'::jsonb,
            'Validate identity number against AcademicOK API'
        );
    end if;

    -- Step 3: User Type Selection (if API fails)
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_SELECT_USER_TYPE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_SELECT_USER_TYPE',
            '{"step": 3, "code": "STEP_SELECT_USER_TYPE", "description": "Usuario selecciona su tipo (estudiante, docente, externo)"}'::jsonb,
            'User selects their type when API validation fails'
        );
    end if;

    -- Step 4: Request Email and Name
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_REQUEST_EMAIL_NAME') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_REQUEST_EMAIL_NAME',
            '{"step": 4, "code": "STEP_REQUEST_EMAIL_NAME", "description": "Solicitar nombre completo y correo electrónico"}'::jsonb,
            'Request full name and email from external users'
        );
    end if;

    -- Step 5: Send OTP
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_SEND_OTP') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_SEND_OTP',
            '{"step": 5, "code": "STEP_SEND_OTP", "description": "Enviar código OTP al correo del usuario"}'::jsonb,
            'Send OTP verification code to user email'
        );
    end if;

    -- Step 6: Verify OTP
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_VERIFY_OTP') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_VERIFY_OTP',
            '{"step": 6, "code": "STEP_VERIFY_OTP", "description": "Verificar código OTP ingresado por el usuario"}'::jsonb,
            'Verify OTP code entered by user'
        );
    end if;

    -- Step 7: Create User Account
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_CREATE_USER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_CREATE_USER',
            '{"step": 7, "code": "STEP_CREATE_USER", "description": "Crear cuenta de usuario en el sistema"}'::jsonb,
            'Create user account after successful OTP verification'
        );
    end if;

    -- Step 8: Registration Complete
    if not exists (select 1 from cht_parameters where prm_code = 'REG_STEP_COMPLETED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_STEPS',
            'REG_STEP_COMPLETED',
            '{"step": 8, "code": "STEP_COMPLETED", "description": "Registro completado exitosamente"}'::jsonb,
            'Registration process completed successfully'
        );
    end if;

    -- =====================================================
    -- Registration Step Messages
    -- =====================================================

    -- Message when waiting for cedula
    if not exists (select 1 from cht_parameters where prm_code = 'REG_MSG_WAITING_CEDULA') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_MESSAGES',
            'REG_MSG_WAITING_CEDULA',
            '{"message": "⏳ Estoy esperando tu número de cédula para continuar con el registro.\n\nPor favor envía tu cédula de 10 dígitos.\n\nEjemplo: 1234567890"}'::jsonb,
            'Reminder message when user sends other text instead of cedula'
        );
    end if;

    -- Message when waiting for user type selection
    if not exists (select 1 from cht_parameters where prm_code = 'REG_MSG_WAITING_USER_TYPE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_MESSAGES',
            'REG_MSG_WAITING_USER_TYPE',
            '{"message": "⏳ Por favor selecciona tu tipo de usuario enviando el número correspondiente:\n\n*1* - 🎓 Estudiante\n*2* - 👨‍🏫 Docente\n*3* - 👤 Usuario externo"}'::jsonb,
            'Reminder when user does not select a valid type'
        );
    end if;

    -- Message when waiting for email and name
    if not exists (select 1 from cht_parameters where prm_code = 'REG_MSG_WAITING_EMAIL_NAME') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_MESSAGES',
            'REG_MSG_WAITING_EMAIL_NAME',
            '{"message": "⏳ Necesito tu nombre completo y correo electrónico.\n\nFormato: *Nombre Completo / correo@email.com*\n\nEjemplo:\nJuan Pérez / juan.perez@gmail.com"}'::jsonb,
            'Reminder when user does not provide valid email and name'
        );
    end if;

    -- Message when waiting for OTP
    if not exists (select 1 from cht_parameters where prm_code = 'REG_MSG_WAITING_OTP') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'REGISTRATION_MESSAGES',
            'REG_MSG_WAITING_OTP',
            '{"message": "⏳ Estoy esperando el código de verificación que te envié por correo.\n\nPor favor ingresa el código de 6 dígitos.\n\nSi no lo recibiste, escribe ''reenviar''."}'::jsonb,
            'Reminder when user sends text instead of OTP code'
        );
    end if;

end $$;

-- =====================================================
-- Stored Procedure: sp_update_registration_step
-- Updates the registration step for a pending registration
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_update_registration_step(
    IN p_whatsapp VARCHAR(50),
    IN p_registration_step VARCHAR(50),
    OUT o_success BOOLEAN,
    OUT o_code VARCHAR(50),
    OUT o_message TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    -- Update registration step
    UPDATE cht_pending_registrations
    SET pnd_registration_step = p_registration_step,
        pnd_updated_at = CURRENT_TIMESTAMP
    WHERE pnd_whatsapp = p_whatsapp;

    -- Check if row was updated
    IF NOT FOUND THEN
        o_success := FALSE;
        o_code := 'ERR_NO_PENDING_REGISTRATION';
        o_message := 'No pending registration found for this WhatsApp number';
        RETURN;
    END IF;

    o_success := TRUE;
    o_code := 'OK';
    o_message := 'Registration step updated successfully';
END;
$$;

-- =====================================================
-- Comment
-- =====================================================
COMMENT ON COLUMN cht_pending_registrations.pnd_registration_step IS 'Current step in the registration flow (STEP_REQUEST_CEDULA, STEP_VALIDATE_ACADEMICOK, etc.)';
COMMENT ON PROCEDURE sp_update_registration_step IS 'Updates the registration step for a pending registration by WhatsApp number';
