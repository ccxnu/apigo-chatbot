-- WhatsApp Module - SQL Functions & Procedures
-- Handles WhatsApp sessions, conversations, and messages

-- =====================================================
-- WHATSAPP SESSIONS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_whatsapp_session
-- Description: Get WhatsApp session by name
-- =====================================================
create or replace function fn_get_whatsapp_session(
    p_session_name varchar
)
returns table (
    wss_id int,
    wss_session_name varchar,
    wss_phone_number varchar,
    wss_device_name varchar,
    wss_platform varchar,
    wss_qr_code text,
    wss_connected boolean,
    wss_last_seen timestamp,
    wss_session_data jsonb,
    wss_active boolean,
    wss_created_at timestamp,
    wss_updated_at timestamp
) as $$
begin
    return query
    select
        ws.wss_id,
        ws.wss_session_name,
        ws.wss_phone_number,
        ws.wss_device_name,
        ws.wss_platform,
        ws.wss_qr_code,
        ws.wss_connected,
        ws.wss_last_seen,
        ws.wss_session_data,
        ws.wss_active,
        ws.wss_created_at,
        ws.wss_updated_at
    from public.cht_whatsapp_sessions ws
    where ws.wss_session_name = p_session_name
    and ws.wss_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_active_whatsapp_session
-- Description: Get the currently active WhatsApp session
-- =====================================================
create or replace function fn_get_active_whatsapp_session()
returns table (
    wss_id int,
    wss_session_name varchar,
    wss_phone_number varchar,
    wss_device_name varchar,
    wss_platform varchar,
    wss_qr_code text,
    wss_connected boolean,
    wss_last_seen timestamp,
    wss_session_data jsonb,
    wss_active boolean,
    wss_created_at timestamp,
    wss_updated_at timestamp
) as $$
begin
    return query
    select
        ws.wss_id,
        ws.wss_session_name,
        ws.wss_phone_number,
        ws.wss_device_name,
        ws.wss_platform,
        ws.wss_qr_code,
        ws.wss_connected,
        ws.wss_last_seen,
        ws.wss_session_data,
        ws.wss_active,
        ws.wss_created_at,
        ws.wss_updated_at
    from public.cht_whatsapp_sessions ws
    where ws.wss_active = true
    and ws.wss_connected = true
    order by ws.wss_last_seen desc nulls last
    limit 1;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_update_whatsapp_session_status
-- Description: Updates WhatsApp session connection status
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_whatsapp_session_status(
    out success boolean,
    out code varchar,
    in p_session_name varchar,
    in p_phone_number varchar default null,
    in p_device_name varchar default null,
    in p_platform varchar default null,
    in p_connected boolean default null
)
language plpgsql
as $$
declare
    v_session_id int;
begin
    success := true;
    code := 'OK';

    -- Get or create session
    select wss_id into v_session_id
    from public.cht_whatsapp_sessions
    where wss_session_name = p_session_name;

    if v_session_id is null then
        -- Create new session
        insert into public.cht_whatsapp_sessions (
            wss_session_name,
            wss_phone_number,
            wss_device_name,
            wss_platform,
            wss_connected
        ) values (
            p_session_name,
            p_phone_number,
            p_device_name,
            p_platform,
            coalesce(p_connected, false)
        )
        returning wss_id into v_session_id;
    else
        -- Update existing session
        update public.cht_whatsapp_sessions
        set
            wss_phone_number = coalesce(p_phone_number, wss_phone_number),
            wss_device_name = coalesce(p_device_name, wss_device_name),
            wss_platform = coalesce(p_platform, wss_platform),
            wss_connected = coalesce(p_connected, wss_connected),
            wss_last_seen = case when p_connected = true then current_timestamp else wss_last_seen end
        where wss_id = v_session_id;
    end if;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_WHATSAPP_SESSION';
        raise notice 'Error updating WhatsApp session: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_whatsapp_qr_code
-- Description: Updates QR code for WhatsApp session
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_whatsapp_qr_code(
    out success boolean,
    out code varchar,
    in p_session_name varchar,
    in p_qr_code text
)
language plpgsql
as $$
begin
    success := true;
    code := 'OK';

    update public.cht_whatsapp_sessions
    set wss_qr_code = p_qr_code
    where wss_session_name = p_session_name;

    if not found then
        -- Create session if doesn't exist
        insert into public.cht_whatsapp_sessions (
            wss_session_name,
            wss_qr_code,
            wss_connected
        ) values (
            p_session_name,
            p_qr_code,
            false
        );
    end if;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_QR_CODE';
        raise notice 'Error updating QR code: %', sqlerrm;
