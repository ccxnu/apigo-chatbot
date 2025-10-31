-- =====================================================
-- Rollback Sender Email Fix
-- =====================================================

-- Revert to old sender email
UPDATE cht_parameters
SET prm_data = jsonb_set(
    prm_data,
    '{senderEmail}',
    '"noreply@ists.edu.ec"'
)
WHERE prm_code = 'EMAIL_CONFIG';
