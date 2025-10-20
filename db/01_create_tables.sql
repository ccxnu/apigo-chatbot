-- Database Schema - Chatbot System

-- Search path
show search_path;
alter database chatbot_db set search_path to public, ex;
alter role postgres set search_path to public, ex;
reset search_path;
set search_path to public, ex;

-- =====================================================
-- Table: cht_parameters (Configuracion)
-- =====================================================
create table if not exists public.cht_parameters (
    prm_id              serial primary key,
    prm_name            varchar(100) not null,
    prm_code            varchar(100) unique not null,
    prm_data            jsonb not null default '{}'::jsonb,
    prm_description     varchar(500),
    prm_active          boolean not null default true,
    prm_created_at      timestamp not null default current_timestamp,
    prm_updated_at      timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_permissions (Seguridad)
-- =====================================================
create table if not exists public.cht_permissions (
    prm_id              serial primary key,
    prm_rol             varchar(50) not null references cht_parameters(prm_code),
    prm_funcionality    varchar(50) not null references cht_parameters(prm_code),
    prm_active          boolean not null default true,
    prm_created_at      timestamp not null default current_timestamp,
    constraint uk_permission_rol_func unique (prm_rol, prm_funcionality)
);

-- =====================================================
-- Table: cht_users (Usuarios)
-- =====================================================
create table if not exists public.cht_users (
    usr_id              serial primary key,
    usr_identity_number varchar(20) unique not null,
    usr_name            varchar(100) not null,
    usr_email           varchar(100) unique not null,
    usr_phone           varchar(20),
    usr_active          boolean not null default true,
    usr_created_at      timestamp not null default current_timestamp,
    usr_updated_at      timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_sessions (Sesiones)
-- =====================================================
create table if not exists public.cht_sessions (
    ssn_id              serial primary key,
    ssn_fk_user         int not null references cht_users(usr_id),
    ssn_origin          varchar(50) not null references cht_parameters(prm_code),
    ssn_token           varchar(500) unique not null,
    ssn_started_at      timestamp not null default current_timestamp,
    ssn_ended_at        timestamp,
    ssn_active          boolean not null default true
);

-- =====================================================
-- Table: cht_documents (Conocimiento - Documento)
-- =====================================================
create table if not exists public.cht_documents (
    doc_id              serial primary key,
    doc_category        varchar(50) not null references cht_parameters(prm_code),
    doc_title           varchar(200) not null,
    doc_summary         text,
    doc_source          varchar(500),
    doc_published_at    timestamp,
    doc_active          boolean not null default true,
    doc_created_at      timestamp not null default current_timestamp,
    doc_updated_at      timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_chunks (Conocimiento - Chunk)
-- =====================================================
create table if not exists public.cht_chunks (
    chk_id              serial primary key,
    chk_fk_document     int not null references cht_documents(doc_id) on delete cascade,
    chk_content         text not null,
    chk_embedding       vector(1536),  -- OpenAI/Claude embedding dimension
    chk_created_at      timestamp not null default current_timestamp,
    chk_updated_at      timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_chunk_statistics (Conocimiento - EstadisticaChunk)
-- =====================================================
create table if not exists public.cht_chunk_statistics (
    cst_id                      serial primary key,
    cst_fk_chunk                int unique not null references cht_chunks(chk_id) on delete cascade,
    cst_usage_count             int not null default 0,
    cst_last_used_at            timestamp,
    cst_precision_atk           float,
    cst_recall_atk              float,
    cst_f1_atk                  float,
    cst_mrr                     float,
    cst_map                     float,
    cst_ndcg                    float,
    cst_staleness_days          int,
    cst_last_refresh_at         timestamp,
    cst_curriculum_coverage_pct float,
    cst_created_at              timestamp not null default current_timestamp,
    cst_updated_at              timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_messages (Sesiones - MensajeChat)
-- =====================================================
create table if not exists public.cht_messages (
    msg_id              serial primary key,
    msg_fk_session      int not null references cht_sessions(ssn_id) on delete cascade,
    msg_rol             varchar(50) not null references cht_parameters(prm_code),
    msg_content         text not null,
    msg_embedding       vector(1536),
    msg_processed_at    timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_message_statistics (Sesiones - EstadisticaMensaje)
-- =====================================================
create table if not exists public.cht_message_statistics (
    mst_id                      serial primary key,
    mst_fk_message              int unique not null references cht_messages(msg_id) on delete cascade,
    mst_latency_ms              int,
    mst_tokens_prompt           int,
    mst_tokens_completion       int,
    mst_cost_usd                decimal(10,6),
    mst_faithfulness_score      float,
    mst_answer_correct_rate     float,
    mst_feedback_rating         int check (mst_feedback_rating between 1 and 5),
    mst_created_at              timestamp not null default current_timestamp
);

-- =====================================================
-- Table: cht_logs (Sistema de Logs)
-- =====================================================
create table if not exists public.cht_logs (
    log_id              serial primary key,
    log_level           varchar(20) not null check (log_level in ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL')),
    log_module          varchar(100) not null,
    log_message         text not null,
    log_details         jsonb,
    log_user_id         int references cht_users(usr_id),
    log_session_id      int references cht_sessions(ssn_id),
    log_ip_address      inet,
    log_created_at      timestamp not null default current_timestamp
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Parameters
create index if not exists idx_cht_parameters_code_active on cht_parameters(prm_code, prm_active);
create index if not exists idx_cht_parameters_data on cht_parameters using gin(prm_data);

-- Permissions
create index if not exists idx_cht_permissions_rol on cht_permissions(prm_rol);
create index if not exists idx_cht_permissions_active on cht_permissions(prm_active);

-- Users
create index if not exists idx_cht_users_email on cht_users(usr_email);
create index if not exists idx_cht_users_rol on cht_users(usr_rol);
create index if not exists idx_cht_users_active on cht_users(usr_active);
create index if not exists idx_cht_users_details on cht_users using gin(usr_details);

-- Sessions
create index if not exists idx_cht_sessions_user on cht_sessions(ssn_fk_user);
create index if not exists idx_cht_sessions_token on cht_sessions(ssn_token);
create index if not exists idx_cht_sessions_active on cht_sessions(ssn_active);
create index if not exists idx_cht_sessions_started on cht_sessions(ssn_started_at desc);

-- Documents
create index if not exists idx_cht_documents_category on cht_documents(doc_category);
create index if not exists idx_cht_documents_active on cht_documents(doc_active);

-- Chunks (CRITICAL for vector search)
create index if not exists idx_cht_chunks_document on cht_chunks(chk_fk_document);
create index if not exists idx_cht_chunks_embedding on cht_chunks using ivfflat(chk_embedding vector_cosine_ops) with (lists = 100);

-- Chunk Statistics
create index if not exists idx_cht_chunk_stats_usage on cht_chunk_statistics(cst_usage_count desc);
create index if not exists idx_cht_chunk_stats_last_used on cht_chunk_statistics(cst_last_used_at desc);

-- Messages
create index if not exists idx_cht_messages_session on cht_messages(msg_fk_session);
create index if not exists idx_cht_messages_processed on cht_messages(msg_processed_at desc);
create index if not exists idx_cht_messages_embedding on cht_messages using ivfflat(msg_embedding vector_cosine_ops) with (lists = 100);

-- Message Statistics
create index if not exists idx_cht_msg_stats_message on cht_message_statistics(mst_fk_message);

-- Logs
create index if not exists idx_cht_logs_level on cht_logs(log_level);
create index if not exists idx_cht_logs_module on cht_logs(log_module);
create index if not exists idx_cht_logs_created on cht_logs(log_created_at desc);
create index if not exists idx_cht_logs_user on cht_logs(log_user_id);
create index if not exists idx_cht_logs_details on cht_logs using gin(log_details);

-- =====================================================
-- Update Timestamp Trigger Function
-- =====================================================
create or replace function fn_update_timestamp()
returns trigger as $$
begin
    -- Dinámicamente actualizar el campo *_updated_at según la tabla
    if TG_TABLE_NAME = 'cht_parameters' then
        new.prm_updated_at = current_timestamp;
    elsif TG_TABLE_NAME = 'cht_users' then
        new.usr_updated_at = current_timestamp;
    elsif TG_TABLE_NAME = 'cht_documents' then
        new.doc_updated_at = current_timestamp;
    elsif TG_TABLE_NAME = 'cht_chunks' then
        new.chk_updated_at = current_timestamp;
    elsif TG_TABLE_NAME = 'cht_chunk_statistics' then
        new.cst_updated_at = current_timestamp;
    end if;
    return new;
end;
$$ language plpgsql;

-- Apply triggers
drop trigger if exists tr_cht_parameters_updated on cht_parameters;
drop trigger if exists tr_cht_users_updated on cht_users;
drop trigger if exists tr_cht_documents_updated on cht_documents;
drop trigger if exists tr_cht_chunks_updated on cht_chunks;
drop trigger if exists tr_cht_chunk_stats_updated on cht_chunk_statistics;

create trigger tr_cht_parameters_updated before update on cht_parameters for each row execute function fn_update_timestamp();
create trigger tr_cht_users_updated before update on cht_users for each row execute function fn_update_timestamp();
create trigger tr_cht_documents_updated before update on cht_documents for each row execute function fn_update_timestamp();
create trigger tr_cht_chunks_updated before update on cht_chunks for each row execute function fn_update_timestamp();
create trigger tr_cht_chunk_stats_updated before update on cht_chunk_statistics for each row execute function fn_update_timestamp();

-- =====================================================
-- Comments
-- =====================================================
comment on table cht_parameters is 'System configuration parameters';
comment on table cht_permissions is 'Role-based permissions matrix';
comment on table cht_users is 'System users with role-based access';
comment on table cht_sessions is 'Active and historical user sessions';
comment on table cht_documents is 'Knowledge base documents';
comment on table cht_chunks is 'Document chunks with embeddings for RAG';
comment on table cht_chunk_statistics is 'Usage and quality metrics for chunks';
comment on table cht_messages is 'Chat messages with AI responses';
comment on table cht_message_statistics is 'Performance metrics for messages';
comment on table cht_logs is 'System audit and error logs';
