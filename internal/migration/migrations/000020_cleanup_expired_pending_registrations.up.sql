-- =====================================================
-- Cleanup Expired Pending Registrations
-- =====================================================

-- Delete pending registrations where OTP has expired (older than 10 minutes)
DELETE FROM cht_pending_registrations
WHERE pnd_otp_expires_at < CURRENT_TIMESTAMP;

-- Log the cleanup
DO $$
DECLARE
    v_count INT;
BEGIN
    GET DIAGNOSTICS v_count = ROW_COUNT;
    RAISE NOTICE 'Cleaned up % expired pending registrations', v_count;
END $$;
