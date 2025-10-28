# OTP-Based Registration System - Integration Guide

## Overview

The new OTP (One-Time Password) registration system adds email verification to the WhatsApp chatbot registration process. This prevents identity theft by ensuring that users can only register with their own cedula when they have access to the associated email address.

## Changes from Previous System

### Before (Old System)
1. User sends cedula via WhatsApp
2. System validates cedula with AcademicOK API
3. User is **immediately registered** and can start chatting
4. ❌ **Problem**: Anyone could register using someone else's cedula from a different phone

### After (New System with OTP)
1. User sends cedula via WhatsApp
2. System validates cedula with AcademicOK API
3. System generates 6-digit OTP code
4. System sends OTP to user's institutional email
5. User enters OTP code via WhatsApp
6. ✅ **Only if OTP matches**: User is registered and can start chatting
7. ✅ **Security**: Prevents unauthorized registration since user must have email access

## Architecture

### Database Layer

**New Tables:**
- `cht_pending_registrations` - Stores users awaiting OTP verification
- `cht_otp_verification_log` - Audit log for OTP verification attempts

**New Stored Procedures:**
- `sp_create_pending_registration` - Create/update pending registration with OTP
- `sp_delete_pending_registration` - Remove pending registration after successful registration
- `fn_verify_otp_code` - Verify OTP code and return user data
- `fn_get_pending_registration_by_whatsapp` - Get pending registration
- `fn_cleanup_expired_pending_registrations` - Clean up expired registrations

### Application Layer

**New Domain Entities** (`domain/registration.go`):
- `PendingRegistration` - Pending registration model
- `OTPVerificationResult` - OTP verification result
- `RegistrationRepository` - Repository interface
- `RegistrationUseCase` - Use case interface
- `OTPMailer` - Email sender interface

**New Repository** (`repository/registration_repository.go`):
- Implements `RegistrationRepository` interface
- Handles database operations for pending registrations

**New Use Case** (`usecase/registration_usecase.go`):
- `InitiateRegistration` - Start registration with OTP generation
- `VerifyAndRegister` - Verify OTP and complete registration
- `GetPendingRegistration` - Get pending registration status
- `ResendOTP` - Generate and send new OTP code

**New Mailer Service** (`internal/mailer/otp_mailer.go`):
- Implements `OTPMailer` interface
- Sends OTP emails via Tikee email service
- Customizable email template

**New WhatsApp Handler** (`internal/whatsapp/handlers/registration_handler.go`):
- Replaces `user_validation_handler.go`
- Handles entire OTP registration flow
- Manages user state (pending vs registered)

## Integration Steps

### 1. Run Database Migrations

Execute the SQL files in order:

```bash
psql -h localhost -U postgres -d chatbot_db -f db/05_registration_otp_tables.sql
psql -h localhost -U postgres -d chatbot_db -f db/06_registration_otp_procedures.sql
psql -h localhost -U postgres -d chatbot_db -f db/07_otp_registration_parameters.sql
```

This creates:
- New tables for pending registrations
- Stored procedures for OTP management
- Error codes and configuration parameters

### 2. Update Dependency Injection

In your `cmd/main.go` or dependency injection setup:

```go
import (
    "api-chatbot/domain"
    "api-chatbot/repository"
    "api-chatbot/usecase"
    "api-chatbot/internal/mailer"
    "api-chatbot/internal/whatsapp/handlers"
)

// Create repositories
regRepo := repository.NewRegistrationRepository(dal)
userRepo := repository.NewWhatsAppUserRepository(dal)

// Create mailer
tikeeURL := "http://20.84.48.225:5056/api/emails/enviarDirecto"
senderEmail := "automatizaciones@tikee.tech"
otpMailer := mailer.NewOTPMailer(httpClient, tikeeURL, senderEmail, paramCache, timeout)

// Create use cases
userUseCase := usecase.NewWhatsAppUserUseCase(userRepo, httpClient, paramCache, timeout)
regUseCase := usecase.NewRegistrationUseCase(regRepo, userRepo, userUseCase, otpMailer, paramCache, timeout)

// Create WhatsApp handlers
registrationHandler := handlers.NewRegistrationHandler(
    regUseCase,
    userUseCase,
    convUseCase,
    whatsappClient,
    paramCache,
    1, // priority
)

// Replace old UserValidationHandler with RegistrationHandler
dispatcher.AddHandler(registrationHandler)
```

### 3. Update Handler Priority

