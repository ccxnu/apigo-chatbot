-- =====================================================
-- OTP Registration System Parameters
-- =====================================================

DO $$
BEGIN
    -- =====================================================
    -- ERROR CODES for OTP Registration
    -- =====================================================

    -- ERR_USER_ALREADY_EXISTS
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_USER_ALREADY_EXISTS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_USER_ALREADY_EXISTS',
                '{"message": "El usuario ya está registrado"}'::jsonb,
                'User already exists');
    END IF;

    -- ERR_IDENTITY_ALREADY_REGISTERED
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_IDENTITY_ALREADY_REGISTERED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_IDENTITY_ALREADY_REGISTERED',
                '{"message": "Esta cédula ya está registrada con otro número de WhatsApp"}'::jsonb,
                'Identity number already registered with different WhatsApp');
    END IF;

    -- ERR_EXTERNAL_USER_INFO_REQUIRED
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_EXTERNAL_USER_INFO_REQUIRED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_EXTERNAL_USER_INFO_REQUIRED',
                '{"message": "Usuario externo - se requiere información adicional"}'::jsonb,
                'External user requires additional information');
    END IF;

    -- ERR_IDENTITY_NOT_FOUND
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_IDENTITY_NOT_FOUND') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_IDENTITY_NOT_FOUND',
                '{"message": "Número de identificación no encontrado en la base de datos institucional"}'::jsonb,
                'Identity number not found in institute database');
    END IF;

    -- ERR_INVALID_IDENTITY
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_INVALID_IDENTITY') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_INVALID_IDENTITY',
                '{"message": "Número de identificación inválido"}'::jsonb,
                'Invalid identity number');
    END IF;

    -- ERR_NO_PENDING_REGISTRATION
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REGISTRATION') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_NO_PENDING_REGISTRATION',
                '{"message": "No hay un registro pendiente para este usuario"}'::jsonb,
                'No pending registration found');
    END IF;

    -- ERR_INVALID_OTP
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_INVALID_OTP') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_INVALID_OTP',
                '{"message": "Código de verificación incorrecto"}'::jsonb,
                'Invalid OTP code');
    END IF;

    -- ERR_OTP_EXPIRED
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_OTP_EXPIRED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_OTP_EXPIRED',
                '{"message": "El código de verificación ha expirado"}'::jsonb,
                'OTP code has expired');
    END IF;

    -- ERR_MAX_ATTEMPTS
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_MAX_ATTEMPTS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_MAX_ATTEMPTS',
                '{"message": "Número máximo de intentos excedido"}'::jsonb,
                'Maximum verification attempts exceeded');
    END IF;

    -- ERR_CREATE_PENDING_REG
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_CREATE_PENDING_REG') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_CREATE_PENDING_REG',
                '{"message": "Error al crear el registro pendiente"}'::jsonb,
                'Error creating pending registration');
    END IF;

    -- ERR_NO_PENDING_REG
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_NO_PENDING_REG') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_NO_PENDING_REG',
                '{"message": "No se encontró un registro pendiente"}'::jsonb,
                'No pending registration found');
    END IF;

    -- ERR_DELETE_PENDING
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_DELETE_PENDING') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_DELETE_PENDING',
                '{"message": "Error al eliminar el registro pendiente"}'::jsonb,
                'Error deleting pending registration');
    END IF;

    -- ERR_INTERNAL
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_INTERNAL') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_INTERNAL',
                '{"message": "Error interno del sistema"}'::jsonb,
                'Internal system error');
    END IF;

    -- =====================================================
    -- CONFIGURATION PARAMETERS
    -- =====================================================

    -- OTP Expiration Time
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'OTP_EXPIRATION_MINUTES') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'OTP_EXPIRATION_MINUTES',
                '{"minutes": 10}'::jsonb,
                'OTP code expiration time in minutes');
    END IF;

    -- Email OTP Template
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EMAIL_OTP_TEMPLATE') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EMAIL_OTP_TEMPLATE',
                '{"subject": "Código de verificación - Chatbot ISTS", "html": ""}'::jsonb,
                'HTML template for OTP verification emails (empty uses default template)');
    END IF;

    -- Tikee Email Service Configuration
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'TIKEE_EMAIL_SERVICE') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'TIKEE_EMAIL_SERVICE',
                '{"url": "http://20.84.48.225:5056/api/emails/enviarDirecto", "sender": "automatizaciones@tikee.tech"}'::jsonb,
                'Tikee email service configuration');
    END IF;

    -- WhatsApp Registration Messages
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'MESSAGE_REQUEST_CEDULA') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('MESSAGES', 'MESSAGE_REQUEST_CEDULA',
                '{"message": "👋 ¡Hola! Bienvenido al asistente virtual del Instituto.\n\nPara poder ayudarte, necesito que te registres primero.\n\nPor favor, envíame tu número de cédula (10 dígitos).\n\nEjemplo: 1234567890"}'::jsonb,
                'Message requesting cedula from user');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'MESSAGE_HELP') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('MESSAGES', 'MESSAGE_HELP',
                '{"message": "👋 *Bienvenido al Asistente del Instituto*\n\n Puedes hacer preguntas sobre:\n• Carreras y programas académicos\n• Requisitos de admisión\n• Horarios y calendario académico\n• Servicios estudiantiles\n• Y mucho más...\n\nEscribe tu pregunta y te ayudaré con gusto."}'::jsonb,
                'Help message shown after registration');
    END IF;

    -- =====================================================
    -- ROLES
    -- =====================================================

    -- ROLE_STUDENT
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ROLE_STUDENT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ROLES', 'ROLE_STUDENT',
                '{"name": "Estudiante", "permissions": []}'::jsonb,
                'Student role');
    END IF;

    -- ROLE_PROFESSOR
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ROLE_PROFESSOR') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ROLES', 'ROLE_PROFESSOR',
                '{"name": "Docente", "permissions": []}'::jsonb,
                'Professor role');
    END IF;

    -- ROLE_EXTERNAL
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ROLE_EXTERNAL') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ROLES', 'ROLE_EXTERNAL',
                '{"name": "Usuario Externo", "permissions": []}'::jsonb,
                'External user role');
    END IF;

END $$;
