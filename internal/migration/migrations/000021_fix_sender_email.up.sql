-- =====================================================
-- Fix Sender Email to AWS SES Verified Address
-- =====================================================

-- Update EMAIL_CONFIG to use AWS SES verified email
UPDATE cht_parameters
SET prm_data = jsonb_set(
    prm_data,
    '{senderEmail}',
    '"automatizaciones@tikee.tech"'
)
WHERE prm_code = 'EMAIL_CONFIG';

-- Verify the update
DO $$
DECLARE
    v_email TEXT;
BEGIN
    SELECT prm_data->>'senderEmail' INTO v_email FROM cht_parameters WHERE prm_code = 'EMAIL_CONFIG';
    RAISE NOTICE 'Sender email updated to AWS SES verified address: %', v_email;
END $$;
