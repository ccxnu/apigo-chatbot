# Complete Registration Flows - All Scenarios

## Overview

The system now supports THREE types of users:
1. **Institute Students** (found in AcademicOK apidatospersona)
2. **Institute Professors** (found in AcademicOK apidatosdocente)
3. **External Users** (NOT in AcademicOK - visitors, partners, etc.)

All require **email OTP verification** to complete registration.

---

## Flow 1: Institute Student Registration ✅

### Step-by-Step

**1. User sends cedula**
```
User: "1234567890"
```

**2. System validates with AcademicOK**
- Calls `apidatospersona` API
- Student found with `careras` array → Role = "ROLE_STUDENT"
- Gets: name, email from institute database

**3. System creates pending registration**
- Generates 6-digit OTP
- Stores in `cht_pending_registrations`
- Sends email to institutional email

**4. Bot confirms**
```
Bot: ✅ ¡Hola Juan Pérez!

He enviado un código de verificación de 6 dígitos a tu correo electrónico:
📧 ju***@instituto.edu.ec

Por favor, revisa tu bandeja de entrada (y también la carpeta de spam)
y envíame el código para completar tu registro.

El código expirará en 10 minutos.

Si no recibes el correo, escribe "reenviar" para generar un nuevo código.
```

**5. User receives email with OTP** (e.g., `456789`)

**6. User sends OTP**
```
User: "456789"
```

**7. System verifies & registers**
- Verifies OTP in database
- Creates user in `cht_users` with `usr_rol = 'ROLE_STUDENT'`
- Creates conversation in `cht_conversations`
- Deletes pending registration

**8. Bot welcomes user**
```
Bot: 🎓 ¡Registro completado, Juan Pérez!

Has sido registrado exitosamente como estudiante.

Ahora puedes hacer preguntas sobre el instituto y recibir asistencia.
```

---

## Flow 2: Institute Professor Registration ✅

### Step-by-Step

**1. User sends cedula**
```
User: "0987654321"
```

**2. System validates with AcademicOK**
- Calls `apidatospersona` API → No `careras` array
- Calls `apidatosdocente` API → Professor found
- Role = "ROLE_PROFESSOR"
- Gets: name, email from institute database

**3-8. Same as Student Flow**
- OTP generation, email, verification, registration
- **Only difference**: Role emoji changes to 👨‍🏫

```
Bot: 👨‍🏫 ¡Registro completado, María García!

Has sido registrado exitosamente como docente.

Ahora puedes hacer preguntas sobre el instituto y recibir asistencia.
```

---

## Flow 3: External User Registration 🆕✅

### Step-by-Step

**1. User sends cedula**
```
User: "9999999999"
```

**2. System validates with AcademicOK**
- Calls `apidatospersona` API → NOT found
- Calls `apidatosdocente` API → NOT found
- Returns: "ERR_EXTERNAL_USER_INFO_REQUIRED"

**3. System creates incomplete pending registration**
- Creates record in `cht_pending_registrations`
- `pnd_user_type = 'external'`
- `pnd_email = ''` (empty - to be collected)
- `pnd_name = ''` (empty - to be collected)
- **NO OTP yet** - will generate after email is provided

**4. Bot asks for name and email**
```
Bot: 👤 No encontré tu cédula en nuestra base de datos institucional.

Sin embargo, puedes registrarte como usuario externo.

Por favor envíame tu información en el siguiente formato:

*Nombre Completo / correo@email.com*

Ejemplo:
Juan Pérez / juan.perez@gmail.com

Recibirás un código de verificación en ese correo para completar tu registro.
```

**5. User provides name and email**
```
User: "Carlos Mendez / carlos.mendez@gmail.com"
```

**6. System extracts and validates**
- Regex extracts email: `carlos.mendez@gmail.com`
- Splits by "/" extracts name: `Carlos Mendez`
- Validates name has at least 2 words ✅

**7. System completes external registration**
- Updates pending registration with name and email
- Generates 6-digit OTP
- Sends email to provided email
- `pnd_role = 'ROLE_EXTERNAL'`

