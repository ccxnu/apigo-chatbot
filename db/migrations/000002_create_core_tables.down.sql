-- Rollback core tables
-- WARNING: This will drop all core tables and their data

DROP TRIGGER IF EXISTS tr_cht_conversations_updated ON cht_conversations;
DROP TRIGGER IF EXISTS tr_cht_whatsapp_sessions_updated ON cht_whatsapp_sessions;
DROP TRIGGER IF EXISTS tr_cht_chunk_stats_updated ON cht_chunk_statistics;
DROP TRIGGER IF EXISTS tr_cht_chunks_updated ON cht_chunks;
DROP TRIGGER IF EXISTS tr_cht_documents_updated ON cht_documents;
DROP TRIGGER IF EXISTS tr_cht_users_updated ON cht_users;
DROP TRIGGER IF EXISTS tr_cht_parameters_updated ON cht_parameters;

DROP FUNCTION IF EXISTS fn_update_timestamp();

DROP TABLE IF EXISTS cht_conversation_messages CASCADE;
DROP TABLE IF EXISTS cht_conversations CASCADE;
DROP TABLE IF EXISTS cht_whatsapp_sessions CASCADE;
DROP TABLE IF EXISTS cht_logs CASCADE;
DROP TABLE IF EXISTS cht_chunk_statistics CASCADE;
DROP TABLE IF EXISTS cht_chunks CASCADE;
DROP TABLE IF EXISTS cht_documents CASCADE;
DROP TABLE IF EXISTS cht_sessions CASCADE;
DROP TABLE IF EXISTS cht_users CASCADE;
DROP TABLE IF EXISTS cht_permissions CASCADE;
DROP TABLE IF EXISTS cht_parameters CASCADE;
