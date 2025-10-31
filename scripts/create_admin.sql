-- =====================================================
-- Create Admin User Script
-- =====================================================
-- Instructions:
-- 1. Edit the values below (username, email, password, name)
-- 2. Run: psql -h localhost -p 5432 -U postgres -d chatbot_db -f scripts/create_admin.sql
--
-- The password will be hashed using bcrypt
-- =====================================================

DO $$
DECLARE
    v_username VARCHAR(50) := 'admin';                    -- CHANGE THIS
    v_email VARCHAR(100) := 'admin@ists.edu.ec';          -- CHANGE THIS
    v_password VARCHAR(100) := 'Admin123!';               -- CHANGE THIS
    v_name VARCHAR(100) := 'System Administrator';        -- CHANGE THIS
    v_role VARCHAR(50) := 'ROLE_ADMIN';
    v_password_hash TEXT;
    v_success BOOLEAN;
    v_code VARCHAR;
    v_admin_id INT;
BEGIN
    -- Check if pgcrypto extension exists (for password hashing)
    CREATE EXTENSION IF NOT EXISTS pgcrypto;

    -- Hash password using crypt (bcrypt algorithm)
    -- Note: In production, you should use your application's bcrypt hasher
    -- This is a workaround for SQL-only admin creation
    v_password_hash := crypt(v_password, gen_salt('bf', 10));

    -- Check if username already exists
    IF EXISTS (SELECT 1 FROM cht_admin_users WHERE adm_username = v_username) THEN
        RAISE EXCEPTION 'Username "%" already exists. Please choose a different username.', v_username;
    END IF;

    -- Check if email already exists
    IF EXISTS (SELECT 1 FROM cht_admin_users WHERE adm_email = v_email) THEN
        RAISE EXCEPTION 'Email "%" already exists. Please use a different email.', v_email;
    END IF;

    -- Check if role exists in parameters table
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = v_role) THEN
        RAISE EXCEPTION 'Role "%" does not exist in cht_parameters table.', v_role;
    END IF;

    -- Insert admin user
    INSERT INTO cht_admin_users (
        adm_username,
        adm_email,
        adm_password_hash,
        adm_name,
        adm_role,
        adm_permissions,
        adm_claims,
        adm_is_active
    ) VALUES (
        v_username,
        v_email,
        v_password_hash,
        v_name,
        v_role,
        '[]'::jsonb,  -- No custom permissions
        '{}'::jsonb,  -- No custom claims
        true
    )
    RETURNING adm_id INTO v_admin_id;

    -- Success message
    RAISE NOTICE '';
    RAISE NOTICE '=====================================================';
    RAISE NOTICE 'Admin user created successfully!';
    RAISE NOTICE '=====================================================';
    RAISE NOTICE 'Admin ID:  %', v_admin_id;
    RAISE NOTICE 'Username:  %', v_username;
    RAISE NOTICE 'Email:     %', v_email;
    RAISE NOTICE 'Name:      %', v_name;
    RAISE NOTICE 'Role:      %', v_role;
    RAISE NOTICE '';
    RAISE NOTICE 'You can now login at: http://localhost:8080/admin/login';
    RAISE NOTICE '=====================================================';

EXCEPTION
    WHEN OTHERS THEN
        RAISE NOTICE '';
        RAISE NOTICE '✗ Error creating admin user: %', SQLERRM;
        RAISE NOTICE '';
        RAISE;
END $$;
