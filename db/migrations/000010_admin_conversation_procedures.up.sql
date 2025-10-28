-- =====================================================
-- Admin Conversation Panel - Stored Procedures
-- =====================================================

-- =====================================================
-- Function: fn_get_all_conversations_for_admin
-- Description: Get all conversations with latest message preview (WhatsApp-like list)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_all_conversations_for_admin(
    p_limit int DEFAULT 50,
    p_offset int DEFAULT 0,
    p_filter varchar DEFAULT NULL -- 'unread', 'blocked', 'active', 'all'
)
RETURNS TABLE (
    cnv_id int,
    cnv_chat_id varchar,
    cnv_phone_number varchar,
    cnv_contact_name varchar,
    cnv_is_group boolean,
    cnv_group_name varchar,
    cnv_last_message_at timestamp,
    cnv_message_count int,
    cnv_unread_count int,
    cnv_blocked boolean,
    cnv_admin_intervened boolean,
    cnv_temporary boolean,
    cnv_expires_at timestamp,
    usr_id int,
    usr_name varchar,
    usr_identity_number varchar,
    usr_rol varchar,
    usr_blocked boolean,
    last_message_preview text,
    last_message_from_me boolean
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.cnv_id,
        c.cnv_chat_id,
        c.cnv_phone_number,
        c.cnv_contact_name,
        c.cnv_is_group,
        c.cnv_group_name,
        c.cnv_last_message_at,
        c.cnv_message_count,
        c.cnv_unread_count,
        c.cnv_blocked,
        c.cnv_admin_intervened,
        c.cnv_temporary,
        c.cnv_expires_at,
        u.usr_id,
        u.usr_name,
        u.usr_identity_number,
        u.usr_rol,
        u.usr_blocked,
        (SELECT cvm_body FROM cht_conversation_messages
         WHERE cvm_fk_conversation = c.cnv_id
         ORDER BY cvm_timestamp DESC LIMIT 1) as last_message_preview,
        (SELECT cvm_from_me FROM cht_conversation_messages
         WHERE cvm_fk_conversation = c.cnv_id
         ORDER BY cvm_timestamp DESC LIMIT 1) as last_message_from_me
    FROM public.cht_conversations c
    LEFT JOIN public.cht_users u ON c.cnv_fk_user = u.usr_id
    WHERE
        c.cnv_active = true
        AND (
            p_filter IS NULL
            OR (p_filter = 'unread' AND c.cnv_unread_count > 0)
            OR (p_filter = 'blocked' AND c.cnv_blocked = true)
            OR (p_filter = 'active' AND c.cnv_blocked = false AND c.cnv_unread_count > 0)
            OR (p_filter = 'all')
        )
    ORDER BY c.cnv_last_message_at DESC NULLS LAST
    LIMIT p_limit
    OFFSET p_offset;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Function: fn_get_conversation_messages
-- Description: Get all messages for a conversation
-- =====================================================
CREATE OR REPLACE FUNCTION fn_get_conversation_messages(
    p_conversation_id int,
    p_limit int DEFAULT 100
)
RETURNS TABLE (
    cvm_id int,
    cvm_message_id varchar,
    cvm_from_me boolean,
    cvm_sender_name varchar,
    cvm_sender_type varchar,
    cvm_message_type varchar,
    cvm_body text,
    cvm_media_url varchar,
    cvm_quoted_message varchar,
    cvm_timestamp bigint,
    cvm_is_forwarded boolean,
    cvm_read boolean,
    cvm_created_at timestamp,
    admin_name varchar
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        m.cvm_id,
        m.cvm_message_id,
        m.cvm_from_me,
        m.cvm_sender_name,
        m.cvm_sender_type,
        m.cvm_message_type,
        m.cvm_body,
        m.cvm_media_url,
        m.cvm_quoted_message,
        m.cvm_timestamp,
        m.cvm_is_forwarded,
        m.cvm_read,
        m.cvm_created_at,
        a.adm_username as admin_name
    FROM public.cht_conversation_messages m
    LEFT JOIN public.cht_admin_users a ON m.cvm_admin_id = a.adm_id
    WHERE m.cvm_fk_conversation = p_conversation_id
    ORDER BY m.cvm_timestamp ASC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Procedure: sp_block_user
-- Description: Block/unblock a user
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_block_user(
    OUT success boolean,
    OUT code varchar,
    IN p_user_id int,
    IN p_blocked boolean,
    IN p_admin_id int,
    IN p_reason text DEFAULT NULL
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Update user
    UPDATE public.cht_users
    SET
        usr_blocked = p_blocked,
        usr_blocked_at = CASE WHEN p_blocked THEN CURRENT_TIMESTAMP ELSE NULL END,
        usr_blocked_by = CASE WHEN p_blocked THEN p_admin_id ELSE NULL END,
        usr_block_reason = CASE WHEN p_blocked THEN p_reason ELSE NULL END,
        usr_updated_at = CURRENT_TIMESTAMP
    WHERE usr_id = p_user_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_USER_NOT_FOUND';
        RETURN;
    END IF;

    -- Also block all conversations from this user
    UPDATE public.cht_conversations
    SET
        cnv_blocked = p_blocked,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_fk_user = p_user_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_BLOCK_USER';
        RAISE NOTICE 'Error blocking user: %', SQLERRM;
END;
$$;

-- =====================================================
-- Procedure: sp_delete_conversation
-- Description: Soft delete a conversation
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_delete_conversation(
    OUT success boolean,
    OUT code varchar,
    IN p_conversation_id int
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    UPDATE public.cht_conversations
    SET
        cnv_active = false,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_id = p_conversation_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        RETURN;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_DELETE_CONVERSATION';
        RAISE NOTICE 'Error deleting conversation: %', SQLERRM;
END;
$$;

-- =====================================================
-- Procedure: sp_send_admin_message
-- Description: Admin sends a message in a conversation
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_send_admin_message(
    OUT success boolean,
    OUT code varchar,
    OUT o_message_id int,
    IN p_conversation_id int,
    IN p_admin_id int,
    IN p_message_id varchar,
    IN p_body text
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_exists boolean;
BEGIN
    success := true;
    code := 'OK';
    o_message_id := NULL;

    -- Check conversation exists
    SELECT EXISTS(
        SELECT 1 FROM public.cht_conversations WHERE cnv_id = p_conversation_id
    ) INTO v_exists;

    IF NOT v_exists THEN
        success := false;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        RETURN;
    END IF;

    -- Insert admin message
    INSERT INTO public.cht_conversation_messages (
        cvm_fk_conversation,
        cvm_message_id,
        cvm_from_me,
        cvm_sender_name,
        cvm_sender_type,
        cvm_message_type,
        cvm_body,
        cvm_timestamp,
        cvm_admin_id,
        cvm_read
    ) VALUES (
        p_conversation_id,
        p_message_id,
        true, -- from_me = true (we're sending)
        'Admin',
        'admin',
        'text',
        p_body,
        EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::bigint,
        p_admin_id,
        true -- admin messages are always "read"
    )
    RETURNING cvm_id INTO o_message_id;

    -- Update conversation
    UPDATE public.cht_conversations
    SET
        cnv_message_count = cnv_message_count + 1,
        cnv_last_message_at = CURRENT_TIMESTAMP,
        cnv_admin_intervened = true,
        cnv_last_admin_message_at = CURRENT_TIMESTAMP,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_id = p_conversation_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_SEND_ADMIN_MESSAGE';
        o_message_id := NULL;
        RAISE NOTICE 'Error sending admin message: %', SQLERRM;
END;
$$;

-- =====================================================
-- Procedure: sp_mark_messages_as_read
-- Description: Mark messages as read by admin
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_mark_messages_as_read(
    OUT success boolean,
    OUT code varchar,
    IN p_conversation_id int
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    -- Mark all unread messages as read
    UPDATE public.cht_conversation_messages
    SET
        cvm_read = true,
        cvm_read_at = CURRENT_TIMESTAMP
    WHERE cvm_fk_conversation = p_conversation_id
      AND cvm_read = false
      AND cvm_sender_type = 'user';

    -- Reset unread count
    UPDATE public.cht_conversations
    SET
        cnv_unread_count = 0,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_id = p_conversation_id;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_MARK_READ';
        RAISE NOTICE 'Error marking messages as read: %', SQLERRM;
END;
$$;

-- =====================================================
-- Procedure: sp_set_conversation_temporary
-- Description: Enable/disable temporary conversation with expiry
-- =====================================================
CREATE OR REPLACE PROCEDURE sp_set_conversation_temporary(
    OUT success boolean,
    OUT code varchar,
    IN p_conversation_id int,
    IN p_temporary boolean,
    IN p_hours_until_expiry int DEFAULT 24
)
LANGUAGE plpgsql
AS $$
BEGIN
    success := true;
    code := 'OK';

    UPDATE public.cht_conversations
    SET
        cnv_temporary = p_temporary,
        cnv_expires_at = CASE
            WHEN p_temporary THEN CURRENT_TIMESTAMP + (p_hours_until_expiry || ' hours')::INTERVAL
            ELSE NULL
        END,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_id = p_conversation_id;

    IF NOT FOUND THEN
        success := false;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        RETURN;
    END IF;

EXCEPTION
    WHEN OTHERS THEN
        success := false;
        code := 'ERR_SET_TEMPORARY';
        RAISE NOTICE 'Error setting temporary conversation: %', SQLERRM;
END;
$$;

-- =====================================================
-- Function: fn_cleanup_expired_conversations
-- Description: Delete expired temporary conversations (for cron job)
-- =====================================================
CREATE OR REPLACE FUNCTION fn_cleanup_expired_conversations()
RETURNS int AS $$
DECLARE
    v_deleted_count int;
BEGIN
    UPDATE public.cht_conversations
    SET
        cnv_active = false,
        cnv_updated_at = CURRENT_TIMESTAMP
    WHERE cnv_temporary = true
      AND cnv_expires_at < CURRENT_TIMESTAMP
      AND cnv_active = true;

    GET DIAGNOSTICS v_deleted_count = ROW_COUNT;

    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON FUNCTION fn_get_all_conversations_for_admin IS 'Get paginated list of conversations for admin panel';
COMMENT ON FUNCTION fn_get_conversation_messages IS 'Get all messages for a conversation';
COMMENT ON PROCEDURE sp_block_user IS 'Block or unblock a user from using the chatbot';
COMMENT ON PROCEDURE sp_delete_conversation IS 'Soft delete a conversation';
COMMENT ON PROCEDURE sp_send_admin_message IS 'Admin sends a message in a conversation';
COMMENT ON PROCEDURE sp_mark_messages_as_read IS 'Mark all messages in conversation as read';
COMMENT ON PROCEDURE sp_set_conversation_temporary IS 'Enable temporary conversation with auto-delete';
COMMENT ON FUNCTION fn_cleanup_expired_conversations IS 'Delete expired temporary conversations (run via cron)';
