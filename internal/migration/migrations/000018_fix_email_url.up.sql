-- =====================================================
-- Fix Email URL - Force update to correct Tikee endpoint
-- =====================================================

-- Force update EMAIL_CONFIG to correct URL
UPDATE cht_parameters
SET prm_data = '{"senderEmail": "automatizaciones@tikee.tech", "tikeeURL": "http://20.84.48.225:5056/api/emails/enviarDirecto"}'::jsonb,
    prm_description = 'Email sender and Tikee API configuration for OTP emails'
WHERE prm_code = 'EMAIL_CONFIG';

-- Verify the update
DO $$
DECLARE
    v_url TEXT;
BEGIN
    SELECT prm_data->>'tikeeURL' INTO v_url FROM cht_parameters WHERE prm_code = 'EMAIL_CONFIG';
    RAISE NOTICE 'EMAIL_CONFIG updated. Current URL: %', v_url;
END $$;
