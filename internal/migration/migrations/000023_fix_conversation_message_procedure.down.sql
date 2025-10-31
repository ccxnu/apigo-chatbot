-- Revert to original procedure (from 000006_whatsapp_procedures.up.sql)
DROP PROCEDURE IF EXISTS sp_create_conversation_message;

CREATE OR REPLACE PROCEDURE sp_create_conversation_message(
    OUT success BOOLEAN,
    OUT code VARCHAR,
    OUT o_cvm_id INT,
    IN p_conversation_id INT,
    IN p_message_id VARCHAR,
    IN p_from_me BOOLEAN,
    IN p_sender_name VARCHAR DEFAULT NULL,
    IN p_sender_type VARCHAR DEFAULT 'user',
    IN p_message_type VARCHAR DEFAULT 'text',
    IN p_body TEXT DEFAULT NULL,
    IN p_media_url VARCHAR DEFAULT NULL,
    IN p_quoted_message VARCHAR DEFAULT NULL,
    IN p_timestamp BIGINT DEFAULT NULL,
    IN p_is_forwarded BOOLEAN DEFAULT FALSE,
    IN p_queue_time_ms INT DEFAULT NULL,
    IN p_prompt_tokens INT DEFAULT NULL,
    IN p_prompt_time_ms INT DEFAULT NULL,
    IN p_completion_tokens INT DEFAULT NULL,
    IN p_completion_time_ms INT DEFAULT NULL,
    IN p_total_tokens INT DEFAULT NULL,
    IN p_total_time_ms INT DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_exists BOOLEAN;
BEGIN
    success := TRUE;
    code := 'OK';
    o_cvm_id := NULL;

    SELECT EXISTS(
        SELECT 1
        FROM public.cht_conversations
        WHERE cnv_id = p_conversation_id
    ) INTO v_exists;

    IF NOT v_exists THEN
        success := FALSE;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        RETURN;
    END IF;

    INSERT INTO public.cht_conversation_messages (
        cvm_fk_conversation,
        cvm_message_id,
        cvm_from_me,
        cvm_sender_name,
        cvm_sender_type,
        cvm_message_type,
        cvm_body,
        cvm_media_url,
        cvm_quoted_message,
        cvm_timestamp,
        cvm_is_forwarded,
        cvm_queue_time_ms,
        cvm_prompt_tokens,
        cvm_prompt_time_ms,
        cvm_completion_tokens,
        cvm_completion_time_ms,
        cvm_total_tokens,
        cvm_total_time_ms
    ) VALUES (
        p_conversation_id,
        p_message_id,
        p_from_me,
        p_sender_name,
        p_sender_type,
        p_message_type,
        p_body,
        p_media_url,
        p_quoted_message,
        COALESCE(p_timestamp, EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::BIGINT),
        p_is_forwarded,
        p_queue_time_ms,
        p_prompt_tokens,
        p_prompt_time_ms,
        p_completion_tokens,
        p_completion_time_ms,
        p_total_tokens,
        p_total_time_ms
    )
    RETURNING cvm_id INTO o_cvm_id;

    UPDATE public.cht_conversations
    SET
        cnv_message_count = cnv_message_count + 1,
        cnv_last_message_at = CURRENT_TIMESTAMP
    WHERE cnv_id = p_conversation_id;

EXCEPTION
    WHEN OTHERS THEN
        success := FALSE;
        code := 'ERR_CREATE_MESSAGE';
        o_cvm_id := NULL;
        RAISE NOTICE 'Error creating message: %', SQLERRM;
END;
$$;

COMMENT ON PROCEDURE sp_create_conversation_message IS 'Create new message in conversation. Returns success, code, and cvm_id';