Replace the old `UserValidationHandler` with the new `RegistrationHandler`:

**Before:**
```go
userValidationHandler := handlers.NewUserValidationHandler(...)
dispatcher.AddHandler(userValidationHandler)
```

**After:**
```go
registrationHandler := handlers.NewRegistrationHandler(...)
dispatcher.AddHandler(registrationHandler)
```

### 4. Configure Email Service

Update the `TIKEE_EMAIL_SERVICE` parameter if needed:

```sql
UPDATE cht_parameters
SET prm_data = '{"url": "http://YOUR_TIKEE_URL/api/emails/enviarDirecto", "sender": "your-email@domain.com"}'::jsonb
WHERE prm_code = 'TIKEE_EMAIL_SERVICE';
```

### 5. Customize OTP Template (Optional)

To use a custom email template:

```sql
UPDATE cht_parameters
SET prm_data = '{
  "subject": "Your Custom Subject",
  "html": "<html>Your custom HTML template with [codigo_otp] placeholder</html>"
}'::jsonb
WHERE prm_code = 'EMAIL_OTP_TEMPLATE';
```

Available placeholders:
- `[nom_usuario]` - User's name
- `[nom_email]` - User's email
- `[codigo_otp]` - 6-digit OTP code
- `[fecha]` - Current date
- `[hora]` - Current time
- `[tipo_usuario]` - User type (miembro de la institución / usuario externo)

## User Flow

### For Institute Users (Students/Professors)

1. **User sends cedula**: `1234567890`
2. **System validates** with AcademicOK API
3. **System sends response**:
   ```
   ✅ ¡Hola Juan Pérez!

   He enviado un código de verificación de 6 dígitos a tu correo electrónico:
   📧 ju***@instituto.edu.ec

   Por favor, revisa tu bandeja de entrada (y también la carpeta de spam)
   y envíame el código para completar tu registro.

   El código expirará en 10 minutos.

   Si no recibes el correo, escribe "reenviar" para generar un nuevo código.
   ```

4. **User receives email** with 6-digit OTP (e.g., `123456`)

5. **User sends OTP**: `123456`

6. **System verifies and registers**:
   ```
   🎓 ¡Registro completado, Juan Pérez!

   Has sido registrado exitosamente como estudiante.

   Ahora puedes hacer preguntas sobre el instituto y recibir asistencia.
   ```

### For External Users

External users (not in AcademicOK) receive:
```
👤 No encontré tu cédula en nuestra base de datos institucional.

Actualmente, el registro está disponible solo para:
• Estudiantes del instituto
• Docentes del instituto

Si eres estudiante o docente, verifica que tu cédula sea correcta.

Si eres un visitante externo, por favor contacta con el departamento
de sistemas para obtener acceso.
```

## OTP Features

### Security Features
- **6-digit random code** generated using crypto/rand
- **10-minute expiration** (configurable)
- **Maximum 5 attempts** per OTP code
- **Audit logging** of all verification attempts
- **One-time use** - code invalidated after successful verification

### User Experience Features
- **Resend capability** - User can request new OTP by writing "reenviar"
- **Email masking** - Shows `ju***@instituto.edu.ec` for privacy
- **Clear instructions** - Guides user through each step
- **Helpful error messages** - Specific feedback for different failure scenarios

### Error Handling

The system handles various error scenarios:

| Error Code | Scenario | User Message |
|------------|----------|--------------|
| `ERR_USER_ALREADY_EXISTS` | User already registered | "Ya estás registrado en el sistema" |
| `ERR_IDENTITY_ALREADY_REGISTERED` | Cedula registered with different WhatsApp | "Esta cédula ya está registrada con otro número de WhatsApp" |
| `ERR_INVALID_OTP` | Wrong OTP code | "Código incorrecto. Por favor verifica e intenta nuevamente" |
| `ERR_OTP_EXPIRED` | OTP code expired | "El código ha expirado. Escribe 'reenviar' para generar un nuevo código" |
| `ERR_MAX_ATTEMPTS` | Too many failed attempts | "Has excedido el número máximo de intentos" |
| `ERR_IDENTITY_NOT_FOUND` | Cedula not in AcademicOK | "No pude validar tu cédula" |

## Configuration Parameters

### OTP Expiration Time
```sql
-- Default: 10 minutes
UPDATE cht_parameters
SET prm_data = '{"minutes": 15}'::jsonb
WHERE prm_code = 'OTP_EXPIRATION_MINUTES';
```

