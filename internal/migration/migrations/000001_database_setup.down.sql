-- Rollback database setup
-- WARNING: This will drop all extensions and schemas

DROP EXTENSION IF EXISTS vector CASCADE;
DROP EXTENSION IF EXISTS pgcrypto CASCADE;
DROP EXTENSION IF NOT EXISTS "uuid-ossp" CASCADE;

-- Note: We don't drop the public schema as it's a default PostgreSQL schema
DROP SCHEMA IF EXISTS ex CASCADE;
