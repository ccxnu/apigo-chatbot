-- =====================================================
-- External User Usage Limits
-- =====================================================
-- This file adds configuration and tracking for limiting
-- external user access to the chatbot

DO $$
BEGIN
    -- =====================================================
    -- CONFIGURATION PARAMETERS
    -- =====================================================

    -- Enable/Disable External User Registration
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_REGISTRATION_ENABLED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_REGISTRATION_ENABLED',
                '{"enabled": true}'::jsonb,
                'Allow external users to register. Set to false to disable external registration.');
    END IF;

    -- Daily Message Limit for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_DAILY_LIMIT',
                '{"limit": 20, "period": "daily"}'::jsonb,
                'Maximum number of messages an external user can send per day. 0 = unlimited.');
    END IF;

    -- Weekly Message Limit for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_WEEKLY_LIMIT',
                '{"limit": 100, "period": "weekly"}'::jsonb,
                'Maximum number of messages an external user can send per week. 0 = unlimited.');
    END IF;

    -- Monthly Message Limit for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_MONTHLY_LIMIT',
                '{"limit": 300, "period": "monthly"}'::jsonb,
                'Maximum number of messages an external user can send per month. 0 = unlimited.');
    END IF;

    -- Total Lifetime Limit for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_TOTAL_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_TOTAL_LIMIT',
                '{"limit": 1000, "enabled": false}'::jsonb,
                'Maximum total messages an external user can send (lifetime). Set enabled=true to activate.');
    END IF;

    -- Auto-expire External Users After Inactivity
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_EXPIRY_DAYS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_EXPIRY_DAYS',
                '{"days": 30, "enabled": true}'::jsonb,
                'Deactivate external users after N days of inactivity. 0 or enabled=false to disable.');
    END IF;

    -- Require Approval for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_REQUIRE_APPROVAL') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_REQUIRE_APPROVAL',
                '{"required": false}'::jsonb,
                'Require admin approval before external users can use the chatbot. Set to true for manual approval.');
    END IF;

    -- Maximum External Users Allowed
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_MAX_COUNT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_MAX_COUNT',
                '{"limit": 500, "enabled": false}'::jsonb,
                'Maximum number of external users allowed in the system. Set enabled=true to activate.');
    END IF;

    -- Access Hours for External Users
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_ACCESS_HOURS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_ACCESS_HOURS',
                '{"enabled": false, "start_hour": 8, "end_hour": 18, "timezone": "America/Guayaquil"}'::jsonb,
                'Restrict external user access to specific hours. enabled=true to activate.');
    END IF;

    -- Rate Limiting (messages per minute)
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_RATE_LIMIT') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('CONFIGURATION', 'EXTERNAL_USER_RATE_LIMIT',
                '{"messages_per_minute": 5, "enabled": true}'::jsonb,
                'Maximum messages per minute for external users. Prevents spam.');
    END IF;

    -- =====================================================
    -- ERROR MESSAGES
    -- =====================================================

    -- Daily Limit Reached
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_DAILY_LIMIT_REACHED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_DAILY_LIMIT_REACHED',
                '{"message": "Has alcanzado tu límite diario de mensajes. Intenta mañana."}'::jsonb,
                'External user exceeded daily message limit');
    END IF;

    -- Weekly Limit Reached
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_WEEKLY_LIMIT_REACHED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_WEEKLY_LIMIT_REACHED',
                '{"message": "Has alcanzado tu límite semanal de mensajes. Intenta la próxima semana."}'::jsonb,
                'External user exceeded weekly message limit');
    END IF;

    -- Monthly Limit Reached
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_MONTHLY_LIMIT_REACHED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_MONTHLY_LIMIT_REACHED',
                '{"message": "Has alcanzado tu límite mensual de mensajes. Intenta el próximo mes."}'::jsonb,
                'External user exceeded monthly message limit');
    END IF;

    -- Total Limit Reached
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_TOTAL_LIMIT_REACHED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_TOTAL_LIMIT_REACHED',
                '{"message": "Has alcanzado el límite total de mensajes permitidos. Contacta al administrador."}'::jsonb,
                'External user exceeded total lifetime limit');
    END IF;

    -- Rate Limit Exceeded
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_RATE_LIMIT_EXCEEDED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_RATE_LIMIT_EXCEEDED',
                '{"message": "Estás enviando mensajes muy rápido. Por favor espera un momento."}'::jsonb,
                'External user sending messages too quickly');
    END IF;

    -- External Registration Disabled
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_EXTERNAL_REG_DISABLED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_EXTERNAL_REG_DISABLED',
                '{"message": "El registro de usuarios externos está temporalmente deshabilitado."}'::jsonb,
                'External user registration is disabled');
    END IF;

    -- Approval Required
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_APPROVAL_REQUIRED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_APPROVAL_REQUIRED',
                '{"message": "Tu cuenta requiere aprobación del administrador. Te notificaremos cuando esté aprobada."}'::jsonb,
                'External user account pending approval');
    END IF;

    -- Outside Access Hours
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_OUTSIDE_ACCESS_HOURS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_OUTSIDE_ACCESS_HOURS',
                '{"message": "Los usuarios externos solo pueden acceder durante el horario de atención."}'::jsonb,
                'External user tried to access outside allowed hours');
    END IF;

    -- Max External Users Reached
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_MAX_EXTERNAL_USERS') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_MAX_EXTERNAL_USERS',
                '{"message": "Se ha alcanzado el límite máximo de usuarios externos. Intenta más tarde."}'::jsonb,
                'Maximum number of external users reached');
    END IF;

    -- Account Expired
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'ERR_ACCOUNT_EXPIRED') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('ERROR_CODES', 'ERR_ACCOUNT_EXPIRED',
                '{"message": "Tu cuenta ha expirado por inactividad. Contacta al administrador para reactivarla."}'::jsonb,
                'External user account expired due to inactivity');
    END IF;

END $$;

-- =====================================================
-- Add approval status column to users table
-- =====================================================
DO $$
BEGIN
    -- Add approval status for external users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'cht_users'
        AND column_name = 'usr_approved'
    ) THEN
        ALTER TABLE cht_users
        ADD COLUMN usr_approved BOOLEAN DEFAULT TRUE;

        COMMENT ON COLUMN cht_users.usr_approved IS 'Whether user is approved (mainly for external users requiring approval)';
    END IF;

    -- Add last activity timestamp
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'cht_users'
        AND column_name = 'usr_last_activity_at'
    ) THEN
        ALTER TABLE cht_users
        ADD COLUMN usr_last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

        COMMENT ON COLUMN cht_users.usr_last_activity_at IS 'Last time user sent a message (for inactivity tracking)';
    END IF;

    -- Add message counter for external users
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'cht_users'
        AND column_name = 'usr_message_count'
    ) THEN
        ALTER TABLE cht_users
        ADD COLUMN usr_message_count INT DEFAULT 0;

        COMMENT ON COLUMN cht_users.usr_message_count IS 'Total number of messages sent by user (for limits)';
    END IF;
END $$;

-- =====================================================
-- Create index on new columns
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_users_approved ON cht_users(usr_approved) WHERE usr_rol = 'ROLE_EXTERNAL';
CREATE INDEX IF NOT EXISTS idx_users_last_activity ON cht_users(usr_last_activity_at) WHERE usr_rol = 'ROLE_EXTERNAL';
CREATE INDEX IF NOT EXISTS idx_users_role_active ON cht_users(usr_rol, usr_active);
