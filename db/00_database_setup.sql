-- Database Setup Script - Chatbot System

-- Create Database (run as postgres superuser)
create database chatbot_db
    with
    owner = postgres
    encoding = 'UTF8'
    tablespace = pg_default
    connection limit = -1;

-- =====================================================
-- Connect to the database
-- =====================================================
-- \c chatbot_db;

-- =====================================================
-- Create Schemas
-- =====================================================
create schema if not exists public;
create schema if not exists ex;

-- =====================================================
-- Install Extensions in 'ex' schema
-- =====================================================
-- Vector extension for embeddings
create extension if not exists vector with schema ex;

-- PGCrypto for encryption functions
create extension if not exists pgcrypto with schema ex;

-- UUID generation
create extension if not exists "uuid-ossp" with schema ex;

-- =====================================================
-- Configure Search Path (KEY para usar sin prefijo)
-- =====================================================

-- For the database (affects all new connections)
alter database chatbot_db set search_path to public, ex;

-- For current session
set search_path to public, ex;

-- =====================================================
-- Grant Usage on Schemas
-- =====================================================
grant usage on schema public to public;
grant usage on schema ex to public;
grant execute on all functions in schema ex to public;
alter default privileges in schema ex grant execute on functions to public;

-- =====================================================
-- Comments
-- =====================================================
comment on schema public is 'Main application schema for all application objects';
comment on schema ex is 'Schema for PostgreSQL extensions to isolate their functions';

-- =====================================================
-- Verification
-- =====================================================
-- Check installed extensions
select
    e.extname as extension_name,
    n.nspname as schema_name,
    e.extversion as version
from pg_extension e
join pg_namespace n on e.extnamespace = n.oid
order by e.extname;

-- Check search_path
select current_setting('search_path') as current_search_path;
-- Expected: public, ex