end;
$$;

-- =====================================================
-- CONVERSATIONS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_conversation_by_chat_id
-- Description: Get conversation by WhatsApp chat ID
-- =====================================================
create or replace function fn_get_conversation_by_chat_id(
    p_chat_id varchar
)
returns table (
    cnv_id int,
    cnv_fk_user int,
    cnv_chat_id varchar,
    cnv_phone_number varchar,
    cnv_contact_name varchar,
    cnv_is_group boolean,
    cnv_group_name varchar,
    cnv_last_message_at timestamp,
    cnv_message_count int,
    cnv_active boolean,
    cnv_created_at timestamp,
    cnv_updated_at timestamp
) as $$
begin
    return query
    select
        c.cnv_id,
        c.cnv_fk_user,
        c.cnv_chat_id,
        c.cnv_phone_number,
        c.cnv_contact_name,
        c.cnv_is_group,
        c.cnv_group_name,
        c.cnv_last_message_at,
        c.cnv_message_count,
        c.cnv_active,
        c.cnv_created_at,
        c.cnv_updated_at
    from public.cht_conversations c
    where c.cnv_chat_id = p_chat_id
    and c.cnv_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_conversation
-- Description: Creates or gets existing conversation
-- Returns: success (boolean), code (varchar), cnv_id (int)
-- =====================================================
create or replace procedure sp_create_conversation(
    out success boolean,
    out code varchar,
    out o_cnv_id int,
    in p_chat_id varchar,
    in p_phone_number varchar,
    in p_contact_name varchar default null,
    in p_is_group boolean default false,
    in p_group_name varchar default null
)
language plpgsql
as $$
begin
    success := true;
    code := 'OK';
    o_cnv_id := null;

    -- Check if conversation already exists
    select cnv_id into o_cnv_id
    from public.cht_conversations
    where cnv_chat_id = p_chat_id;

    if o_cnv_id is null then
        -- Create new conversation
        insert into public.cht_conversations (
            cnv_chat_id,
            cnv_phone_number,
            cnv_contact_name,
            cnv_is_group,
            cnv_group_name,
            cnv_message_count
        ) values (
            p_chat_id,
            p_phone_number,
            p_contact_name,
            p_is_group,
            p_group_name,
            0
        )
        returning cnv_id into o_cnv_id;
    end if;

exception
    when others then
        success := false;
        code := 'ERR_CREATE_CONVERSATION';
        o_cnv_id := null;
        raise notice 'Error creating conversation: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_link_user_to_conversation
-- Description: Links a validated user to a conversation
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_link_user_to_conversation(
    out success boolean,
    out code varchar,
    in p_chat_id varchar,
    in p_identity_number varchar
)
language plpgsql
as $$
declare
    v_user_id int;
    v_cnv_id int;
begin
    success := true;
    code := 'OK';

    -- Get user by identity number
    select usr_id into v_user_id
    from public.cht_users
    where usr_identity_number = p_identity_number
    and usr_active = true;

    if v_user_id is null then
        success := false;
        code := 'ERR_USER_NOT_FOUND';
        return;
    end if;

    -- Get conversation
    select cnv_id into v_cnv_id
    from public.cht_conversations
    where cnv_chat_id = p_chat_id;

    if v_cnv_id is null then
        success := false;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        return;
    end if;

    -- Link user to conversation
    update public.cht_conversations
    set cnv_fk_user = v_user_id
    where cnv_id = v_cnv_id;

exception
    when others then
        success := false;
        code := 'ERR_LINK_USER_CONVERSATION';
        raise notice 'Error linking user to conversation: %', sqlerrm;
end;
$$;

-- =====================================================
-- MESSAGES SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_conversation_history
-- Description: Get message history for a conversation
-- =====================================================
create or replace function fn_get_conversation_history(
    p_chat_id varchar,
    p_limit int default 20
)
returns table (
    cvm_id int,
    cvm_fk_conversation int,
    cvm_message_id varchar,
    cvm_from_me boolean,
    cvm_sender_name varchar,
    cvm_message_type varchar,
    cvm_body text,
    cvm_media_url varchar,
    cvm_quoted_message varchar,
    cvm_timestamp bigint,
    cvm_is_forwarded boolean,
    cvm_metadata jsonb,
    cvm_created_at timestamp
) as $$
declare
    v_cnv_id int;
