# OTP Registration System - Implementation Summary

## Overview

Complete OTP-based email verification system for WhatsApp chatbot registration. Prevents identity theft by requiring email verification before completing user registration.

## Files Created

### Database Layer

1. **`db/05_registration_otp_tables.sql`**
   - Tables: `cht_pending_registrations`, `cht_otp_verification_log`
   - Indexes and triggers

2. **`db/06_registration_otp_procedures.sql`**
   - Stored procedures following your standard (OUT success, OUT code first)
   - `sp_create_pending_registration(OUT success, OUT code, OUT o_pending_id, IN p_identity_number, ...)`
   - `sp_delete_pending_registration(OUT success, OUT code, IN p_pending_id)`
   - `fn_verify_otp_code(...)` - Returns table with verification result
   - `fn_get_pending_registration_by_whatsapp(...)` - Returns pending registration
   - `fn_cleanup_expired_pending_registrations()` - Maintenance function

3. **`db/07_otp_registration_parameters.sql`**
   - 13 error codes
   - Configuration parameters
   - Role definitions

### Application Layer

4. **`domain/registration.go`**
   - Domain entities and interfaces
   - Follows your Result pattern

5. **`repository/registration_repository.go`**
   - Uses DAL with `QueryRows` and `ExecProc`
   - Follows your repository standards

6. **`usecase/registration_usecase.go`**
   - Business logic for OTP flow
   - Crypto-secure OTP generation

7. **`internal/mailer/otp_mailer.go`**
   - Sends emails via Tikee service
   - Customizable HTML templates

8. **`internal/whatsapp/handlers/registration_handler.go`**
   - Complete WhatsApp handler
   - Replaces `user_validation_handler.go`

### Documentation

9. **`docs/OTP_REGISTRATION_INTEGRATION.md`**
   - Integration guide
   - Configuration instructions

10. **`docs/OTP_IMPLEMENTATION_SUMMARY.md`** (this file)
    - Quick reference

## Key Standards Followed

### Stored Procedures
✅ **OUT parameters first**: `OUT success BOOLEAN, OUT code VARCHAR, OUT o_xxx, IN p_xxx`
✅ **IN parameters with defaults at end**: `IN p_optional DEFAULT NULL`
✅ **Naming**: `sp_*` for procedures, `fn_*` for functions
✅ **Error handling**: Always set `success := FALSE` and `code := 'ERR_*'` in EXCEPTION block
✅ **RAISE NOTICE**: Log errors with `RAISE NOTICE 'Error: %', SQLERRM`

### Repository Layer
✅ **DAL usage**: `dal.ExecProc[ResultType]()`and `dal.QueryRows[Type]()`
✅ **Error wrapping**: `fmt.Errorf("failed to execute %s: %w", procName, err)`
✅ **Naming**: `sp*` constants for procedures, `fn*` for functions
✅ **Result checking**: Check `!result.Success` for procedure results

### Use Case Layer
✅ **Context timeout**: Create context with timeout from `uc.contextTimeout`
✅ **Result pattern**: Return `d.Result[T]` with `d.Success(data)` or `d.Error[T](uc.paramCache, code)`
✅ **Logging**: Use `logger.LogError()`, `logger.LogWarn()`, `logger.LogInfo()`
✅ **Error codes**: Return error codes from parameter cache

## Integration Checklist

### 1. Database Setup
```bash
# Run migrations in order
psql -h localhost -U postgres -d chatbot_db -f db/05_registration_otp_tables.sql
psql -h localhost -U postgres -d chatbot_db -f db/06_registration_otp_procedures.sql
psql -h localhost -U postgres -d chatbot_db -f db/07_otp_registration_parameters.sql
```

### 2. Dependency Injection (cmd/main.go or similar)
```go
// Repositories
regRepo := repository.NewRegistrationRepository(dal)
userRepo := repository.NewWhatsAppUserRepository(dal)

// Mailer
tikeeURL := "http://20.84.48.225:5056/api/emails/enviarDirecto"
senderEmail := "automatizaciones@tikee.tech"
otpMailer := mailer.NewOTPMailer(httpClient, tikeeURL, senderEmail, paramCache, timeout)

// Use Cases
userUseCase := usecase.NewWhatsAppUserUseCase(userRepo, httpClient, paramCache, timeout)
regUseCase := usecase.NewRegistrationUseCase(regRepo, userRepo, userUseCase, otpMailer, paramCache, timeout)
convUseCase := usecase.NewConversationUseCase(...) // existing

// Handler
registrationHandler := handlers.NewRegistrationHandler(
    regUseCase,
    userUseCase,
    convUseCase,
    whatsappClient,
    paramCache,
    1, // priority
)

// Replace old handler
// OLD: dispatcher.AddHandler(userValidationHandler)
// NEW:
dispatcher.AddHandler(registrationHandler)
```

