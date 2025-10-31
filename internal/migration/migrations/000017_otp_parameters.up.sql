-- =====================================================
-- OTP Registration Parameters
-- =====================================================

DO $$
BEGIN
    -- Update EMAIL_CONFIG with Tikee URL and sender email
    INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
    VALUES (
        'EMAIL_CONFIGURATION',
        'EMAIL_CONFIG',
        '{"senderEmail": "automatizaciones@tikee.tech", "tikeeURL": "http://20.84.48.225:5056/api/emails/enviarDirecto"}'::jsonb,
        'Email sender and Tikee API configuration for OTP emails'
    )
    ON CONFLICT (prm_code)
    DO UPDATE SET
        prm_data = '{"senderEmail": "automatizaciones@tikee.tech", "tikeeURL": "http://20.84.48.225:5056/api/emails/enviarDirecto"}'::jsonb,
        prm_description = 'Email sender and Tikee API configuration for OTP emails';

    -- Add OTP_EXPIRATION_MINUTES parameter
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'OTP_EXPIRATION_MINUTES') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES (
            'OTP_CONFIGURATION',
            'OTP_EXPIRATION_MINUTES',
            '{"minutes": 10}'::jsonb,
            'OTP code expiration time in minutes'
        );
    END IF;

    -- Add OTP-related error codes
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_INVALID_OTP') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_INVALID_OTP', '{"message": "❌ Código incorrecto. Por favor verifica e intenta nuevamente."}'::jsonb, 'Invalid OTP code');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_OTP_EXPIRED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_OTP_EXPIRED', '{"message": "⏰ El código ha expirado. Escribe reenviar para generar un nuevo código."}'::jsonb, 'OTP code has expired');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_MAX_ATTEMPTS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_MAX_ATTEMPTS', '{"message": "🚫 Has excedido el número máximo de intentos. Escribe reenviar para generar un nuevo código."}'::jsonb, 'Maximum OTP verification attempts exceeded');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REG') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_NO_PENDING_REG', '{"message": "❌ No tienes un registro pendiente. Por favor envía tu cédula para iniciar el registro."}'::jsonb, 'No pending registration found');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REGISTRATION') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_NO_PENDING_REGISTRATION', '{"message": "❌ No tienes un registro pendiente. Por favor envía tu cédula para iniciar el registro."}'::jsonb, 'No pending registration found (alias)');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_IDENTITY_ALREADY_REGISTERED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_IDENTITY_ALREADY_REGISTERED', '{"message": "❌ Esta cédula ya está registrada con otro número de WhatsApp."}'::jsonb, 'Identity number already registered with different WhatsApp');
    END IF;

END $$;

-- Comments
COMMENT ON TABLE cht_parameters IS 'System parameters including OTP configuration and error codes for registration';