begin
    -- Get conversation ID
    select cnv_id into v_cnv_id
    from public.cht_conversations
    where cnv_chat_id = p_chat_id;

    if v_cnv_id is null then
        return;
    end if;

    return query
    select
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
    from public.cht_conversation_messages m
    where m.cvm_fk_conversation = v_cnv_id
    order by m.cvm_timestamp desc
    limit p_limit;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_conversation_message
-- Description: Creates a new message in a conversation
-- Returns: success (boolean), code (varchar), cvm_id (int)
-- =====================================================
create or replace procedure sp_create_conversation_message(
    out success boolean,
    out code varchar,
    out o_cvm_id int,
    in p_conversation_id int,
    in p_message_id varchar,
    in p_from_me boolean,
    in p_sender_name varchar default null,
    in p_sender_type varchar default 'user',
    in p_message_type varchar default 'text',
    in p_body text default null,
    in p_media_url varchar default null,
    in p_quoted_message varchar default null,
    in p_timestamp bigint default null,
    in p_is_forwarded boolean default false,
    in p_queue_time_ms int default null,
    in p_prompt_tokens int default null,
    in p_prompt_time_ms int default null,
    in p_completion_tokens int default null,
    in p_completion_time_ms int default null,
    in p_total_tokens int default null,
    in p_total_time_ms int default null
)
language plpgsql
as $$
declare
    v_exists boolean;
begin
    success := true;
    code := 'OK';
    o_cvm_id := null;

    select exists(
        select 1
        from public.cht_conversations
        where cnv_id = p_conversation_id
    ) into v_exists;

    if not v_exists then
        success := false;
        code := 'ERR_CONVERSATION_NOT_FOUND';
        return;
    end if;

    insert into public.cht_conversation_messages (
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
    ) values (
        p_conversation_id,
        p_message_id,
        p_from_me,
        p_sender_name,
        p_sender_type,
        p_message_type,
        p_body,
        p_media_url,
        p_quoted_message,
        coalesce(p_timestamp, extract(epoch from current_timestamp)::bigint),
        p_is_forwarded,
        p_queue_time_ms,
        p_prompt_tokens,
        p_prompt_time_ms,
        p_completion_tokens,
        p_completion_time_ms,
        p_total_tokens,
        p_total_time_ms
    )
    returning cvm_id into o_cvm_id;

    -- Update conversation message count and last message time
    update public.cht_conversations
    set
        cnv_message_count = cnv_message_count + 1,
        cnv_last_message_at = current_timestamp
    where cnv_id = p_conversation_id;

exception
    when others then
        success := false;
        code := 'ERR_CREATE_MESSAGE';
        o_cvm_id := null;
        raise notice 'Error creating message: %', sqlerrm;
end;
$$;

-- =====================================================
-- USERS SECTION
-- =====================================================

-- =====================================================
-- Function: fn_get_user_by_identity
-- Description: Get user by identity number
-- =====================================================
create or replace function fn_get_user_by_identity(
    p_identity_number varchar
)
returns table (
    usr_id int,
    usr_identity_number varchar,
    usr_name varchar,
    usr_email varchar,
    usr_phone varchar,
    usr_rol varchar,
    usr_details jsonb,
    usr_whatsapp varchar,
    usr_active boolean,
    usr_created_at timestamp,
    usr_updated_at timestamp
) as $$
begin
    return query
    select
        u.usr_id,
        u.usr_identity_number,
        u.usr_name,
        u.usr_email,
        u.usr_phone,
        u.usr_rol,
        u.usr_details,
        u.usr_whatsapp,
        u.usr_active,
        u.usr_created_at,
        u.usr_updated_at
    from public.cht_users u
    where u.usr_identity_number = p_identity_number
    and u.usr_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Function: fn_get_user_by_whatsapp
