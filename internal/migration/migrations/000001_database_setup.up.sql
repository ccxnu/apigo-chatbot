-- Database Setup - Extensions and Schemas
-- This is the foundational migration

-- =====================================================
-- Create Schemas
-- =====================================================
CREATE SCHEMA IF NOT EXISTS public;
CREATE SCHEMA IF NOT EXISTS ex;

-- =====================================================
-- Install Extensions in 'ex' schema
-- =====================================================
-- Vector extension for embeddings
CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA ex;

-- PGCrypto for encryption functions
CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA ex;

-- UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA ex;

-- =====================================================
-- Configure Search Path
-- =====================================================
-- For the database (affects all new connections)
ALTER DATABASE chatbot_db SET search_path TO public, ex;

-- =====================================================
-- Grant Usage on Schemas
-- =====================================================
GRANT USAGE ON SCHEMA public TO public;
GRANT USAGE ON SCHEMA ex TO public;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA ex TO public;
ALTER DEFAULT PRIVILEGES IN SCHEMA ex GRANT EXECUTE ON FUNCTIONS TO public;

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON SCHEMA public IS 'Main application schema for all application objects';
COMMENT ON SCHEMA ex IS 'Schema for PostgreSQL extensions to isolate their functions';
