-- =====================================================
-- External User Limit Enforcement Procedures
-- =====================================================

-- =====================================================
-- Function: fn_check_external_user_limits
-- Description: Check if external user has exceeded any usage limits
-- Returns: allowed (boolean), code (varchar), remaining (jsonb)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_check_external_user_limits(
    p_user_id INT,
    p_whatsapp VARCHAR
)
RETURNS TABLE (
    allowed BOOLEAN,
    code VARCHAR,
    message TEXT,
    daily_remaining INT,
    weekly_remaining INT,
    monthly_remaining INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_user RECORD;
    v_daily_count INT;
    v_weekly_count INT;
    v_monthly_count INT;
    v_daily_limit INT;
    v_weekly_limit INT;
    v_monthly_limit INT;
    v_total_limit INT;
    v_total_limit_enabled BOOLEAN;
    v_approval_required BOOLEAN;
    v_recent_messages INT;
    v_rate_limit INT;
    v_rate_limit_enabled BOOLEAN;
BEGIN
    -- Get user information
    SELECT * INTO v_user
    FROM cht_users
    WHERE usr_id = p_user_id
    AND usr_active = TRUE;

    -- User not found or inactive
    IF NOT FOUND THEN
        RETURN QUERY SELECT FALSE, 'ERR_USER_NOT_FOUND'::VARCHAR, 'Usuario no encontrado'::TEXT, 0, 0, 0;
        RETURN;
    END IF;

    -- Only check limits for external users
    IF v_user.usr_rol != 'ROLE_EXTERNAL' THEN
        RETURN QUERY SELECT TRUE, 'OK'::VARCHAR, 'Sin límites'::TEXT, 999999, 999999, 999999;
        RETURN;
    END IF;

    -- Check if user is approved (if approval is required)
    SELECT (prm_data->>'required')::BOOLEAN INTO v_approval_required
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_REQUIRE_APPROVAL'
    AND prm_active = TRUE;

    IF v_approval_required AND NOT v_user.usr_approved THEN
        RETURN QUERY SELECT FALSE, 'ERR_APPROVAL_REQUIRED'::VARCHAR, 'Requiere aprobación'::TEXT, 0, 0, 0;
        RETURN;
    END IF;

    -- Get daily limit
    SELECT (prm_data->>'limit')::INT INTO v_daily_limit
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT'
    AND prm_active = TRUE;
    v_daily_limit := COALESCE(v_daily_limit, 0);

    -- Get weekly limit
    SELECT (prm_data->>'limit')::INT INTO v_weekly_limit
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT'
    AND prm_active = TRUE;
    v_weekly_limit := COALESCE(v_weekly_limit, 0);

    -- Get monthly limit
    SELECT (prm_data->>'limit')::INT INTO v_monthly_limit
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT'
    AND prm_active = TRUE;
    v_monthly_limit := COALESCE(v_monthly_limit, 0);

    -- Get total limit settings
    SELECT
        (prm_data->>'limit')::INT,
        COALESCE((prm_data->>'enabled')::BOOLEAN, FALSE)
    INTO v_total_limit, v_total_limit_enabled
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_TOTAL_LIMIT'
    AND prm_active = TRUE;

    -- Get rate limit settings
    SELECT
        (prm_data->>'messages_per_minute')::INT,
        COALESCE((prm_data->>'enabled')::BOOLEAN, FALSE)
    INTO v_rate_limit, v_rate_limit_enabled
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_RATE_LIMIT'
    AND prm_active = TRUE;

    -- Check rate limit (messages in last minute)
    IF v_rate_limit_enabled AND v_rate_limit > 0 THEN
        SELECT COUNT(*) INTO v_recent_messages
        FROM cht_conversation_messages m
        JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
        WHERE c.cnv_fk_user = p_user_id
        AND m.cvm_from_me = FALSE
        AND m.cvm_created_at > CURRENT_TIMESTAMP - INTERVAL '1 minute';

        IF v_recent_messages >= v_rate_limit THEN
            RETURN QUERY SELECT FALSE, 'ERR_RATE_LIMIT_EXCEEDED'::VARCHAR,
                'Demasiados mensajes'::TEXT, 0, 0, 0;
            RETURN;
        END IF;
    END IF;

    -- Check total lifetime limit
    IF v_total_limit_enabled AND v_total_limit > 0 THEN
        IF v_user.usr_message_count >= v_total_limit THEN
            RETURN QUERY SELECT FALSE, 'ERR_TOTAL_LIMIT_REACHED'::VARCHAR,
                'Límite total alcanzado'::TEXT, 0, 0, 0;
            RETURN;
        END IF;
    END IF;

    -- Count messages in different periods
    -- Daily count (today)
    SELECT COUNT(*) INTO v_daily_count
    FROM cht_conversation_messages m
    JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
    WHERE c.cnv_fk_user = p_user_id
    AND m.cvm_from_me = FALSE
    AND m.cvm_created_at >= CURRENT_DATE;

    -- Weekly count (last 7 days)
    SELECT COUNT(*) INTO v_weekly_count
    FROM cht_conversation_messages m
    JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
    WHERE c.cnv_fk_user = p_user_id
    AND m.cvm_from_me = FALSE
    AND m.cvm_created_at >= CURRENT_DATE - INTERVAL '7 days';

    -- Monthly count (current month)
    SELECT COUNT(*) INTO v_monthly_count
    FROM cht_conversation_messages m
    JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
    WHERE c.cnv_fk_user = p_user_id
    AND m.cvm_from_me = FALSE
    AND m.cvm_created_at >= DATE_TRUNC('month', CURRENT_DATE);

    -- Check daily limit
    IF v_daily_limit > 0 AND v_daily_count >= v_daily_limit THEN
        RETURN QUERY SELECT FALSE, 'ERR_DAILY_LIMIT_REACHED'::VARCHAR,
            'Límite diario alcanzado'::TEXT,
            0,
            GREATEST(0, v_weekly_limit - v_weekly_count),
            GREATEST(0, v_monthly_limit - v_monthly_count);
        RETURN;
    END IF;

    -- Check weekly limit
    IF v_weekly_limit > 0 AND v_weekly_count >= v_weekly_limit THEN
        RETURN QUERY SELECT FALSE, 'ERR_WEEKLY_LIMIT_REACHED'::VARCHAR,
            'Límite semanal alcanzado'::TEXT,
            GREATEST(0, v_daily_limit - v_daily_count),
            0,
            GREATEST(0, v_monthly_limit - v_monthly_count);
        RETURN;
    END IF;

    -- Check monthly limit
    IF v_monthly_limit > 0 AND v_monthly_count >= v_monthly_limit THEN
        RETURN QUERY SELECT FALSE, 'ERR_MONTHLY_LIMIT_REACHED'::VARCHAR,
            'Límite mensual alcanzado'::TEXT,
            GREATEST(0, v_daily_limit - v_daily_count),
            GREATEST(0, v_weekly_limit - v_weekly_count),
            0;
        RETURN;
    END IF;

    -- All checks passed - user is allowed
    RETURN QUERY SELECT
        TRUE,
        'OK'::VARCHAR,
        'Permitido'::TEXT,
        CASE WHEN v_daily_limit > 0 THEN GREATEST(0, v_daily_limit - v_daily_count) ELSE 999999 END,
        CASE WHEN v_weekly_limit > 0 THEN GREATEST(0, v_weekly_limit - v_weekly_count) ELSE 999999 END,
        CASE WHEN v_monthly_limit > 0 THEN GREATEST(0, v_monthly_limit - v_monthly_count) ELSE 999999 END;
END;
$$;

-- =====================================================
-- Procedure: sp_update_user_activity
-- Description: Update user's last activity and message count
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_update_user_activity(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    IN p_user_id INT
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := TRUE;
    code := 'OK';

    UPDATE cht_users
    SET usr_last_activity_at = CURRENT_TIMESTAMP,
        usr_message_count = usr_message_count + 1
    WHERE usr_id = p_user_id;

    IF NOT FOUND THEN
        success := FALSE;
        code := 'ERR_USER_NOT_FOUND';
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_UPDATE_ACTIVITY';
        RAISE NOTICE 'Error updating user activity: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_get_external_user_stats
-- Description: Get usage statistics for an external user
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_external_user_stats(
    p_user_id INT
)
RETURNS TABLE (
    daily_count INT,
    weekly_count INT,
    monthly_count INT,
    total_count INT,
    daily_limit INT,
    weekly_limit INT,
    monthly_limit INT,
    last_activity TIMESTAMP,
    approved BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_daily_limit INT;
    v_weekly_limit INT;
    v_monthly_limit INT;
    v_user RECORD;
BEGIN
    -- Get user info
    SELECT usr_message_count, usr_last_activity_at, usr_approved
    INTO v_user
    FROM cht_users
    WHERE usr_id = p_user_id;

    -- Get limits from parameters
    SELECT (prm_data->>'limit')::INT INTO v_daily_limit
    FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT' AND prm_active = TRUE;

    SELECT (prm_data->>'limit')::INT INTO v_weekly_limit
    FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT' AND prm_active = TRUE;

    SELECT (prm_data->>'limit')::INT INTO v_monthly_limit
    FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT' AND prm_active = TRUE;

    RETURN QUERY
    SELECT
        (SELECT COUNT(*)::INT
         FROM cht_conversation_messages m
         JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
         WHERE c.cnv_fk_user = p_user_id
         AND m.cvm_from_me = FALSE
         AND m.cvm_created_at >= CURRENT_DATE) AS daily_count,

        (SELECT COUNT(*)::INT
         FROM cht_conversation_messages m
         JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
         WHERE c.cnv_fk_user = p_user_id
         AND m.cvm_from_me = FALSE
         AND m.cvm_created_at >= CURRENT_DATE - INTERVAL '7 days') AS weekly_count,

        (SELECT COUNT(*)::INT
         FROM cht_conversation_messages m
         JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
         WHERE c.cnv_fk_user = p_user_id
         AND m.cvm_from_me = FALSE
         AND m.cvm_created_at >= DATE_TRUNC('month', CURRENT_DATE)) AS monthly_count,

        v_user.usr_message_count AS total_count,
        COALESCE(v_daily_limit, 0) AS daily_limit,
        COALESCE(v_weekly_limit, 0) AS weekly_limit,
        COALESCE(v_monthly_limit, 0) AS monthly_limit,
        v_user.usr_last_activity_at AS last_activity,
        v_user.usr_approved AS approved;
END;
$$;

-- =====================================================
-- Procedure: sp_deactivate_expired_external_users
-- Description: Deactivate external users who have been inactive
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_deactivate_expired_external_users(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT o_deactivated_count INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_expiry_days INT;
    v_expiry_enabled BOOLEAN;
BEGIN
    success := TRUE;
    code := 'OK';
    o_deactivated_count := 0;

    -- Get expiry settings
    SELECT
        (prm_data->>'days')::INT,
        COALESCE((prm_data->>'enabled')::BOOLEAN, FALSE)
    INTO v_expiry_days, v_expiry_enabled
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_EXPIRY_DAYS'
    AND prm_active = TRUE;

    -- Only proceed if expiry is enabled
    IF NOT v_expiry_enabled OR v_expiry_days IS NULL OR v_expiry_days = 0 THEN
        RETURN;
    END IF;

    -- Deactivate users
    UPDATE cht_users
    SET usr_active = FALSE
    WHERE usr_rol = 'ROLE_EXTERNAL'
    AND usr_active = TRUE
    AND usr_last_activity_at < CURRENT_TIMESTAMP - (v_expiry_days || ' days')::INTERVAL;

    GET DIAGNOSTICS o_deactivated_count = ROW_COUNT;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_DEACTIVATE_USERS';
        o_deactivated_count := 0;
        RAISE NOTICE 'Error deactivating expired users: %', SQLERRM;
END;
$$;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON FUNCTION fn_check_external_user_limits IS 'Check if external user has exceeded usage limits. Returns allowed status and remaining counts.';
COMMENT ON PROCEDURE sp_update_user_activity IS 'Update user last activity timestamp and increment message count';
COMMENT ON FUNCTION fn_get_external_user_stats IS 'Get detailed usage statistics for an external user';
COMMENT ON PROCEDURE sp_deactivate_expired_external_users IS 'Deactivate external users inactive for configured period';
