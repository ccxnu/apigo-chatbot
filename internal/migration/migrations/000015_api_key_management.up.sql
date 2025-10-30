-- =====================================================
-- API Key Management System
-- Migration: 000015_api_key_management.up.sql
-- Purpose: Create procedures for API key management and usage tracking
-- =====================================================

-- =====================================================
-- Table: cht_api_usage
-- Description: Track API usage for external API endpoints
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_api_usage (
    usg_id              SERIAL PRIMARY KEY,
    usg_api_key_id      INT NOT NULL REFERENCES cht_api_keys(key_id) ON DELETE CASCADE,
    usg_endpoint        VARCHAR(100) NOT NULL,
    usg_method          VARCHAR(10) NOT NULL,
    usg_status_code     INT NOT NULL,
    usg_tokens_used     INT DEFAULT 0,
    usg_request_time_ms INT DEFAULT 0,
    usg_ip_address      VARCHAR(45),
    usg_user_agent      TEXT,
    usg_error_message   TEXT,
    usg_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_usage_key_id ON cht_api_usage(usg_api_key_id);
CREATE INDEX IF NOT EXISTS idx_api_usage_created ON cht_api_usage(usg_created_at);
CREATE INDEX IF NOT EXISTS idx_api_usage_endpoint ON cht_api_usage(usg_endpoint);
CREATE INDEX IF NOT EXISTS idx_api_usage_status ON cht_api_usage(usg_status_code);

-- =====================================================
-- Stored Procedure: sp_create_api_key
-- Description: Create a new API key for external integrations
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_create_api_key(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT key_id INT,
    IN p_name VARCHAR(100),
    IN p_value TEXT,
    IN p_type VARCHAR(50) DEFAULT 'external_api',
    IN p_claims JSONB DEFAULT '{}'::JSONB,
    IN p_rate_limit INT DEFAULT 1000,
    IN p_allowed_ips JSONB DEFAULT '[]'::JSONB,
    IN p_permissions JSONB DEFAULT '[]'::JSONB,
    IN p_expires_at TIMESTAMP DEFAULT NULL,
    IN p_created_by INT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Check if API key value already exists
    IF EXISTS (SELECT 1 FROM cht_api_keys WHERE key_value = p_value) THEN
        success := false;
        code := 'ERR_API_KEY_EXISTS';
        key_id := NULL;
        RETURN;
    END IF;

    -- Insert API key
    INSERT INTO cht_api_keys (
        key_name, key_value, key_type, key_claims, key_rate_limit,
        key_allowed_ips, key_permissions, key_expires_at, key_created_by
    ) VALUES (
        p_name, p_value, p_type, p_claims, p_rate_limit,
        p_allowed_ips, p_permissions, p_expires_at, p_created_by
    )
    RETURNING key_id INTO key_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_CREATE_API_KEY';
        RAISE NOTICE 'Error creating API key: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_get_api_key_by_value
-- Description: Get API key details by value (for authentication)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_api_key_by_value(
    p_value TEXT
)
RETURNS TABLE (
    key_id INT,
    key_name VARCHAR(100),
    key_value TEXT,
    key_type VARCHAR(50),
    key_claims JSONB,
    key_rate_limit INT,
    key_allowed_ips JSONB,
    key_permissions JSONB,
    key_is_active BOOLEAN,
    key_expires_at TIMESTAMP,
    key_last_used_at TIMESTAMP,
    key_created_by INT,
    key_created_at TIMESTAMP,
    key_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        k.key_id, k.key_name, k.key_value, k.key_type, k.key_claims,
        k.key_rate_limit, k.key_allowed_ips, k.key_permissions,
        k.key_is_active, k.key_expires_at, k.key_last_used_at,
        k.key_created_by, k.key_created_at, k.key_updated_at
    FROM cht_api_keys k
    WHERE k.key_value = p_value;
END;
$$;

-- =====================================================
-- Function: fn_get_all_api_keys
-- Description: Get all API keys (for admin listing)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_all_api_keys()
RETURNS TABLE (
    key_id INT,
    key_name VARCHAR(100),
    key_value TEXT,
    key_type VARCHAR(50),
    key_claims JSONB,
    key_rate_limit INT,
    key_allowed_ips JSONB,
    key_permissions JSONB,
    key_is_active BOOLEAN,
    key_expires_at TIMESTAMP,
    key_last_used_at TIMESTAMP,
    key_created_by INT,
    key_created_at TIMESTAMP,
    key_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        k.key_id, k.key_name, k.key_value, k.key_type, k.key_claims,
        k.key_rate_limit, k.key_allowed_ips, k.key_permissions,
        k.key_is_active, k.key_expires_at, k.key_last_used_at,
        k.key_created_by, k.key_created_at, k.key_updated_at
    FROM cht_api_keys k
    ORDER BY k.key_created_at DESC;
END;
$$;

-- =====================================================
-- Function: fn_get_api_key_by_id
-- Description: Get API key details by ID
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_api_key_by_id(
    p_key_id INT
)
RETURNS TABLE (
    key_id INT,
    key_name VARCHAR(100),
    key_value TEXT,
    key_type VARCHAR(50),
    key_claims JSONB,
    key_rate_limit INT,
    key_allowed_ips JSONB,
    key_permissions JSONB,
    key_is_active BOOLEAN,
    key_expires_at TIMESTAMP,
    key_last_used_at TIMESTAMP,
    key_created_by INT,
    key_created_at TIMESTAMP,
    key_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        k.key_id, k.key_name, k.key_value, k.key_type, k.key_claims,
        k.key_rate_limit, k.key_allowed_ips, k.key_permissions,
        k.key_is_active, k.key_expires_at, k.key_last_used_at,
        k.key_created_by, k.key_created_at, k.key_updated_at
    FROM cht_api_keys k
    WHERE k.key_id = p_key_id;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_update_api_key_last_used
-- Description: Update last used timestamp for API key
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_update_api_key_last_used(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_key_id INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    UPDATE cht_api_keys
    SET key_last_used_at = CURRENT_TIMESTAMP
    WHERE key_id = p_key_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_API_KEY_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_UPDATE_API_KEY';
        RAISE NOTICE 'Error updating API key last used: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_update_api_key
-- Description: Update API key details
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_update_api_key(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_key_id INT,
    IN p_name VARCHAR(100) DEFAULT NULL,
    IN p_rate_limit INT DEFAULT NULL,
    IN p_allowed_ips JSONB DEFAULT NULL,
    IN p_permissions JSONB DEFAULT NULL,
    IN p_is_active BOOLEAN DEFAULT NULL,
    IN p_expires_at TIMESTAMP DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Update only non-null fields
    UPDATE cht_api_keys
    SET key_name = COALESCE(p_name, key_name),
        key_rate_limit = COALESCE(p_rate_limit, key_rate_limit),
        key_allowed_ips = COALESCE(p_allowed_ips, key_allowed_ips),
        key_permissions = COALESCE(p_permissions, key_permissions),
        key_is_active = COALESCE(p_is_active, key_is_active),
        key_expires_at = COALESCE(p_expires_at, key_expires_at)
    WHERE key_id = p_key_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_API_KEY_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_UPDATE_API_KEY';
        RAISE NOTICE 'Error updating API key: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_delete_api_key
-- Description: Delete (deactivate) an API key
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_delete_api_key(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_key_id INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Soft delete by setting is_active to false
    UPDATE cht_api_keys
    SET key_is_active = false
    WHERE key_id = p_key_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_API_KEY_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_DELETE_API_KEY';
        RAISE NOTICE 'Error deleting API key: %', SQLERRM;
END;
$$;

-- =====================================================
-- Stored Procedure: sp_track_api_usage
-- Description: Track API usage for analytics and billing
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_track_api_usage(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT usage_id INT,
    IN p_api_key_id INT,
    IN p_endpoint VARCHAR(100),
    IN p_method VARCHAR(10),
    IN p_status_code INT,
    IN p_tokens_used INT DEFAULT 0,
    IN p_request_time_ms INT DEFAULT 0,
    IN p_ip_address VARCHAR(45) DEFAULT NULL,
    IN p_user_agent TEXT DEFAULT NULL,
    IN p_error_message TEXT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    INSERT INTO cht_api_usage (
        usg_api_key_id, usg_endpoint, usg_method, usg_status_code,
        usg_tokens_used, usg_request_time_ms, usg_ip_address,
        usg_user_agent, usg_error_message
    ) VALUES (
        p_api_key_id, p_endpoint, p_method, p_status_code,
        p_tokens_used, p_request_time_ms, p_ip_address,
        p_user_agent, p_error_message
    )
    RETURNING usg_id INTO usage_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_TRACK_USAGE';
        RAISE NOTICE 'Error tracking API usage: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_get_api_usage_stats
-- Description: Get usage statistics for an API key
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_api_usage_stats(
    p_api_key_id INT,
    p_from_date TIMESTAMP DEFAULT NULL,
    p_to_date TIMESTAMP DEFAULT NULL
)
RETURNS TABLE (
    total_requests BIGINT,
    total_tokens BIGINT,
    avg_response_time NUMERIC,
    success_rate NUMERIC,
    requests_by_endpoint JSONB,
    requests_by_status JSONB
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_from_date TIMESTAMP;
    v_to_date TIMESTAMP;
BEGIN
    -- Default to last 30 days if not specified
    v_from_date := COALESCE(p_from_date, CURRENT_TIMESTAMP - INTERVAL '30 days');
    v_to_date := COALESCE(p_to_date, CURRENT_TIMESTAMP);

    RETURN QUERY
    SELECT
        COUNT(*)::BIGINT as total_requests,
        SUM(usg_tokens_used)::BIGINT as total_tokens,
        ROUND(AVG(usg_request_time_ms)::NUMERIC, 2) as avg_response_time,
        ROUND((COUNT(*) FILTER (WHERE usg_status_code < 400)::NUMERIC / COUNT(*)::NUMERIC * 100), 2) as success_rate,
        (SELECT jsonb_object_agg(usg_endpoint, count)
         FROM (
             SELECT usg_endpoint, COUNT(*) as count
             FROM cht_api_usage
             WHERE usg_api_key_id = p_api_key_id
               AND usg_created_at BETWEEN v_from_date AND v_to_date
             GROUP BY usg_endpoint
         ) endpoint_counts
        ) as requests_by_endpoint,
        (SELECT jsonb_object_agg(usg_status_code::TEXT, count)
         FROM (
             SELECT usg_status_code, COUNT(*) as count
             FROM cht_api_usage
             WHERE usg_api_key_id = p_api_key_id
               AND usg_created_at BETWEEN v_from_date AND v_to_date
             GROUP BY usg_status_code
         ) status_counts
        ) as requests_by_status
    FROM cht_api_usage
    WHERE usg_api_key_id = p_api_key_id
      AND usg_created_at BETWEEN v_from_date AND v_to_date;
END;
$$;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON TABLE cht_api_usage IS 'Tracks API usage for external API endpoints (analytics and billing)';
COMMENT ON COLUMN cht_api_usage.usg_tokens_used IS 'Number of tokens used in the request (for LLM calls)';
COMMENT ON COLUMN cht_api_usage.usg_request_time_ms IS 'Request processing time in milliseconds';
