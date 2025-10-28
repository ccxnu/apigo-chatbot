-- Core Database Tables
-- This migration creates all core tables for the chatbot system

SET search_path TO public, ex;

-- =====================================================
-- Table: cht_parameters (Configuration)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_parameters (
    prm_id              SERIAL PRIMARY KEY,
    prm_name            VARCHAR(100) NOT NULL,
    prm_code            VARCHAR(100) UNIQUE NOT NULL,
    prm_data            JSONB NOT NULL DEFAULT '{}'::JSONB,
    prm_description     VARCHAR(500),
    prm_active          BOOLEAN NOT NULL DEFAULT TRUE,
    prm_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    prm_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_permissions (Security)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_permissions (
    prm_id              SERIAL PRIMARY KEY,
    prm_rol             VARCHAR(50) NOT NULL REFERENCES cht_parameters(prm_code),
    prm_funcionality    VARCHAR(50) NOT NULL REFERENCES cht_parameters(prm_code),
    prm_active          BOOLEAN NOT NULL DEFAULT TRUE,
    prm_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_permission_rol_func UNIQUE (prm_rol, prm_funcionality)
);

-- =====================================================
-- Table: cht_users (Users)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_users (
    usr_id              SERIAL PRIMARY KEY,
    usr_identity_number VARCHAR(20) UNIQUE NOT NULL,
    usr_name            VARCHAR(100) NOT NULL,
    usr_email           VARCHAR(100) UNIQUE NOT NULL,
    usr_phone           VARCHAR(20),
    usr_rol             VARCHAR(50) REFERENCES cht_parameters(prm_code),
    usr_details         JSONB DEFAULT '{}'::JSONB,
    usr_whatsapp        VARCHAR(50),
    usr_active          BOOLEAN NOT NULL DEFAULT TRUE,
    usr_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    usr_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_sessions (Sessions)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_sessions (
    ssn_id              SERIAL PRIMARY KEY,
    ssn_fk_user         INT NOT NULL REFERENCES cht_users(usr_id),
    ssn_origin          VARCHAR(50) NOT NULL REFERENCES cht_parameters(prm_code),
    ssn_token           VARCHAR(500) UNIQUE NOT NULL,
    ssn_started_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ssn_ended_at        TIMESTAMP,
    ssn_active          BOOLEAN NOT NULL DEFAULT TRUE
);

-- =====================================================
-- Table: cht_documents (Knowledge - Document)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_documents (
    doc_id              SERIAL PRIMARY KEY,
    doc_category        VARCHAR(50) NOT NULL REFERENCES cht_parameters(prm_code),
    doc_title           VARCHAR(200) NOT NULL,
    doc_summary         TEXT,
    doc_source          VARCHAR(500),
    doc_published_at    TIMESTAMP,
    doc_active          BOOLEAN NOT NULL DEFAULT TRUE,
    doc_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    doc_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_chunks (Knowledge - Chunk)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_chunks (
    chk_id              SERIAL PRIMARY KEY,
    chk_fk_document     INT NOT NULL REFERENCES cht_documents(doc_id) ON DELETE CASCADE,
    chk_content         TEXT NOT NULL,
    chk_embedding       VECTOR(1536),
    chk_fts_vector      TSVECTOR GENERATED ALWAYS AS (to_tsvector('spanish'::regconfig, chk_content)) STORED,
    chk_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    chk_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_chunk_statistics (Knowledge - Statistics)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_chunk_statistics (
    cst_id                      SERIAL PRIMARY KEY,
    cst_fk_chunk                INT UNIQUE NOT NULL REFERENCES cht_chunks(chk_id) ON DELETE CASCADE,
    cst_usage_count             INT NOT NULL DEFAULT 0,
    cst_last_used_at            TIMESTAMP,
    cst_precision_atk           FLOAT,
    cst_recall_atk              FLOAT,
    cst_f1_atk                  FLOAT,
    cst_mrr                     FLOAT,
    cst_map                     FLOAT,
    cst_ndcg                    FLOAT,
    cst_staleness_days          INT,
    cst_last_refresh_at         TIMESTAMP,
    cst_curriculum_coverage_pct FLOAT,
    cst_created_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cst_updated_at              TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_logs (System Logs)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_logs (
    log_id              SERIAL PRIMARY KEY,
    log_level           VARCHAR(20) NOT NULL CHECK (log_level IN ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL')),
    log_module          VARCHAR(100) NOT NULL,
    log_message         TEXT NOT NULL,
    log_details         JSONB,
    log_user_id         INT REFERENCES cht_users(usr_id),
    log_session_id      INT REFERENCES cht_sessions(ssn_id),
    log_ip_address      INET,
    log_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_whatsapp_sessions (WhatsApp Connection)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_whatsapp_sessions (
    wss_id              SERIAL PRIMARY KEY,
    wss_session_name    VARCHAR(50) UNIQUE NOT NULL,
    wss_phone_number    VARCHAR(50),
    wss_device_name     VARCHAR(100),
    wss_platform        VARCHAR(50),
    wss_qr_code         TEXT,
    wss_connected       BOOLEAN NOT NULL DEFAULT FALSE,
    wss_last_seen       TIMESTAMP,
    wss_session_data    JSONB DEFAULT '{}'::JSONB,
    wss_active          BOOLEAN NOT NULL DEFAULT TRUE,
    wss_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    wss_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Table: cht_conversations (WhatsApp Conversations)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_conversations (
    cnv_id              SERIAL PRIMARY KEY,
    cnv_fk_user         INT REFERENCES cht_users(usr_id),
    cnv_chat_id         VARCHAR(100) NOT NULL,
    cnv_phone_number    VARCHAR(50) NOT NULL,
    cnv_contact_name    VARCHAR(100),
    cnv_is_group        BOOLEAN NOT NULL DEFAULT FALSE,
    cnv_group_name      VARCHAR(100),
    cnv_last_message_at TIMESTAMP,
    cnv_message_count   INT NOT NULL DEFAULT 0,
    cnv_active          BOOLEAN NOT NULL DEFAULT TRUE,
    cnv_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    cnv_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_chat_id UNIQUE (cnv_chat_id)
);

-- =====================================================
-- Table: cht_conversation_messages (WhatsApp Messages)
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_conversation_messages (
    cvm_id              SERIAL PRIMARY KEY,
    cvm_fk_conversation INT NOT NULL REFERENCES cht_conversations(cnv_id) ON DELETE CASCADE,
    cvm_message_id      VARCHAR(100) UNIQUE NOT NULL,
    cvm_from_me         BOOLEAN NOT NULL DEFAULT FALSE,
    cvm_sender_name     VARCHAR(100),
    cvm_message_type    VARCHAR(20) NOT NULL CHECK (cvm_message_type IN ('text', 'image', 'document', 'audio', 'video', 'sticker')),
    cvm_body            TEXT,
    cvm_media_url       VARCHAR(500),
    cvm_quoted_message  VARCHAR(100),
    cvm_timestamp       BIGINT NOT NULL,
    cvm_is_forwarded    BOOLEAN NOT NULL DEFAULT FALSE,
    cvm_metadata        JSONB DEFAULT '{}'::JSONB,
    cvm_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Parameters
CREATE INDEX IF NOT EXISTS idx_cht_parameters_code_active ON cht_parameters(prm_code, prm_active);
CREATE INDEX IF NOT EXISTS idx_cht_parameters_data ON cht_parameters USING GIN(prm_data);

-- Permissions
CREATE INDEX IF NOT EXISTS idx_cht_permissions_rol ON cht_permissions(prm_rol);
CREATE INDEX IF NOT EXISTS idx_cht_permissions_active ON cht_permissions(prm_active);

-- Users
CREATE INDEX IF NOT EXISTS idx_cht_users_email ON cht_users(usr_email);
CREATE INDEX IF NOT EXISTS idx_cht_users_rol ON cht_users(usr_rol);
CREATE INDEX IF NOT EXISTS idx_cht_users_active ON cht_users(usr_active);
CREATE INDEX IF NOT EXISTS idx_cht_users_details ON cht_users USING GIN(usr_details);

-- Sessions
CREATE INDEX IF NOT EXISTS idx_cht_sessions_user ON cht_sessions(ssn_fk_user);
CREATE INDEX IF NOT EXISTS idx_cht_sessions_token ON cht_sessions(ssn_token);
CREATE INDEX IF NOT EXISTS idx_cht_sessions_active ON cht_sessions(ssn_active);
CREATE INDEX IF NOT EXISTS idx_cht_sessions_started ON cht_sessions(ssn_started_at DESC);

-- Documents
CREATE INDEX IF NOT EXISTS idx_cht_documents_category ON cht_documents(doc_category);
CREATE INDEX IF NOT EXISTS idx_cht_documents_active ON cht_documents(doc_active);

-- Chunks (CRITICAL for vector search)
CREATE INDEX IF NOT EXISTS idx_cht_chunks_document ON cht_chunks(chk_fk_document);
CREATE INDEX IF NOT EXISTS idx_cht_chunks_embedding ON cht_chunks USING IVFFLAT(chk_embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS chk_fts_idx ON public.cht_chunks USING GIN (chk_fts_vector);

-- Chunk Statistics
CREATE INDEX IF NOT EXISTS idx_cht_chunk_stats_usage ON cht_chunk_statistics(cst_usage_count DESC);
CREATE INDEX IF NOT EXISTS idx_cht_chunk_stats_last_used ON cht_chunk_statistics(cst_last_used_at DESC);

-- Logs
CREATE INDEX IF NOT EXISTS idx_cht_logs_level ON cht_logs(log_level);
CREATE INDEX IF NOT EXISTS idx_cht_logs_module ON cht_logs(log_module);
CREATE INDEX IF NOT EXISTS idx_cht_logs_created ON cht_logs(log_created_at DESC);
CREATE INDEX IF NOT EXISTS idx_cht_logs_user ON cht_logs(log_user_id);
CREATE INDEX IF NOT EXISTS idx_cht_logs_details ON cht_logs USING GIN(log_details);

-- WhatsApp Sessions
CREATE INDEX IF NOT EXISTS idx_cht_whatsapp_sessions_name ON cht_whatsapp_sessions(wss_session_name);
CREATE INDEX IF NOT EXISTS idx_cht_whatsapp_sessions_connected ON cht_whatsapp_sessions(wss_connected);
CREATE INDEX IF NOT EXISTS idx_cht_whatsapp_sessions_active ON cht_whatsapp_sessions(wss_active);

-- Conversations
CREATE INDEX IF NOT EXISTS idx_cht_conversations_user ON cht_conversations(cnv_fk_user);
CREATE INDEX IF NOT EXISTS idx_cht_conversations_chat_id ON cht_conversations(cnv_chat_id);
CREATE INDEX IF NOT EXISTS idx_cht_conversations_phone ON cht_conversations(cnv_phone_number);
CREATE INDEX IF NOT EXISTS idx_cht_conversations_last_msg ON cht_conversations(cnv_last_message_at DESC);
CREATE INDEX IF NOT EXISTS idx_cht_conversations_active ON cht_conversations(cnv_active);

-- Conversation Messages
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_conversation ON cht_conversation_messages(cvm_fk_conversation);
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_message_id ON cht_conversation_messages(cvm_message_id);
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_timestamp ON cht_conversation_messages(cvm_timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_from_me ON cht_conversation_messages(cvm_from_me);
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_type ON cht_conversation_messages(cvm_message_type);

-- =====================================================
-- Update Timestamp Trigger Function
-- =====================================================
CREATE OR REPLACE FUNCTION fn_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'cht_parameters' THEN
        NEW.prm_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_users' THEN
        NEW.usr_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_documents' THEN
        NEW.doc_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_chunks' THEN
        NEW.chk_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_chunk_statistics' THEN
        NEW.cst_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_whatsapp_sessions' THEN
        NEW.wss_updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_TABLE_NAME = 'cht_conversations' THEN
        NEW.cnv_updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers
DROP TRIGGER IF EXISTS tr_cht_parameters_updated ON cht_parameters;
DROP TRIGGER IF EXISTS tr_cht_users_updated ON cht_users;
DROP TRIGGER IF EXISTS tr_cht_documents_updated ON cht_documents;
DROP TRIGGER IF EXISTS tr_cht_chunks_updated ON cht_chunks;
DROP TRIGGER IF EXISTS tr_cht_chunk_stats_updated ON cht_chunk_statistics;
DROP TRIGGER IF EXISTS tr_cht_whatsapp_sessions_updated ON cht_whatsapp_sessions;
DROP TRIGGER IF EXISTS tr_cht_conversations_updated ON cht_conversations;

CREATE TRIGGER tr_cht_parameters_updated BEFORE UPDATE ON cht_parameters FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_users_updated BEFORE UPDATE ON cht_users FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_documents_updated BEFORE UPDATE ON cht_documents FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_chunks_updated BEFORE UPDATE ON cht_chunks FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_chunk_stats_updated BEFORE UPDATE ON cht_chunk_statistics FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_whatsapp_sessions_updated BEFORE UPDATE ON cht_whatsapp_sessions FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();
CREATE TRIGGER tr_cht_conversations_updated BEFORE UPDATE ON cht_conversations FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON TABLE cht_parameters IS 'System configuration parameters';
COMMENT ON TABLE cht_permissions IS 'Role-based permissions matrix';
COMMENT ON TABLE cht_users IS 'System users with role-based access';
COMMENT ON TABLE cht_sessions IS 'Active and historical user sessions';
COMMENT ON TABLE cht_documents IS 'Knowledge base documents';
COMMENT ON TABLE cht_chunks IS 'Document chunks with embeddings for RAG';
COMMENT ON TABLE cht_chunk_statistics IS 'Usage and quality metrics for chunks';
COMMENT ON TABLE cht_logs IS 'System audit and error logs';
COMMENT ON TABLE cht_whatsapp_sessions IS 'WhatsApp connection sessions and device information';
COMMENT ON TABLE cht_conversations IS 'WhatsApp conversations and chat metadata';
COMMENT ON TABLE cht_conversation_messages IS 'Individual WhatsApp messages in conversations';
