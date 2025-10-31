-- =====================================================
-- Rollback Email URL Fix
-- =====================================================

-- Revert to old URL (if needed for rollback)
UPDATE cht_parameters
SET prm_data = '{"sender": "noreply@example.com"}'::jsonb,
    prm_description = 'Email sender configuration'
WHERE prm_code = 'EMAIL_CONFIG';
