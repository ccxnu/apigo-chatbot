-- =====================================================
-- Admin Authentication System
-- Migration: 07_admin_authentication.sql
-- Purpose: Create tables and procedures for admin authentication
--          with JWT tokens (access + refresh) and custom claims
-- =====================================================

-- =====================================================
-- Table: cht_admin_users
-- Description: Administrative users with enhanced security
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_admin_users (
    adm_id              SERIAL PRIMARY KEY,
    adm_username        VARCHAR(50) UNIQUE NOT NULL,
    adm_email           VARCHAR(100) UNIQUE NOT NULL,
    adm_password_hash   TEXT NOT NULL,  -- bcrypt hash
    adm_name            VARCHAR(100) NOT NULL,
    adm_role            VARCHAR(50) NOT NULL REFERENCES cht_parameters(prm_code),
    adm_permissions     JSONB DEFAULT '[]'::JSONB,  -- Custom permissions array
    adm_claims          JSONB DEFAULT '{}'::JSONB,  -- Custom JWT claims
    adm_is_active       BOOLEAN NOT NULL DEFAULT true,
    adm_is_locked       BOOLEAN NOT NULL DEFAULT false,
    adm_failed_attempts INT NOT NULL DEFAULT 0,
    adm_last_login      TIMESTAMP,
    adm_last_login_ip   VARCHAR(45),  -- IPv6 compatible
    adm_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    adm_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_admin_users_username ON cht_admin_users(adm_username);
CREATE INDEX IF NOT EXISTS idx_admin_users_email ON cht_admin_users(adm_email);
CREATE INDEX IF NOT EXISTS idx_admin_users_role ON cht_admin_users(adm_role);

-- =====================================================
-- Table: cht_refresh_tokens
-- Description: Refresh tokens for JWT authentication
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_refresh_tokens (
    rft_id              SERIAL PRIMARY KEY,
    rft_admin_id        INT NOT NULL REFERENCES cht_admin_users(adm_id) ON DELETE CASCADE,
    rft_token           TEXT UNIQUE NOT NULL,  -- Refresh token (JWT)
    rft_token_family    UUID NOT NULL DEFAULT ex.uuid_generate_v4(),  -- For token rotation
    rft_user_agent      TEXT,
    rft_ip_address      VARCHAR(45),
    rft_expires_at      TIMESTAMP NOT NULL,
    rft_is_revoked      BOOLEAN NOT NULL DEFAULT false,
    rft_revoked_at      TIMESTAMP,
    rft_revoked_reason  VARCHAR(255),
    rft_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_admin_id ON cht_refresh_tokens(rft_admin_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON cht_refresh_tokens(rft_token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_family ON cht_refresh_tokens(rft_token_family);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON cht_refresh_tokens(rft_expires_at);

-- =====================================================
-- Table: cht_api_keys
-- Description: API keys for external integrations
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_api_keys (
    key_id              SERIAL PRIMARY KEY,
    key_name            VARCHAR(100) NOT NULL,
    key_value           TEXT UNIQUE NOT NULL,  -- API key (hashed)
    key_type            VARCHAR(50) NOT NULL DEFAULT 'external_api',  -- external_api, internal, webhook
    key_claims          JSONB DEFAULT '{}'::JSONB,  -- Custom claims for API
    key_rate_limit      INT DEFAULT 1000,  -- Requests per hour
    key_allowed_ips     JSONB DEFAULT '[]'::JSONB,  -- Whitelist IPs
    key_permissions     JSONB DEFAULT '[]'::JSONB,  -- Allowed endpoints/actions
    key_is_active       BOOLEAN NOT NULL DEFAULT true,
    key_expires_at      TIMESTAMP,
    key_last_used_at    TIMESTAMP,
    key_created_by      INT REFERENCES cht_admin_users(adm_id),
    key_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    key_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_api_keys_value ON cht_api_keys(key_value);
CREATE INDEX IF NOT EXISTS idx_api_keys_type ON cht_api_keys(key_type);
CREATE INDEX IF NOT EXISTS idx_api_keys_active ON cht_api_keys(key_is_active);

-- =====================================================
-- Table: cht_auth_logs
-- Description: Authentication audit trail
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_auth_logs (
    log_id              SERIAL PRIMARY KEY,
    log_user_id         INT REFERENCES cht_admin_users(adm_id) ON DELETE CASCADE,
    log_username        VARCHAR(50),
    log_action          VARCHAR(50) NOT NULL,  -- login, logout, refresh, failed_login, password_reset
    log_status          VARCHAR(20) NOT NULL,  -- success, failure
    log_ip_address      VARCHAR(45),
    log_user_agent      TEXT,
    log_details         JSONB DEFAULT '{}'::JSONB,
    log_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_auth_logs_user_id ON cht_auth_logs(log_user_id);
CREATE INDEX IF NOT EXISTS idx_auth_logs_action ON cht_auth_logs(log_action);
CREATE INDEX IF NOT EXISTS idx_auth_logs_created ON cht_auth_logs(log_created_at);

-- =====================================================
-- Triggers for updated_at
-- =====================================================
CREATE TRIGGER tr_cht_admin_users_updated
    BEFORE UPDATE ON cht_admin_users
    FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();

CREATE TRIGGER tr_cht_api_keys_updated
    BEFORE UPDATE ON cht_api_keys
    FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();

-- =====================================================
-- Stored Procedure: sp_create_admin_user
-- Description: Create a new admin user with hashed password
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_create_admin_user(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT admin_id INT,
    IN p_username VARCHAR(50),
    IN p_email VARCHAR(100),
    IN p_password_hash TEXT,
    IN p_name VARCHAR(100),
    IN p_role VARCHAR(50),
    IN p_permissions JSONB DEFAULT '[]'::JSONB,
    IN p_claims JSONB DEFAULT '{}'::JSONB
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Check if username exists
    IF EXISTS (SELECT 1 FROM cht_admin_users WHERE adm_username = p_username) THEN
        success := false;
        code := 'ERR_USERNAME_EXISTS';
        admin_id := NULL;
        RETURN;
    END IF;

    -- Check if email exists
    IF EXISTS (SELECT 1 FROM cht_admin_users WHERE adm_email = p_email) THEN
        success := false;
        code := 'ERR_EMAIL_EXISTS';
        admin_id := NULL;
        RETURN;
    END IF;

    -- Insert admin user
    INSERT INTO cht_admin_users (
        adm_username, adm_email, adm_password_hash, adm_name,
        adm_role, adm_permissions, adm_claims
    ) VALUES (
        p_username, p_email, p_password_hash, p_name,
        p_role, p_permissions, p_claims
    )
    RETURNING adm_id INTO admin_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_CREATE_ADMIN';
        RAISE NOTICE 'Error creating admin user: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_get_admin_by_username
-- Description: Get admin user by username
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_admin_by_username(
    p_username VARCHAR(50)
)
RETURNS TABLE (
    adm_id INT,
    adm_username VARCHAR(50),
    adm_email VARCHAR(100),
    adm_password_hash TEXT,
    adm_name VARCHAR(100),
    adm_role VARCHAR(50),
    adm_permissions JSONB,
    adm_claims JSONB,
    adm_is_active BOOLEAN,
    adm_is_locked BOOLEAN,
    adm_failed_attempts INT,
    adm_last_login TIMESTAMP,
    adm_last_login_ip VARCHAR(45),
    adm_created_at TIMESTAMP,
    adm_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.adm_id, a.adm_username, a.adm_email, a.adm_password_hash,
        a.adm_name, a.adm_role, a.adm_permissions, a.adm_claims,
        a.adm_is_active, a.adm_is_locked, a.adm_failed_attempts,
        a.adm_last_login, a.adm_last_login_ip, a.adm_created_at, a.adm_updated_at
    FROM cht_admin_users a
    WHERE a.adm_username = p_username;
END;
$$;

-- =====================================================
-- Function: fn_get_admin_by_id
-- Description: Get admin user by ID
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_admin_by_id(
    p_admin_id INT
)
RETURNS TABLE (
    adm_id INT,
    adm_username VARCHAR(50),
    adm_email VARCHAR(100),
    adm_password_hash TEXT,
    adm_name VARCHAR(100),
    adm_role VARCHAR(50),
    adm_permissions JSONB,
    adm_claims JSONB,
    adm_is_active BOOLEAN,
    adm_is_locked BOOLEAN,
    adm_failed_attempts INT,
    adm_last_login TIMESTAMP,
    adm_last_login_ip VARCHAR(45),
    adm_created_at TIMESTAMP,
    adm_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        a.adm_id, a.adm_username, a.adm_email, a.adm_password_hash,
        a.adm_name, a.adm_role, a.adm_permissions, a.adm_claims,
        a.adm_is_active, a.adm_is_locked, a.adm_failed_attempts,
        a.adm_last_login, a.adm_last_login_ip, a.adm_created_at, a.adm_updated_at
    FROM cht_admin_users a
    WHERE a.adm_id = p_admin_id;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_update_admin_login
-- Description: Update admin last login info
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_update_admin_login(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_admin_id INT,
    IN p_ip_address VARCHAR(45),
    IN p_reset_failed_attempts BOOLEAN DEFAULT true
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    IF p_reset_failed_attempts THEN
        UPDATE cht_admin_users
        SET adm_last_login = CURRENT_TIMESTAMP,
            adm_last_login_ip = p_ip_address,
            adm_failed_attempts = 0
        WHERE adm_id = p_admin_id;
    ELSE
        UPDATE cht_admin_users
        SET adm_last_login = CURRENT_TIMESTAMP,
            adm_last_login_ip = p_ip_address
        WHERE adm_id = p_admin_id;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_UPDATE_LOGIN';
        RAISE NOTICE 'Error updating admin login: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_increment_failed_attempts
-- Description: Increment failed login attempts and lock if needed
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_increment_failed_attempts(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT is_locked BOOLEAN,
    IN p_username VARCHAR(50)
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_attempts INT;
BEGIN
    success := true;
    code := 'OK';

    -- Increment failed attempts
    UPDATE cht_admin_users
    SET adm_failed_attempts = adm_failed_attempts + 1,
        adm_is_locked = CASE WHEN adm_failed_attempts + 1 >= 5 THEN true ELSE adm_is_locked END
    WHERE adm_username = p_username
    RETURNING adm_failed_attempts, adm_is_locked INTO v_attempts, is_locked;

    IF is_locked THEN
        success := false;
        code := 'ERR_ACCOUNT_LOCKED';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_INCREMENT_ATTEMPTS';
        RAISE NOTICE 'Error incrementing failed attempts: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_store_refresh_token
-- Description: Store a new refresh token
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_store_refresh_token(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT token_id INT,
    IN p_admin_id INT,
    IN p_token TEXT,
    IN p_token_family UUID,
    IN p_user_agent TEXT,
    IN p_ip_address VARCHAR(45),
    IN p_expires_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    INSERT INTO cht_refresh_tokens (
        rft_admin_id, rft_token, rft_token_family, rft_user_agent,
        rft_ip_address, rft_expires_at
    ) VALUES (
        p_admin_id, p_token, p_token_family, p_user_agent,
        p_ip_address, p_expires_at
    )
    RETURNING rft_id INTO token_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_STORE_TOKEN';
        RAISE NOTICE 'Error storing refresh token: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_get_refresh_token
-- Description: Get refresh token details
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_refresh_token(
    p_token TEXT
)
RETURNS TABLE (
    rft_id INT,
    rft_admin_id INT,
    rft_token TEXT,
    rft_token_family UUID,
    rft_expires_at TIMESTAMP,
    rft_is_revoked BOOLEAN,
    adm_username VARCHAR(50),
    adm_email VARCHAR(100),
    adm_name VARCHAR(100),
    adm_role VARCHAR(50),
    adm_permissions JSONB,
    adm_claims JSONB,
    adm_is_active BOOLEAN
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        rt.rft_id, rt.rft_admin_id, rt.rft_token, rt.rft_token_family,
        rt.rft_expires_at, rt.rft_is_revoked,
        a.adm_username, a.adm_email, a.adm_name, a.adm_role,
        a.adm_permissions, a.adm_claims, a.adm_is_active
    FROM cht_refresh_tokens rt
    INNER JOIN cht_admin_users a ON rt.rft_admin_id = a.adm_id
    WHERE rt.rft_token = p_token;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_revoke_refresh_token
-- Description: Revoke a refresh token
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_revoke_refresh_token(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_token TEXT,
    IN p_reason VARCHAR(255) DEFAULT 'user_logout'
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    UPDATE cht_refresh_tokens
    SET rft_is_revoked = true,
        rft_revoked_at = CURRENT_TIMESTAMP,
        rft_revoked_reason = p_reason
    WHERE rft_token = p_token;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_TOKEN_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_REVOKE_TOKEN';
        RAISE NOTICE 'Error revoking refresh token: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_revoke_token_family
-- Description: Revoke all tokens in a family (security breach)
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_revoke_token_family(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT revoked_count INT,
    IN p_token_family UUID,
    IN p_reason VARCHAR(255) DEFAULT 'security_breach'
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    UPDATE cht_refresh_tokens
    SET rft_is_revoked = true,
        rft_revoked_at = CURRENT_TIMESTAMP,
        rft_revoked_reason = p_reason
    WHERE rft_token_family = p_token_family
      AND rft_is_revoked = false;

    GET DIAGNOSTICS revoked_count = ROW_COUNT;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_REVOKE_FAMILY';
        RAISE NOTICE 'Error revoking token family: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_log_auth_event
-- Description: Log authentication events
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_log_auth_event(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_user_id INT,
    IN p_username VARCHAR(50),
    IN p_action VARCHAR(50),
    IN p_status VARCHAR(20),
    IN p_ip_address VARCHAR(45),
    IN p_user_agent TEXT,
    IN p_details JSONB DEFAULT '{}'::JSONB
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    INSERT INTO cht_auth_logs (
        log_user_id, log_username, log_action, log_status,
        log_ip_address, log_user_agent, log_details
    ) VALUES (
        p_user_id, p_username, p_action, p_status,
        p_ip_address, p_user_agent, p_details
    );

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_LOG_EVENT';
        RAISE NOTICE 'Error logging auth event: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_cleanup_expired_tokens
-- Description: Clean up expired refresh tokens (run periodically)
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_cleanup_expired_tokens(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT deleted_count INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    DELETE FROM cht_refresh_tokens
    WHERE rft_expires_at < CURRENT_TIMESTAMP
       OR (rft_is_revoked = true AND rft_revoked_at < CURRENT_TIMESTAMP - INTERVAL '30 days');

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_CLEANUP_TOKENS';
        RAISE NOTICE 'Error cleaning up tokens: %', SQLERRM;
END;
$$;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON TABLE cht_admin_users IS 'Administrative users with role-based access and custom claims';
COMMENT ON TABLE cht_refresh_tokens IS 'JWT refresh tokens with rotation support';
COMMENT ON TABLE cht_api_keys IS 'API keys for external integrations with custom claims';
COMMENT ON TABLE cht_auth_logs IS 'Authentication audit trail for security monitoring';

COMMENT ON COLUMN cht_admin_users.adm_claims IS 'Custom JWT claims for flexible authorization';
COMMENT ON COLUMN cht_admin_users.adm_permissions IS 'Custom permissions array for fine-grained access control';
COMMENT ON COLUMN cht_refresh_tokens.rft_token_family IS 'Token family ID for rotation and security breach detection';
COMMENT ON COLUMN cht_api_keys.key_claims IS 'Custom claims for API key (similar to JWT claims)';
