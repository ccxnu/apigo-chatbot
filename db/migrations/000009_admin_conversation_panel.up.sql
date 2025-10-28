-- =====================================================
-- Admin Conversation Panel - Schema Updates
-- =====================================================

-- Add blocked status to users
ALTER TABLE public.cht_users
ADD COLUMN IF NOT EXISTS usr_blocked boolean NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS usr_blocked_at timestamp,
ADD COLUMN IF NOT EXISTS usr_blocked_by int REFERENCES cht_admin_users(adm_id),
ADD COLUMN IF NOT EXISTS usr_block_reason text;

-- Add conversation management fields
ALTER TABLE public.cht_conversations
ADD COLUMN IF NOT EXISTS cnv_blocked boolean NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS cnv_admin_intervened boolean NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS cnv_last_admin_message_at timestamp,
ADD COLUMN IF NOT EXISTS cnv_temporary boolean NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS cnv_expires_at timestamp,
ADD COLUMN IF NOT EXISTS cnv_unread_count int NOT NULL DEFAULT 0;

-- Add message sender type (user vs admin)
ALTER TABLE public.cht_conversation_messages
ADD COLUMN IF NOT EXISTS cvm_sender_type varchar(20) NOT NULL DEFAULT 'user' CHECK (cvm_sender_type IN ('user', 'admin', 'bot')),
ADD COLUMN IF NOT EXISTS cvm_admin_id int REFERENCES cht_admin_users(adm_id),
ADD COLUMN IF NOT EXISTS cvm_read boolean NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS cvm_read_at timestamp;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_blocked ON cht_users(usr_blocked);
CREATE INDEX IF NOT EXISTS idx_users_whatsapp ON cht_users(usr_whatsapp);

CREATE INDEX IF NOT EXISTS idx_conversations_blocked ON cht_conversations(cnv_blocked);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message ON cht_conversations(cnv_last_message_at DESC);
CREATE INDEX IF NOT EXISTS idx_conversations_unread ON cht_conversations(cnv_unread_count) WHERE cnv_unread_count > 0;
CREATE INDEX IF NOT EXISTS idx_conversations_temporary ON cht_conversations(cnv_temporary, cnv_expires_at) WHERE cnv_temporary = true;
CREATE INDEX IF NOT EXISTS idx_conversations_user ON cht_conversations(cnv_fk_user);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON cht_conversation_messages(cvm_fk_conversation, cvm_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_messages_sender_type ON cht_conversation_messages(cvm_sender_type);
CREATE INDEX IF NOT EXISTS idx_messages_read ON cht_conversation_messages(cvm_read) WHERE cvm_read = false;

-- Comments
COMMENT ON COLUMN cht_users.usr_blocked IS 'User is blocked from using the chatbot';
COMMENT ON COLUMN cht_users.usr_blocked_at IS 'When user was blocked';
COMMENT ON COLUMN cht_users.usr_blocked_by IS 'Admin who blocked the user';
COMMENT ON COLUMN cht_users.usr_block_reason IS 'Reason for blocking';

COMMENT ON COLUMN cht_conversations.cnv_blocked IS 'Conversation is blocked (user cannot send messages)';
COMMENT ON COLUMN cht_conversations.cnv_admin_intervened IS 'Admin has sent messages in this conversation';
COMMENT ON COLUMN cht_conversations.cnv_last_admin_message_at IS 'Last time admin sent a message';
COMMENT ON COLUMN cht_conversations.cnv_temporary IS 'Conversation will auto-delete after expiry';
COMMENT ON COLUMN cht_conversations.cnv_expires_at IS 'When temporary conversation expires';
COMMENT ON COLUMN cht_conversations.cnv_unread_count IS 'Number of unread messages from user';

COMMENT ON COLUMN cht_conversation_messages.cvm_sender_type IS 'Who sent the message: user, admin, or bot';
COMMENT ON COLUMN cht_conversation_messages.cvm_admin_id IS 'Admin who sent the message (if sender_type=admin)';
COMMENT ON COLUMN cht_conversation_messages.cvm_read IS 'Message has been read by admin';
COMMENT ON COLUMN cht_conversation_messages.cvm_read_at IS 'When message was read';
