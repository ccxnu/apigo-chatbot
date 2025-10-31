-- =====================================================
-- Analytics Stored Procedures
-- =====================================================
-- Purpose: Provide comprehensive analytics for admin dashboard
-- Created: 2025-10-22
-- =====================================================

-- =====================================================
-- 1. COST ANALYTICS
-- =====================================================
-- Calculate cost based on token usage and pricing from parameters
CREATE OR REPLACE FUNCTION fn_get_cost_analytics(
    p_start_date timestamp DEFAULT NULL,
    p_end_date timestamp DEFAULT NULL
)
RETURNS TABLE (
    period_start timestamp,
    period_end timestamp,
    total_cost numeric,
    llm_cost numeric,
    embedding_cost numeric,
    prompt_tokens bigint,
    completion_tokens bigint,
    total_tokens bigint,
    embedding_tokens bigint,
    conversation_count bigint,
    cost_per_conversation numeric,
    avg_tokens_per_conversation numeric
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_prompt_price numeric;
    v_completion_price numeric;
    v_embedding_price numeric;
    v_cost_config jsonb;
BEGIN
    -- Get pricing configuration from parameters
    SELECT prm_data INTO v_cost_config
    FROM cht_parameters
    WHERE prm_code = 'COST_CONFIG'
    LIMIT 1;

    -- Extract pricing (default values if config doesn't exist)
    v_prompt_price := COALESCE((v_cost_config->>'llmPricePerMillionInputTokens')::numeric, 0.50);
    v_completion_price := COALESCE((v_cost_config->>'llmPricePerMillionOutputTokens')::numeric, 1.50);
    v_embedding_price := COALESCE((v_cost_config->>'embeddingPricePerMillionTokens')::numeric, 0.13);

    -- Default date range if not provided (current month)
    IF p_start_date IS NULL THEN
        p_start_date := date_trunc('month', CURRENT_TIMESTAMP);
    END IF;
    IF p_end_date IS NULL THEN
        p_end_date := CURRENT_TIMESTAMP;
    END IF;

    RETURN QUERY
    SELECT
        p_start_date as period_start,
        p_end_date as period_end,
        -- Total cost
        COALESCE(ROUND(
            (COALESCE(SUM(cvm_prompt_tokens), 0)::numeric / 1000000.0 * v_prompt_price) +
            (COALESCE(SUM(cvm_completion_tokens), 0)::numeric / 1000000.0 * v_completion_price),
            2
        ), 0.0) as total_cost,
        -- LLM cost (prompt + completion)
        COALESCE(ROUND(
            (COALESCE(SUM(cvm_prompt_tokens), 0)::numeric / 1000000.0 * v_prompt_price) +
            (COALESCE(SUM(cvm_completion_tokens), 0)::numeric / 1000000.0 * v_completion_price),
            2
        ), 0.0) as llm_cost,
        -- Embedding cost (placeholder - will need actual embedding count)
        ROUND(0.0, 2) as embedding_cost,
        -- Token statistics
        COALESCE(SUM(cvm_prompt_tokens), 0)::bigint as prompt_tokens,
        COALESCE(SUM(cvm_completion_tokens), 0)::bigint as completion_tokens,
        COALESCE(SUM(cvm_total_tokens), 0)::bigint as total_tokens,
        0::bigint as embedding_tokens, -- Placeholder
        -- Conversation count
        COUNT(DISTINCT cvm_fk_conversation)::bigint as conversation_count,
        -- Cost per conversation
        CASE
            WHEN COUNT(DISTINCT cvm_fk_conversation) > 0 THEN
                ROUND(
                    (
                        (COALESCE(SUM(cvm_prompt_tokens), 0)::numeric / 1000000.0 * v_prompt_price) +
                        (COALESCE(SUM(cvm_completion_tokens), 0)::numeric / 1000000.0 * v_completion_price)
                    ) / COUNT(DISTINCT cvm_fk_conversation)::numeric,
                    4
                )
            ELSE 0.0
        END as cost_per_conversation,
        -- Average tokens per conversation
        CASE
            WHEN COUNT(DISTINCT cvm_fk_conversation) > 0 THEN
                ROUND(COALESCE(SUM(cvm_total_tokens), 0)::numeric / COUNT(DISTINCT cvm_fk_conversation)::numeric, 0)
            ELSE 0.0
        END as avg_tokens_per_conversation
    FROM cht_conversation_messages
    WHERE cvm_created_at BETWEEN p_start_date AND p_end_date
        AND cvm_sender_type = 'bot'; -- Only count bot responses (LLM usage)
END;
$$;

-- =====================================================
-- 2. TOKEN USAGE ANALYTICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_token_usage(
    p_period varchar DEFAULT 'month', -- 'day', 'week', 'month', 'year', 'all'
    p_group_by varchar DEFAULT 'day' -- 'hour', 'day', 'week'
)
RETURNS TABLE (
    period_label text,
    period_start timestamp,
    period_end timestamp,
    prompt_tokens bigint,
    completion_tokens bigint,
    total_tokens bigint,
    message_count bigint,
    conversation_count bigint,
    avg_tokens_per_message numeric
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    -- Determine date range based on period
    CASE p_period
        WHEN 'day' THEN
            v_start_date := date_trunc('day', v_end_date);
        WHEN 'week' THEN
            v_start_date := date_trunc('week', v_end_date);
        WHEN 'month' THEN
            v_start_date := date_trunc('month', v_end_date);
        WHEN 'year' THEN
            v_start_date := date_trunc('year', v_end_date);
        WHEN 'all' THEN
            v_start_date := '2020-01-01'::timestamp;
        ELSE
            v_start_date := date_trunc('month', v_end_date);
    END CASE;

    -- Group and aggregate based on grouping parameter
    IF p_group_by = 'hour' THEN
        RETURN QUERY
        SELECT
            to_char(date_trunc('hour', cvm_created_at), 'YYYY-MM-DD HH24:00') as period_label,
            date_trunc('hour', cvm_created_at) as period_start,
            date_trunc('hour', cvm_created_at) + interval '1 hour' as period_end,
            COALESCE(SUM(cvm_prompt_tokens), 0)::bigint as prompt_tokens,
            COALESCE(SUM(cvm_completion_tokens), 0)::bigint as completion_tokens,
            COALESCE(SUM(cvm_total_tokens), 0)::bigint as total_tokens,
            COUNT(*)::bigint as message_count,
            COUNT(DISTINCT cvm_fk_conversation)::bigint as conversation_count,
            COALESCE(ROUND(AVG(cvm_total_tokens), 0), 0.0) as avg_tokens_per_message
        FROM cht_conversation_messages
        WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
            AND cvm_sender_type = 'bot'
        GROUP BY date_trunc('hour', cvm_created_at)
        ORDER BY date_trunc('hour', cvm_created_at) DESC;
    ELSIF p_group_by = 'week' THEN
        RETURN QUERY
        SELECT
            to_char(date_trunc('week', cvm_created_at), 'YYYY "Week" IW') as period_label,
            date_trunc('week', cvm_created_at) as period_start,
            date_trunc('week', cvm_created_at) + interval '1 week' as period_end,
            COALESCE(SUM(cvm_prompt_tokens), 0)::bigint as prompt_tokens,
            COALESCE(SUM(cvm_completion_tokens), 0)::bigint as completion_tokens,
            COALESCE(SUM(cvm_total_tokens), 0)::bigint as total_tokens,
            COUNT(*)::bigint as message_count,
            COUNT(DISTINCT cvm_fk_conversation)::bigint as conversation_count,
            COALESCE(ROUND(AVG(cvm_total_tokens), 0), 0.0) as avg_tokens_per_message
        FROM cht_conversation_messages
        WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
            AND cvm_sender_type = 'bot'
        GROUP BY date_trunc('week', cvm_created_at)
        ORDER BY date_trunc('week', cvm_created_at) DESC;
    ELSE -- Default to 'day'
        RETURN QUERY
        SELECT
            to_char(date_trunc('day', cvm_created_at), 'YYYY-MM-DD') as period_label,
            date_trunc('day', cvm_created_at) as period_start,
            date_trunc('day', cvm_created_at) + interval '1 day' as period_end,
            COALESCE(SUM(cvm_prompt_tokens), 0)::bigint as prompt_tokens,
            COALESCE(SUM(cvm_completion_tokens), 0)::bigint as completion_tokens,
            COALESCE(SUM(cvm_total_tokens), 0)::bigint as total_tokens,
            COUNT(*)::bigint as message_count,
            COUNT(DISTINCT cvm_fk_conversation)::bigint as conversation_count,
            COALESCE(ROUND(AVG(cvm_total_tokens), 0), 0.0) as avg_tokens_per_message
        FROM cht_conversation_messages
        WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
            AND cvm_sender_type = 'bot'
        GROUP BY date_trunc('day', cvm_created_at)
        ORDER BY date_trunc('day', cvm_created_at) DESC;
    END IF;
END;
$$;

-- =====================================================
-- 3. ACTIVE USERS ANALYTICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_active_users(
    p_period varchar DEFAULT 'month' -- 'day', 'week', 'month', 'all'
)
RETURNS TABLE (
    period varchar,
    total_users bigint,
    active_users bigint,
    new_users bigint,
    returning_users bigint,
    students bigint,
    professors bigint,
    external bigint,
    avg_messages_per_user numeric,
    avg_sessions_per_user numeric
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
    v_prev_period_start timestamp;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    -- Determine date range
    CASE p_period
        WHEN 'day' THEN
            v_start_date := date_trunc('day', v_end_date);
            v_prev_period_start := v_start_date - interval '1 day';
        WHEN 'week' THEN
            v_start_date := date_trunc('week', v_end_date);
            v_prev_period_start := v_start_date - interval '1 week';
        WHEN 'month' THEN
            v_start_date := date_trunc('month', v_end_date);
            v_prev_period_start := v_start_date - interval '1 month';
        WHEN 'all' THEN
            v_start_date := '2020-01-01'::timestamp;
            v_prev_period_start := v_start_date;
        ELSE
            v_start_date := date_trunc('month', v_end_date);
            v_prev_period_start := v_start_date - interval '1 month';
    END CASE;

    RETURN QUERY
    WITH active_in_period AS (
        SELECT DISTINCT u.usr_id, u.usr_rol, u.usr_created_at
        FROM cht_users u
        INNER JOIN cht_conversations c ON c.cnv_fk_user = u.usr_id
        INNER JOIN cht_conversation_messages m ON m.cvm_fk_conversation = c.cnv_id
        WHERE m.cvm_created_at BETWEEN v_start_date AND v_end_date
    ),
    message_stats AS (
        SELECT
            c.cnv_fk_user,
            COUNT(m.cvm_id) as message_count,
            COUNT(DISTINCT c.cnv_id) as session_count
        FROM cht_conversations c
        INNER JOIN cht_conversation_messages m ON m.cvm_fk_conversation = c.cnv_id
        WHERE m.cvm_created_at BETWEEN v_start_date AND v_end_date
            AND m.cvm_sender_type = 'user'
        GROUP BY c.cnv_fk_user
    )
    SELECT
        p_period as period,
        -- Total users ever created
        (SELECT COUNT(*) FROM cht_users)::bigint as total_users,
        -- Active users in this period
        COUNT(DISTINCT ap.usr_id)::bigint as active_users,
        -- New users (created in this period)
        COUNT(DISTINCT ap.usr_id) FILTER (WHERE ap.usr_created_at >= v_start_date)::bigint as new_users,
        -- Returning users (created before this period)
        COUNT(DISTINCT ap.usr_id) FILTER (WHERE ap.usr_created_at < v_start_date)::bigint as returning_users,
        -- By role
        COUNT(DISTINCT ap.usr_id) FILTER (WHERE ap.usr_rol = 'ROLE_STUDENT')::bigint as students,
        COUNT(DISTINCT ap.usr_id) FILTER (WHERE ap.usr_rol = 'ROLE_PROFESSOR')::bigint as professors,
        COUNT(DISTINCT ap.usr_id) FILTER (WHERE ap.usr_rol = 'ROLE_EXTERNAL')::bigint as external,
        -- Engagement metrics
        COALESCE(ROUND(AVG(ms.message_count), 1), 0.0) as avg_messages_per_user,
        COALESCE(ROUND(AVG(ms.session_count), 1), 0.0) as avg_sessions_per_user
    FROM active_in_period ap
    LEFT JOIN message_stats ms ON ms.cnv_fk_user = ap.usr_id;
END;
$$;

-- =====================================================
-- 4. CONVERSATION METRICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_conversation_metrics(
    p_period varchar DEFAULT 'month'
)
RETURNS TABLE (
    period varchar,
    total_conversations bigint,
    active_conversations bigint,
    new_conversations bigint,
    avg_messages_per_conversation numeric,
    conversations_with_admin_help bigint,
    admin_intervention_rate numeric,
    blocked_conversations bigint,
    temporary_conversations bigint
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    CASE p_period
        WHEN 'day' THEN v_start_date := date_trunc('day', v_end_date);
        WHEN 'week' THEN v_start_date := date_trunc('week', v_end_date);
        WHEN 'month' THEN v_start_date := date_trunc('month', v_end_date);
        WHEN 'all' THEN v_start_date := '2020-01-01'::timestamp;
        ELSE v_start_date := date_trunc('month', v_end_date);
    END CASE;

    RETURN QUERY
    WITH conversation_stats AS (
        SELECT
            c.cnv_id,
            c.cnv_created_at,
            c.cnv_admin_intervened,
            c.cnv_blocked,
            c.cnv_temporary,
            COUNT(m.cvm_id) as message_count
        FROM cht_conversations c
        LEFT JOIN cht_conversation_messages m ON m.cvm_fk_conversation = c.cnv_id
        GROUP BY c.cnv_id
    ),
    period_conversations AS (
        SELECT *
        FROM conversation_stats
        WHERE cnv_created_at >= v_start_date
    )
    SELECT
        p_period as period,
        -- Total conversations ever
        (SELECT COUNT(*) FROM cht_conversations)::bigint as total_conversations,
        -- Active in period (have messages in period)
        COUNT(DISTINCT pc.cnv_id)::bigint as active_conversations,
        -- New conversations created in period
        COUNT(DISTINCT pc.cnv_id) FILTER (WHERE pc.cnv_created_at >= v_start_date)::bigint as new_conversations,
        -- Average messages
        COALESCE(ROUND(AVG(pc.message_count), 1), 0.0) as avg_messages_per_conversation,
        -- Admin intervention
        COUNT(DISTINCT pc.cnv_id) FILTER (WHERE pc.cnv_admin_intervened = true)::bigint as conversations_with_admin_help,
        CASE
            WHEN COUNT(DISTINCT pc.cnv_id) > 0 THEN
                ROUND(
                    COUNT(DISTINCT pc.cnv_id) FILTER (WHERE pc.cnv_admin_intervened = true)::numeric /
                    COUNT(DISTINCT pc.cnv_id)::numeric,
                    3
                )
            ELSE 0.0
        END as admin_intervention_rate,
        -- Status counts
        COUNT(DISTINCT pc.cnv_id) FILTER (WHERE pc.cnv_blocked = true)::bigint as blocked_conversations,
        COUNT(DISTINCT pc.cnv_id) FILTER (WHERE pc.cnv_temporary = true)::bigint as temporary_conversations
    FROM period_conversations pc;
END;
$$;

-- =====================================================
-- 5. MESSAGE ANALYTICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_message_analytics(
    p_period varchar DEFAULT 'month'
)
RETURNS TABLE (
    period varchar,
    total_messages bigint,
    user_messages bigint,
    bot_messages bigint,
    admin_messages bigint,
    avg_messages_per_day numeric,
    peak_hour int,
    peak_hour_count bigint
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
    v_days_in_period numeric;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    CASE p_period
        WHEN 'day' THEN
            v_start_date := date_trunc('day', v_end_date);
            v_days_in_period := 1;
        WHEN 'week' THEN
            v_start_date := date_trunc('week', v_end_date);
            v_days_in_period := 7;
        WHEN 'month' THEN
            v_start_date := date_trunc('month', v_end_date);
            v_days_in_period := EXTRACT(day FROM v_end_date - v_start_date);
        WHEN 'all' THEN
            v_start_date := '2020-01-01'::timestamp;
            v_days_in_period := EXTRACT(day FROM v_end_date - v_start_date);
        ELSE
            v_start_date := date_trunc('month', v_end_date);
            v_days_in_period := EXTRACT(day FROM v_end_date - v_start_date);
    END CASE;

    RETURN QUERY
    WITH message_counts AS (
        SELECT
            COUNT(*) as total,
            COUNT(*) FILTER (WHERE cvm_sender_type = 'user') as users,
            COUNT(*) FILTER (WHERE cvm_sender_type = 'bot') as bots,
            COUNT(*) FILTER (WHERE cvm_sender_type = 'admin') as admins
        FROM cht_conversation_messages
        WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
    ),
    hourly_distribution AS (
        SELECT
            EXTRACT(hour FROM cvm_created_at)::int as hour,
            COUNT(*) as count
        FROM cht_conversation_messages
        WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
        GROUP BY EXTRACT(hour FROM cvm_created_at)
        ORDER BY count DESC
        LIMIT 1
    )
    SELECT
        p_period as period,
        mc.total::bigint as total_messages,
        mc.users::bigint as user_messages,
        mc.bots::bigint as bot_messages,
        mc.admins::bigint as admin_messages,
        ROUND(mc.total::numeric / GREATEST(v_days_in_period, 1), 1) as avg_messages_per_day,
        COALESCE(hd.hour, 0) as peak_hour,
        COALESCE(hd.count, 0)::bigint as peak_hour_count
    FROM message_counts mc
    LEFT JOIN hourly_distribution hd ON true;
END;
$$;

-- =====================================================
-- 6. TOP QUERIES ANALYTICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_top_queries(
    p_period varchar DEFAULT 'month',
    p_limit int DEFAULT 20,
    p_min_similarity numeric DEFAULT 0.5
)
RETURNS TABLE (
    query_text text,
    query_count bigint,
    avg_similarity numeric,
    last_asked timestamp,
    has_good_answer boolean
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    CASE p_period
        WHEN 'day' THEN v_start_date := date_trunc('day', v_end_date);
        WHEN 'week' THEN v_start_date := date_trunc('week', v_end_date);
        WHEN 'month' THEN v_start_date := date_trunc('month', v_end_date);
        WHEN 'all' THEN v_start_date := '2020-01-01'::timestamp;
        ELSE v_start_date := date_trunc('month', v_end_date);
    END CASE;

    RETURN QUERY
    SELECT
        LEFT(cvm_body, 200) as query_text,
        COUNT(*)::bigint as query_count,
        COALESCE(ROUND(AVG(cvm_rag_best_similarity), 3), 0.0) as avg_similarity,
        MAX(cvm_created_at) as last_asked,
        COALESCE(AVG(cvm_rag_best_similarity), 0.0) >= p_min_similarity as has_good_answer
    FROM cht_conversation_messages
    WHERE cvm_created_at BETWEEN v_start_date AND v_end_date
        AND cvm_sender_type = 'user'
        AND cvm_body IS NOT NULL
        AND LENGTH(cvm_body) > 10 -- Filter out very short messages
    GROUP BY LEFT(cvm_body, 200)
    ORDER BY query_count DESC
    LIMIT p_limit;
END;
$$;

-- =====================================================
-- 7. KNOWLEDGE BASE USAGE ANALYTICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_knowledge_usage(
    p_period varchar DEFAULT 'month'
)
RETURNS TABLE (
    chunk_id int,
    document_title varchar,
    usage_count bigint,
    avg_similarity numeric,
    last_used timestamp
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_start_date timestamp;
    v_end_date timestamp;
BEGIN
    v_end_date := CURRENT_TIMESTAMP;

    CASE p_period
        WHEN 'day' THEN v_start_date := date_trunc('day', v_end_date);
        WHEN 'week' THEN v_start_date := date_trunc('week', v_end_date);
        WHEN 'month' THEN v_start_date := date_trunc('month', v_end_date);
        WHEN 'all' THEN v_start_date := '2020-01-01'::timestamp;
        ELSE v_start_date := date_trunc('month', v_end_date);
    END CASE;

    RETURN QUERY
    SELECT
        cs.chs_chk_id as chunk_id,
        d.dct_title as document_title,
        cs.chs_access_count::bigint as usage_count,
        COALESCE(ROUND(cs.chs_avg_similarity, 3), 0.0) as avg_similarity,
        cs.chs_last_accessed as last_used
    FROM cht_chunk_statistics cs
    INNER JOIN cht_chunks chk ON chk.chk_id = cs.chs_chk_id
    INNER JOIN cht_documents d ON d.dct_id = chk.chk_dct_id
    WHERE cs.chs_last_accessed BETWEEN v_start_date AND v_end_date
    ORDER BY cs.chs_access_count DESC
    LIMIT 50;
END;
$$;

-- =====================================================
-- 8. SYSTEM HEALTH METRICS
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_system_health()
RETURNS TABLE (
    metric_name varchar,
    metric_value numeric,
    metric_unit varchar
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    WITH response_times AS (
        SELECT
            COALESCE(AVG(cvm_prompt_time_ms + cvm_completion_time_ms), 0.0)::numeric as avg_llm_time,
            COALESCE(PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY cvm_prompt_time_ms + cvm_completion_time_ms), 0.0)::numeric as p95_llm_time,
            COALESCE(PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY cvm_prompt_time_ms + cvm_completion_time_ms), 0.0)::numeric as p99_llm_time
        FROM cht_conversation_messages
        WHERE cvm_created_at >= CURRENT_TIMESTAMP - interval '24 hours'
            AND cvm_sender_type = 'bot'
            AND cvm_prompt_time_ms IS NOT NULL
    ),
    error_count AS (
        SELECT COUNT(*) as errors
        FROM cht_conversation_messages
        WHERE cvm_created_at >= CURRENT_TIMESTAMP - interval '24 hours'
            AND cvm_body LIKE '%error%'
    )
    SELECT 'avg_llm_response_time'::varchar, COALESCE(ROUND(rt.avg_llm_time, 0), 0.0), 'ms'::varchar FROM response_times rt
    UNION ALL
    SELECT 'p95_llm_response_time'::varchar, COALESCE(ROUND(rt.p95_llm_time, 0), 0.0), 'ms'::varchar FROM response_times rt
    UNION ALL
    SELECT 'p99_llm_response_time'::varchar, COALESCE(ROUND(rt.p99_llm_time, 0), 0.0), 'ms'::varchar FROM response_times rt
    UNION ALL
    SELECT 'errors_last_24h'::varchar, ec.errors::numeric, 'count'::varchar FROM error_count ec
    UNION ALL
    SELECT 'total_conversations'::varchar, COUNT(*)::numeric, 'count'::varchar FROM cht_conversations
    UNION ALL
    SELECT 'total_users'::varchar, COUNT(*)::numeric, 'count'::varchar FROM cht_users;
END;
$$;

-- =====================================================
-- 9. DASHBOARD OVERVIEW (Combined Metrics)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_analytics_overview()
RETURNS TABLE (
    cost_this_month numeric,
    tokens_this_month bigint,
    active_users_today bigint,
    conversations_this_month bigint,
    messages_today bigint,
    avg_response_time_ms numeric,
    admin_intervention_rate numeric,
    last_updated timestamp
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_month_start timestamp := date_trunc('month', CURRENT_TIMESTAMP);
    v_today_start timestamp := date_trunc('day', CURRENT_TIMESTAMP);
    v_now timestamp := CURRENT_TIMESTAMP;
BEGIN
    RETURN QUERY
    SELECT
        -- Cost this month
        COALESCE((SELECT ca.total_cost FROM fn_get_cost_analytics(v_month_start, v_now) ca), 0.0) as cost_this_month,
        -- Tokens this month
        COALESCE((SELECT SUM(tu.total_tokens)::bigint FROM fn_get_token_usage('month', 'day') tu), 0) as tokens_this_month,
        -- Active users today
        COALESCE((SELECT au.active_users FROM fn_get_active_users('day') au), 0) as active_users_today,
        -- Conversations this month
        COALESCE((SELECT cm.new_conversations FROM fn_get_conversation_metrics('month') cm), 0) as conversations_this_month,
        -- Messages today
        COALESCE((SELECT ma.total_messages FROM fn_get_message_analytics('day') ma), 0) as messages_today,
        -- Average response time
        COALESCE((SELECT sh.metric_value FROM fn_get_system_health() sh WHERE sh.metric_name = 'avg_llm_response_time'), 0.0) as avg_response_time_ms,
        -- Admin intervention rate
        COALESCE((SELECT cm2.admin_intervention_rate FROM fn_get_conversation_metrics('month') cm2), 0.0) as admin_intervention_rate,
        -- Last updated
        v_now as last_updated;
END;
$$;

-- =====================================================
-- Grant permissions
-- =====================================================
GRANT EXECUTE ON FUNCTION fn_get_cost_analytics TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_token_usage TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_active_users TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_conversation_metrics TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_message_analytics TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_top_queries TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_knowledge_usage TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_system_health TO postgres;
GRANT EXECUTE ON FUNCTION fn_get_analytics_overview TO postgres;
