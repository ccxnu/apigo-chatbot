-- Revert fn_get_conversation_history to original version without cvm_sender_type
CREATE OR REPLACE FUNCTION fn_get_conversation_history(
    p_chat_id VARCHAR,
    p_limit INT DEFAULT 20
)
RETURNS TABLE (
    cvm_id INT,
    cvm_fk_conversation INT,
    cvm_message_id VARCHAR,
    cvm_from_me BOOLEAN,
    cvm_sender_name VARCHAR,
    cvm_message_type VARCHAR,
    cvm_body TEXT,
    cvm_media_url VARCHAR,
    cvm_quoted_message VARCHAR,
    cvm_timestamp BIGINT,
    cvm_is_forwarded BOOLEAN,
    cvm_metadata JSONB,
    cvm_created_at TIMESTAMP
) AS $$
DECLARE
    v_cnv_id INT;
BEGIN
    -- Get conversation ID
    SELECT cnv_id INTO v_cnv_id
    FROM public.cht_conversations
    WHERE cnv_chat_id = p_chat_id;

    IF v_cnv_id IS NULL THEN
        RETURN;
    END IF;

    RETURN QUERY
    SELECT
        m.cvm_id,
        m.cvm_fk_conversation,
        m.cvm_message_id,
        m.cvm_from_me,
        m.cvm_sender_name,
        m.cvm_message_type,
        m.cvm_body,
        m.cvm_media_url,
        m.cvm_quoted_message,
        m.cvm_timestamp,
        m.cvm_is_forwarded,
        m.cvm_metadata,
        m.cvm_created_at
    FROM public.cht_conversation_messages m
    WHERE m.cvm_fk_conversation = v_cnv_id
    ORDER BY m.cvm_timestamp DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION fn_get_conversation_history IS 'Get message history for a conversation';