### Email Service URL
```sql
UPDATE cht_parameters
SET prm_data = '{
  "url": "http://20.84.48.225:5056/api/emails/enviarDirecto",
  "sender": "automatizaciones@tikee.tech"
}'::jsonb
WHERE prm_code = 'TIKEE_EMAIL_SERVICE';
```

## Maintenance

### Clean Up Expired Registrations

Run periodically (e.g., via cron job):

```sql
SELECT fn_cleanup_expired_pending_registrations();
-- Returns number of deleted records
```

This removes pending registrations older than 24 hours.

### Monitor OTP Verification Attempts

Check failed attempts:

```sql
SELECT
    p.pnd_whatsapp,
    p.pnd_name,
    p.pnd_email,
    COUNT(l.ovl_id) as failed_attempts,
    MAX(l.ovl_attempted_at) as last_attempt
FROM cht_pending_registrations p
JOIN cht_otp_verification_log l ON l.ovl_fk_pending = p.pnd_id
WHERE l.ovl_success = FALSE
GROUP BY p.pnd_id
HAVING COUNT(l.ovl_id) >= 3
ORDER BY failed_attempts DESC;
```

### View Pending Registrations

```sql
SELECT
    pnd_whatsapp,
    pnd_name,
    pnd_email,
    pnd_role,
    pnd_user_type,
    pnd_otp_attempts,
    pnd_otp_expires_at,
    pnd_created_at
FROM cht_pending_registrations
WHERE pnd_verified = FALSE
ORDER BY pnd_created_at DESC;
```

## Testing

### Manual Test Flow

1. Start WhatsApp chatbot
2. Send message from unregistered number: `1234567890`
3. Check email for OTP code
4. Send OTP code: `123456`
5. Verify registration completed
6. Try sending another message to confirm user is now registered

### Test Scenarios

- ✅ Valid cedula + correct OTP → Success
- ✅ Valid cedula + wrong OTP → Error message, retry allowed
- ✅ Expired OTP → Error message, resend option
- ✅ Resend request → New OTP generated
- ✅ Already registered user → Skip registration
- ✅ External user (invalid cedula) → Informative message

## Rollback Plan

If you need to rollback to the old system:

1. **Restore old handler**:
   ```go
   userValidationHandler := handlers.NewUserValidationHandler(...)
   dispatcher.AddHandler(userValidationHandler)
   ```

2. **Keep new tables** (data is preserved)

3. **Optional**: Drop new tables if not needed:
   ```sql
   DROP TABLE IF EXISTS cht_otp_verification_log CASCADE;
   DROP TABLE IF EXISTS cht_pending_registrations CASCADE;
   ```

## Troubleshooting

### Emails not being sent

Check Tikee service configuration:
```sql
SELECT prm_data FROM cht_parameters WHERE prm_code = 'TIKEE_EMAIL_SERVICE';
```

Verify service is accessible:
```bash
curl -X POST http://20.84.48.225:5056/api/emails/enviarDirecto \
  -H "Content-Type: application/json" \
  -d '{"test": "connectivity"}'
```

### OTP codes expiring too quickly

Increase expiration time:
```sql
UPDATE cht_parameters
SET prm_data = '{"minutes": 15}'::jsonb
WHERE prm_code = 'OTP_EXPIRATION_MINUTES';
```

### Users stuck in pending state

Force cleanup:
```sql
DELETE FROM cht_pending_registrations
WHERE pnd_whatsapp = '+593999999999'
AND pnd_verified = FALSE;
```

Then user can start registration again.

## Security Considerations

1. **OTP codes are cryptographically random** - Uses `crypto/rand` for generation
2. **Codes expire after 10 minutes** - Reduces window for brute force
3. **Maximum 5 attempts** - Prevents brute force attacks
4. **Audit logging** - All attempts logged with timestamps
5. **Email verification** - Proves ownership of institutional email
6. **One WhatsApp per cedula** - Prevents duplicate registrations

## Future Enhancements

Potential improvements for future iterations:

- [ ] SMS OTP as fallback for users without email access
- [ ] Rate limiting on OTP generation (prevent spam)
- [ ] Admin panel to manually approve external users
- [ ] Multi-language support for messages
- [ ] Biometric verification (future)
- [ ] Integration with institutional SSO

## Support

For issues or questions, contact:
- Development Team: [your-email@domain.com]
- Documentation: `/docs/OTP_REGISTRATION_INTEGRATION.md`