-- Description: Get user by WhatsApp number
-- =====================================================
create or replace function fn_get_user_by_whatsapp(
    p_whatsapp varchar
)
returns table (
    usr_id int,
    usr_identity_number varchar,
    usr_name varchar,
    usr_email varchar,
    usr_phone varchar,
    usr_rol varchar,
    usr_details jsonb,
    usr_whatsapp varchar,
    usr_active boolean,
    usr_created_at timestamp,
    usr_updated_at timestamp
) as $$
begin
    return query
    select
        u.usr_id,
        u.usr_identity_number,
        u.usr_name,
        u.usr_email,
        u.usr_phone,
        u.usr_rol,
        u.usr_details,
        u.usr_whatsapp,
        u.usr_active,
        u.usr_created_at,
        u.usr_updated_at
    from public.cht_users u
    where u.usr_whatsapp = p_whatsapp
    and u.usr_active = true;
end;
$$ language plpgsql;

-- =====================================================
-- Procedure: sp_create_user
-- Description: Creates a new user from WhatsApp registration
-- Returns: success (boolean), code (varchar), usr_id (int)
-- =====================================================
create or replace procedure sp_create_user(
    out success boolean,
    out code varchar,
    out o_usr_id int,
    in p_identity_number varchar,
    in p_name varchar,
    in p_email varchar,
    in p_phone varchar default null,
    in p_rol varchar default 'ROLE_STUDENT',
    in p_whatsapp varchar default null,
    in p_details jsonb default '{}'::jsonb
)
language plpgsql
as $$
begin
    success := true;
    code := 'OK';
    o_usr_id := null;

    -- Check if user already exists
    select usr_id into o_usr_id
    from public.cht_users
    where usr_identity_number = p_identity_number;

    if o_usr_id is not null then
        success := false;
        code := 'ERR_USER_ALREADY_EXISTS';
        return;
    end if;

    -- Create new user
    insert into public.cht_users (
        usr_identity_number,
        usr_name,
        usr_email,
        usr_phone,
        usr_rol,
        usr_whatsapp,
        usr_details
    ) values (
        p_identity_number,
        p_name,
        p_email,
        p_phone,
        p_rol,
        p_whatsapp,
        p_details
    )
    returning usr_id into o_usr_id;

exception
    when others then
        success := false;
        code := 'ERR_CREATE_USER';
        o_usr_id := null;
        raise notice 'Error creating user: %', sqlerrm;
end;
$$;

-- =====================================================
-- Procedure: sp_update_user_whatsapp
-- Description: Updates user's WhatsApp number
-- Returns: success (boolean), code (varchar)
-- =====================================================
create or replace procedure sp_update_user_whatsapp(
    out success boolean,
    out code varchar,
    in p_identity_number varchar,
    in p_whatsapp varchar
)
language plpgsql
as $$
begin
    success := true;
    code := 'OK';

    update public.cht_users
    set usr_whatsapp = p_whatsapp
    where usr_identity_number = p_identity_number;

    if not found then
        success := false;
        code := 'ERR_USER_NOT_FOUND';
    end if;

exception
    when others then
        success := false;
        code := 'ERR_UPDATE_USER_WHATSAPP';
        raise notice 'Error updating user WhatsApp: %', sqlerrm;
end;
$$;

-- =====================================================
-- Comments
-- =====================================================

-- WhatsApp Sessions
comment on function fn_get_whatsapp_session(varchar) is 'Get WhatsApp session by session name';
comment on function fn_get_active_whatsapp_session() is 'Get currently active WhatsApp session';
comment on procedure sp_update_whatsapp_session_status is 'Update WhatsApp session connection status. Returns success and code';
comment on procedure sp_update_whatsapp_qr_code is 'Update QR code for WhatsApp session. Returns success and code';

-- Conversations
comment on function fn_get_conversation_by_chat_id(varchar) is 'Get conversation by WhatsApp chat ID';
comment on procedure sp_create_conversation is 'Create or get existing conversation. Returns success, code, and cnv_id';
comment on procedure sp_link_user_to_conversation is 'Link validated user to conversation. Returns success and code';

-- Messages
comment on function fn_get_conversation_history(varchar, int) is 'Get message history for a conversation';
comment on procedure sp_create_conversation_message is 'Create new message in conversation. Returns success, code, and cvm_id';

-- Users
comment on function fn_get_user_by_identity(varchar) is 'Get user by identity number';
comment on function fn_get_user_by_whatsapp(varchar) is 'Get user by WhatsApp number';
comment on procedure sp_create_user is 'Create new user from WhatsApp registration. Returns success, code, and usr_id';
comment on procedure sp_update_user_whatsapp is 'Update user WhatsApp number. Returns success and code';