### 3. Configuration
Update parameters if needed:
```sql
-- OTP expiration (default 10 minutes)
UPDATE cht_parameters
SET prm_data = '{"minutes": 15}'::jsonb
WHERE prm_code = 'OTP_EXPIRATION_MINUTES';

-- Email service URL
UPDATE cht_parameters
SET prm_data = '{
  "url": "http://20.84.48.225:5056/api/emails/enviarDirecto",
  "sender": "automatizaciones@tikee.tech"
}'::jsonb
WHERE prm_code = 'TIKEE_EMAIL_SERVICE';
```

### 4. Testing
1. Send cedula from unregistered WhatsApp
2. Check email for OTP code
3. Send OTP code
4. Verify successful registration
5. Test error scenarios (wrong OTP, expired OTP, resend)

## User Flow

```
User → Sends cedula (1234567890)
     ↓
System → Validates with AcademicOK API
       → Generates 6-digit OTP
       → Sends email to user's institutional email
     ↓
User → Receives email with OTP (e.g., 456789)
     → Sends OTP via WhatsApp
     ↓
System → Verifies OTP
       → Creates user account
       → Creates conversation
       → Sends welcome message
     ↓
User → Can now chat with bot ✅
```

## Security Features

- ✅ Crypto-random 6-digit OTP
- ✅ 10-minute expiration (configurable)
- ✅ Maximum 5 attempts per OTP
- ✅ Audit logging of all attempts
- ✅ One-time use (invalidated after success)
- ✅ Email ownership verification
- ✅ Prevents duplicate registrations

## Maintenance

### Cleanup Expired Registrations
```sql
-- Run periodically (e.g., daily cron job)
SELECT fn_cleanup_expired_pending_registrations();
```

### Monitor Failed Attempts
```sql
SELECT
    p.pnd_whatsapp,
    p.pnd_email,
    COUNT(l.ovl_id) as failed_attempts
FROM cht_pending_registrations p
JOIN cht_otp_verification_log l ON l.ovl_fk_pending = p.pnd_id
WHERE l.ovl_success = FALSE
GROUP BY p.pnd_id
HAVING COUNT(l.ovl_id) >= 3;
```

### View Pending Registrations
```sql
SELECT
    pnd_whatsapp,
    pnd_name,
    pnd_email,
    pnd_otp_attempts,
    pnd_otp_expires_at,
    pnd_created_at
FROM cht_pending_registrations
WHERE pnd_verified = FALSE
ORDER BY pnd_created_at DESC;
```

## Error Codes

| Code | Description |
|------|-------------|
| `ERR_USER_ALREADY_EXISTS` | User already registered |
| `ERR_IDENTITY_ALREADY_REGISTERED` | Cedula registered with different WhatsApp |
| `ERR_INVALID_OTP` | Wrong OTP code |
| `ERR_OTP_EXPIRED` | OTP expired (> 10 minutes) |
| `ERR_MAX_ATTEMPTS` | Too many failed attempts (> 5) |
| `ERR_NO_PENDING_REGISTRATION` | No pending registration found |
| `ERR_EXTERNAL_USER_INFO_REQUIRED` | User not in AcademicOK database |
| `ERR_IDENTITY_NOT_FOUND` | Cedula validation failed |

## Differences from Old System

### Old System ❌
- User sends cedula
- Validates with API
- **Immediately registered** ← Security risk!
- Anyone could register with someone else's cedula

### New System ✅
- User sends cedula
- Validates with API
- Generates OTP → sends to email
- User must enter OTP from email
- **Only then registered** ← Secure!
- Proves user has access to institutional email

## File Removal

The following file can be **deprecated** (but keep for rollback):
- `internal/whatsapp/handlers/user_validation_handler.go`

The new `registration_handler.go` replaces it with OTP functionality.

## Rollback Plan

If needed, rollback is simple:

1. **In code**: Switch back to old handler
   ```go
   // Replace
   registrationHandler := handlers.NewRegistrationHandler(...)
   // With
   userValidationHandler := handlers.NewUserValidationHandler(...)
   ```

2. **Keep database tables** (data preserved for future use)

3. **Optional**: Drop tables if not needed
   ```sql
   DROP TABLE IF EXISTS cht_otp_verification_log CASCADE;
   DROP TABLE IF EXISTS cht_pending_registrations CASCADE;
   ```

## Support

For detailed integration steps, see:
- `/docs/OTP_REGISTRATION_INTEGRATION.md` - Full integration guide

For questions or issues, contact the development team.