**8. Bot confirms**
```
Bot: ✅ ¡Hola Carlos Mendez!

He enviado un código de verificación de 6 dígitos a tu correo electrónico:
📧 ca***@gmail.com

Por favor, revisa tu bandeja de entrada (y también la carpeta de spam)
y envíame el código para completar tu registro.

El código expirará en 10 minutos.

Si no recibes el correo, escribe "reenviar" para generar un nuevo código.
```

**9. User receives email with OTP** (e.g., `789012`)

**10. User sends OTP**
```
User: "789012"
```

**11. System verifies & registers**
- Verifies OTP
- Creates user with `usr_rol = 'ROLE_EXTERNAL'`
- User can now access chatbot

**12. Bot welcomes external user**
```
Bot: 👤 ¡Registro completado, Carlos Mendez!

Has sido registrado exitosamente como usuario externo.

Ahora puedes hacer preguntas sobre el instituto y recibir asistencia.
```

---

## Edge Case: Registered User Sends "reenviar" ❌→✅

### Scenario
```
User: "reenviar"  [Already registered user]
```

### What Happens

**1. RegistrationHandler checks**
```go
func (h *RegistrationHandler) Match(ctx, msg) bool {
    result := h.userUseCase.GetUserByWhatsApp(ctx, msg.From)
    if result.Success && result.Data != nil {
        return false  // ❌ Handler doesn't match
    }
    return true
}
```

**2. Handler doesn't match → Falls through to next handler**
- Message goes to CommandHandler or RAGHandler
- User gets normal response (RAG answer or "command not found")

**3. No error or confusion** ✅
- System gracefully handles registered users
- They can use chatbot normally

---

## Edge Case: Invalid Email Format

### Scenario
```
User: "Carlos Mendez email gmail.com"  [Missing @]
```

### What Happens

**1. System tries to extract email**
```go
extractEmail(text) → ""  // No match
```

**2. Bot reminds user of format**
```
Bot: 👤 No encontré tu cédula en nuestra base de datos institucional.

Sin embargo, puedes registrarte como usuario externo.

Por favor envíame tu información en el siguiente formato:

*Nombre Completo / correo@email.com*

Ejemplo:
Juan Pérez / juan.perez@gmail.com
```

---

## Edge Case: Name Without Slash

### Scenario
```
User: "Carlos Mendez carlos@gmail.com"  [No slash separator]
```

### What Happens

**1. System tries to extract name**
```go
extractName(text) → ""  // No "/" found
```

**2. Bot reminds user of format** (same as above)

---

## Edge Case: OTP Expiration During External Registration

### Scenario
```
External user provides email → receives OTP → waits 15 minutes → sends OTP
```

### What Happens

**1. User sends expired OTP**
```
User: "789012"  [Code is > 10 minutes old]
```

**2. System verifies**
```sql
SELECT * FROM fn_verify_otp_code(...)
-- Returns: success=FALSE, code='ERR_OTP_EXPIRED'
```

**3. Bot responds**
```
Bot: ⏰ El código ha expirado. Escribe 'reenviar' para generar un nuevo código.
```

**4. User requests new OTP**
```
User: "reenviar"
```

**5. System regenerates OTP**
- Generates new 6-digit code
- Updates `pnd_otp_code` and `pnd_otp_expires_at`
- Sends new email
- Old code is invalid

---

## Comparison Table

| Feature | Institute Users | External Users |
|---------|----------------|----------------|
| **Cedula validation** | ✅ Via AcademicOK API | ❌ Not in system |
| **Email source** | From AcademicOK database | User provides |
| **Name source** | From AcademicOK database | User provides |
| **Role** | ROLE_STUDENT or ROLE_PROFESSOR | ROLE_EXTERNAL |
| **OTP required** | ✅ Yes | ✅ Yes |
| **OTP timing** | After cedula validation | After user provides email |
| **Email format** | @instituto.edu.ec | Any valid email |
| **Steps required** | 1. Cedula<br/>2. OTP | 1. Cedula<br/>2. Name + Email<br/>3. OTP |

---

## Database States

