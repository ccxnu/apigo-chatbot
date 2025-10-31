-- Debug OTP expiration times
SELECT
    pnd_whatsapp,
    pnd_email,
    pnd_otp_code,
    pnd_otp_generated_at,
    pnd_otp_expires_at,
    CURRENT_TIMESTAMP as current_time,
    (pnd_otp_expires_at - CURRENT_TIMESTAMP) as time_until_expiration,
    pnd_otp_attempts,
    pnd_verified,
    pnd_created_at
FROM cht_pending_registrations
WHERE pnd_whatsapp = '593959423327@s.whatsapp.net'
ORDER BY pnd_created_at DESC
LIMIT 1;
