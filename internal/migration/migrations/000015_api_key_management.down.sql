-- =====================================================
-- API Key Management System
-- Migration: 000015_api_key_management.down.sql
-- Purpose: Rollback API key management procedures and usage tracking
-- =====================================================

-- Drop functions
DROP FUNCTION IF EXISTS fn_get_api_usage_stats(INT, TIMESTAMP, TIMESTAMP);
DROP FUNCTION IF EXISTS fn_get_api_key_by_id(INT);
DROP FUNCTION IF EXISTS fn_get_all_api_keys();
DROP FUNCTION IF EXISTS fn_get_api_key_by_value(TEXT);

-- Drop procedures
DROP PROCEDURE IF EXISTS sp_track_api_usage(BOOLEAN, VARCHAR, INT, INT, VARCHAR, VARCHAR, INT, INT, INT, VARCHAR, TEXT, TEXT);
DROP PROCEDURE IF EXISTS sp_delete_api_key(BOOLEAN, VARCHAR, INT);
DROP PROCEDURE IF EXISTS sp_update_api_key(BOOLEAN, VARCHAR, INT, VARCHAR, INT, JSONB, JSONB, BOOLEAN, TIMESTAMP);
DROP PROCEDURE IF EXISTS sp_update_api_key_last_used(BOOLEAN, VARCHAR, INT);
DROP PROCEDURE IF EXISTS sp_create_api_key(BOOLEAN, VARCHAR, INT, VARCHAR, TEXT, VARCHAR, JSONB, INT, JSONB, JSONB, TIMESTAMP, INT);

-- Drop indexes
DROP INDEX IF EXISTS idx_api_usage_status;
DROP INDEX IF EXISTS idx_api_usage_endpoint;
DROP INDEX IF EXISTS idx_api_usage_created;
DROP INDEX IF EXISTS idx_api_usage_key_id;

-- Drop table
DROP TABLE IF EXISTS cht_api_usage;