### Institute User - After Cedula
```sql
-- cht_pending_registrations
pnd_identity_number | pnd_whatsapp  | pnd_name     | pnd_email              | pnd_role        | pnd_user_type | pnd_otp_code
--------------------|---------------|--------------|------------------------|-----------------|---------------|-------------
1234567890          | +593999999999 | Juan Pérez   | juan@instituto.edu.ec  | ROLE_STUDENT    | institute     | 456789
```

### External User - After Cedula (Before Email)
```sql
-- cht_pending_registrations
pnd_identity_number | pnd_whatsapp  | pnd_name | pnd_email | pnd_role        | pnd_user_type | pnd_otp_code
--------------------|---------------|----------|-----------|-----------------|---------------|-------------
9999999999          | +593999999999 | (empty)  | (empty)   | ROLE_EXTERNAL   | external      | (empty)
```

### External User - After Email Provided
```sql
-- cht_pending_registrations
pnd_identity_number | pnd_whatsapp  | pnd_name       | pnd_email           | pnd_role        | pnd_user_type | pnd_otp_code
--------------------|---------------|----------------|---------------------|-----------------|---------------|-------------
9999999999          | +593999999999 | Carlos Mendez  | carlos@gmail.com    | ROLE_EXTERNAL   | external      | 789012
```

### After Registration Complete (All Users)
```sql
-- cht_users
usr_identity_number | usr_whatsapp  | usr_name       | usr_email           | usr_rol
--------------------|---------------|----------------|---------------------|----------------
1234567890          | +593999999999 | Juan Pérez     | juan@instituto.edu  | ROLE_STUDENT
9999999999          | +593888888888 | Carlos Mendez  | carlos@gmail.com    | ROLE_EXTERNAL

-- cht_pending_registrations (deleted after registration)
(empty)
```

---

## Security Benefits

### For Institute Users
- ✅ Cedula validated against official AcademicOK database
- ✅ Email verified (must have access to institutional email)
- ✅ Cannot register with someone else's cedula + phone combination

### For External Users
- ✅ Email verified (must have access to provided email)
- ✅ Cannot impersonate institute members (cedula not in database)
- ✅ Clearly marked as "external" for auditing purposes
- ✅ Same OTP security as institute users

---

## Error Handling Summary

| Error Code | When | User Message |
|-----------|------|--------------|
| `ERR_USER_ALREADY_EXISTS` | User tries to register again | "Ya estás registrado" |
| `ERR_IDENTITY_ALREADY_REGISTERED` | Cedula registered with different WhatsApp | "Esta cédula ya está registrada" |
| `ERR_INVALID_OTP` | Wrong OTP code | "Código incorrecto" |
| `ERR_OTP_EXPIRED` | OTP > 10 minutes old | "El código ha expirado" |
| `ERR_MAX_ATTEMPTS` | > 5 failed attempts | "Máximo de intentos excedido" |
| `ERR_NO_PENDING_REGISTRATION` | No pending registration | "No tienes un registro pendiente" |
| `ERR_EXTERNAL_USER_INFO_REQUIRED` | Cedula not in AcademicOK | Request name + email |

---

## Testing Checklist

### Institute Student
- [ ] Valid cedula → OTP sent → Correct OTP → Registered ✅
- [ ] Valid cedula → OTP sent → Wrong OTP → Error + retry
- [ ] Valid cedula → OTP sent → Expired OTP → Resend option
- [ ] Valid cedula → OTP sent → "reenviar" → New OTP

### Institute Professor
- [ ] Valid professor cedula → OTP sent → Registered as professor ✅
- [ ] Same OTP flow as students

### External User
- [ ] Invalid cedula → Request email → Provide name/email → OTP sent → Registered ✅
- [ ] Invalid cedula → Request email → Invalid format → Re-ask
- [ ] Invalid cedula → Request email → Valid format → OTP → Verify → Success

### Edge Cases
- [ ] Registered user sends "reenviar" → Falls through to other handlers ✅
- [ ] External user provides only email (no name) → Re-ask
- [ ] External user provides malformed email → Re-ask
- [ ] Any user exceeds 5 OTP attempts → Must request new OTP

---

This completes ALL registration scenarios including external users! 🎉
