-- =====================================================
-- OTP-Based Registration System
-- =====================================================
-- This migration adds tables for secure WhatsApp registration
-- with email OTP verification to prevent identity theft

-- =====================================================
-- Table: cht_pending_registrations
-- Stores users awaiting OTP verification
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_pending_registrations (
    pnd_id              SERIAL PRIMARY KEY,
    pnd_identity_number VARCHAR(20) NOT NULL,
    pnd_whatsapp        VARCHAR(50) NOT NULL,
    pnd_name            VARCHAR(100),
    pnd_email           VARCHAR(100),
    pnd_phone           VARCHAR(20),
    pnd_role            VARCHAR(50),          -- ROLE_STUDENT, ROLE_PROFESSOR, ROLE_EXTERNAL
    pnd_user_type       VARCHAR(20) NOT NULL, -- 'institute' or 'external'
    pnd_details         JSONB DEFAULT '{}'::JSONB,
    pnd_otp_code        VARCHAR(6),
    pnd_otp_generated_at TIMESTAMP,
    pnd_otp_expires_at  TIMESTAMP,
    pnd_otp_attempts    INT DEFAULT 0,
    pnd_verified        BOOLEAN DEFAULT FALSE,
    pnd_created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pnd_updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Ensure one pending registration per whatsapp/cedula combination
    CONSTRAINT uk_pending_whatsapp_identity UNIQUE (pnd_whatsapp, pnd_identity_number)
);

-- =====================================================
-- Table: cht_otp_verification_log
-- Audit log for OTP verification attempts
-- =====================================================
CREATE TABLE IF NOT EXISTS public.cht_otp_verification_log (
    ovl_id              SERIAL PRIMARY KEY,
    ovl_fk_pending      INT NOT NULL REFERENCES cht_pending_registrations(pnd_id) ON DELETE CASCADE,
    ovl_code_attempted  VARCHAR(6) NOT NULL,
    ovl_success         BOOLEAN NOT NULL,
    ovl_ip_address      INET,
    ovl_attempted_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Indexes for Performance
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_pending_reg_identity ON cht_pending_registrations(pnd_identity_number);
CREATE INDEX IF NOT EXISTS idx_pending_reg_whatsapp ON cht_pending_registrations(pnd_whatsapp);
CREATE INDEX IF NOT EXISTS idx_pending_reg_verified ON cht_pending_registrations(pnd_verified);
CREATE INDEX IF NOT EXISTS idx_pending_reg_otp_expires ON cht_pending_registrations(pnd_otp_expires_at);
CREATE INDEX IF NOT EXISTS idx_otp_log_pending ON cht_otp_verification_log(ovl_fk_pending);
CREATE INDEX IF NOT EXISTS idx_otp_log_attempted ON cht_otp_verification_log(ovl_attempted_at DESC);

-- =====================================================
-- Update Timestamp Trigger
-- =====================================================
DROP TRIGGER IF EXISTS tr_cht_pending_registrations_updated ON cht_pending_registrations;

CREATE OR REPLACE FUNCTION fn_update_pending_reg_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.pnd_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_cht_pending_registrations_updated
    BEFORE UPDATE ON cht_pending_registrations
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_pending_reg_timestamp();

-- =====================================================
-- Comments
-- =====================================================
COMMENT ON TABLE cht_pending_registrations IS 'Pending user registrations awaiting OTP verification';
COMMENT ON TABLE cht_otp_verification_log IS 'Audit log for OTP verification attempts';
COMMENT ON COLUMN cht_pending_registrations.pnd_user_type IS 'Type: institute (student/professor) or external';
COMMENT ON COLUMN cht_pending_registrations.pnd_otp_attempts IS 'Number of failed OTP verification attempts';
COMMENT ON COLUMN cht_pending_registrations.pnd_verified IS 'Whether OTP has been successfully verified';
