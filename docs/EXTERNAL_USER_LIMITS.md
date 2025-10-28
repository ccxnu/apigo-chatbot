## External User Usage Limits - Configuration Guide

## Overview

The system now provides comprehensive controls to limit external user access to the chatbot. This prevents abuse and ensures fair resource allocation between institute members and external visitors.

---

## Available Limit Types

### 1. **Time-Based Message Limits**

| Limit Type | Parameter | Default | Description |
|------------|-----------|---------|-------------|
| **Daily** | `EXTERNAL_USER_DAILY_LIMIT` | 20 messages | Messages per day (resets at midnight) |
| **Weekly** | `EXTERNAL_USER_WEEKLY_LIMIT` | 100 messages | Messages per 7-day rolling period |
| **Monthly** | `EXTERNAL_USER_MONTHLY_LIMIT` | 300 messages | Messages per calendar month |
| **Lifetime** | `EXTERNAL_USER_TOTAL_LIMIT` | 1000 (disabled) | Total messages ever (requires enabling) |

### 2. **Rate Limiting**

| Parameter | Default | Description |
|-----------|---------|-------------|
| `EXTERNAL_USER_RATE_LIMIT` | 5 msg/min | Prevents spam by limiting messages per minute |

### 3. **Access Controls**

| Control | Parameter | Default | Description |
|---------|-----------|---------|-------------|
| **Registration** | `EXTERNAL_REGISTRATION_ENABLED` | `true` | Enable/disable external user registration |
| **Approval** | `EXTERNAL_USER_REQUIRE_APPROVAL` | `false` | Require admin approval before access |
| **Max Users** | `EXTERNAL_USER_MAX_COUNT` | 500 (disabled) | Maximum total external users allowed |
| **Access Hours** | `EXTERNAL_USER_ACCESS_HOURS` | Disabled | Restrict access to business hours |
| **Auto-Expiry** | `EXTERNAL_USER_EXPIRY_DAYS` | 30 days | Deactivate after inactivity |

---

## Configuration Examples

### Example 1: Conservative Limits (Recommended for Public Access)

```sql
-- Daily limit: 10 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 10, "period": "daily"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT';

-- Weekly limit: 50 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 50, "period": "weekly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT';

-- Monthly limit: 150 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 150, "period": "monthly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT';

-- Rate limit: 3 messages per minute
UPDATE cht_parameters
SET prm_data = '{"messages_per_minute": 3, "enabled": true}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_RATE_LIMIT';

-- Auto-expire after 15 days of inactivity
UPDATE cht_parameters
SET prm_data = '{"days": 15, "enabled": true}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_EXPIRY_DAYS';
```

### Example 2: Generous Limits (For Partners/Collaborators)

```sql
-- Daily limit: 50 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 50, "period": "daily"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT';

-- Weekly limit: 300 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 300, "period": "weekly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT';

-- Monthly limit: 1000 messages
UPDATE cht_parameters
SET prm_data = '{"limit": 1000, "period": "monthly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT';

-- No auto-expiry
UPDATE cht_parameters
SET prm_data = '{"days": 0, "enabled": false}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_EXPIRY_DAYS';
```

### Example 3: Strict Control (Approval Required)

```sql
-- Require admin approval for all external users
UPDATE cht_parameters
SET prm_data = '{"required": true}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_REQUIRE_APPROVAL';

-- Limit total external users to 100
UPDATE cht_parameters
SET prm_data = '{"limit": 100, "enabled": true}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_MAX_COUNT';

-- Restrict access to business hours (8 AM - 6 PM)
UPDATE cht_parameters
SET prm_data = '{
  "enabled": true,
  "start_hour": 8,
  "end_hour": 18,
  "timezone": "America/Guayaquil"
}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_ACCESS_HOURS';
```

### Example 4: Disable External Registration

```sql
-- Completely disable external user registration
UPDATE cht_parameters
SET prm_data = '{"enabled": false}'::jsonb
WHERE prm_code = 'EXTERNAL_REGISTRATION_ENABLED';
```

### Example 5: No Limits (Unlimited Access)

```sql
-- Set all limits to 0 (unlimited)
UPDATE cht_parameters
SET prm_data = '{"limit": 0, "period": "daily"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT';

UPDATE cht_parameters
SET prm_data = '{"limit": 0, "period": "weekly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_WEEKLY_LIMIT';

UPDATE cht_parameters
SET prm_data = '{"limit": 0, "period": "monthly"}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_MONTHLY_LIMIT';

-- Disable rate limiting
UPDATE cht_parameters
SET prm_data = '{"messages_per_minute": 0, "enabled": false}'::jsonb
WHERE prm_code = 'EXTERNAL_USER_RATE_LIMIT';
```

---

## How It Works

### Limit Checking Process

When an external user sends a message:

1. **System calls** `fn_check_external_user_limits(user_id, whatsapp)`

2. **Function checks** (in order):
   - ✅ Is user approved? (if approval required)
   - ✅ Rate limit (messages in last minute)
   - ✅ Total lifetime limit (if enabled)
   - ✅ Daily limit
   - ✅ Weekly limit
   - ✅ Monthly limit

3. **If any limit exceeded**:
   - Returns `allowed = FALSE` with specific error code
   - Bot sends appropriate message to user

4. **If all checks pass**:
   - Returns `allowed = TRUE` with remaining counts
   - Message is processed normally
   - User activity is updated via `sp_update_user_activity(user_id)`

### Activity Tracking

Each time an external user sends a message:

```sql
-- Update activity and increment counter
CALL sp_update_user_activity(success, code, user_id);
```

This updates:
- `usr_last_activity_at` → Current timestamp
- `usr_message_count` → Increments by 1

---

## User Experience

### When Daily Limit Reached

```
User: "¿Cuándo abren las inscripciones?"

Bot: 📊 Has alcanzado tu límite diario de mensajes (20/20).

Podrás enviar más mensajes mañana.

Como usuario externo, tienes las siguientes cuotas:
• Diario: 20 mensajes
• Semanal: 100 mensajes
• Mensual: 300 mensajes

Te quedan disponibles:
• Esta semana: 45 mensajes
• Este mes: 215 mensajes
```

### When Rate Limit Exceeded

```
User: "pregunta1"
User: "pregunta2"
User: "pregunta3"
User: "pregunta4"
User: "pregunta5"
User: "pregunta6"  ← Too fast!

Bot: ⏳ Estás enviando mensajes muy rápido.

Por favor espera un momento antes de enviar otro mensaje.

Límite: 5 mensajes por minuto
```

### When Approval Required

```
[User completes OTP verification]

Bot: ✅ Registro completado, Carlos Mendez!

⏳ Tu cuenta requiere aprobación del administrador.

Te notificaremos por WhatsApp cuando tu cuenta esté aprobada y puedas empezar a usar el asistente.

Gracias por tu paciencia.
```

### When Account Expired

```
User: [Tries to send message after 30 days inactive]

Bot: ⚠️ Tu cuenta ha expirado por inactividad (30 días sin uso).

Si deseas reactivar tu cuenta, por favor contacta al administrador del instituto.
```

---

## Admin Operations

### Check User Stats

```sql
-- Get usage statistics for a specific external user
SELECT * FROM fn_get_external_user_stats(123);

-- Returns:
-- daily_count   | 15
-- weekly_count  | 87
-- monthly_count | 245
-- total_count   | 1340
-- daily_limit   | 20
-- weekly_limit  | 100
-- monthly_limit | 300
-- last_activity | 2025-01-20 14:30:00
-- approved      | true
```

### Approve External User

```sql
-- Approve a pending external user
UPDATE cht_users
SET usr_approved = TRUE
WHERE usr_id = 123;
```

### Manually Reset User Limits

```sql
-- Reset message counter for a user
UPDATE cht_users
SET usr_message_count = 0
WHERE usr_id = 123;

-- Note: Daily/weekly/monthly counts are based on actual message timestamps,
-- so you would need to delete messages to reset those.
```

### Deactivate User Manually

```sql
-- Deactivate a specific external user
UPDATE cht_users
SET usr_active = FALSE
WHERE usr_id = 123;
```

### Run Automated Cleanup

```sql
-- Deactivate all expired external users (run via cron)
CALL sp_deactivate_expired_external_users(success, code, deactivated_count);

-- Check result
SELECT success, code, deactivated_count;
```

### View All External Users

```sql
-- Get list of all external users with stats
SELECT
    u.usr_id,
    u.usr_name,
    u.usr_email,
    u.usr_whatsapp,
    u.usr_message_count,
    u.usr_last_activity_at,
    u.usr_approved,
    u.usr_active,
    u.usr_created_at
FROM cht_users u
WHERE u.usr_rol = 'ROLE_EXTERNAL'
ORDER BY u.usr_created_at DESC;
```

### View External Users Approaching Limits

```sql
-- Find users close to daily limit (>80%)
WITH user_stats AS (
    SELECT * FROM fn_get_external_user_stats(usr_id)
    FROM cht_users
    WHERE usr_rol = 'ROLE_EXTERNAL'
    AND usr_active = TRUE
)
SELECT
    u.usr_id,
    u.usr_name,
    s.daily_count,
    s.daily_limit,
    ROUND((s.daily_count::NUMERIC / s.daily_limit) * 100, 1) as usage_pct
FROM cht_users u
JOIN user_stats s ON TRUE
WHERE s.daily_limit > 0
AND (s.daily_count::NUMERIC / s.daily_limit) > 0.8
ORDER BY usage_pct DESC;
```

---

## Monitoring & Alerts

### Dashboard Queries

```sql
-- Total external users
SELECT COUNT(*) as total_external_users
FROM cht_users
WHERE usr_rol = 'ROLE_EXTERNAL'
AND usr_active = TRUE;

-- External users registered today
SELECT COUNT(*) as new_today
FROM cht_users
WHERE usr_rol = 'ROLE_EXTERNAL'
AND usr_created_at >= CURRENT_DATE;

-- External users pending approval
SELECT COUNT(*) as pending_approval
FROM cht_users
WHERE usr_rol = 'ROLE_EXTERNAL'
AND usr_approved = FALSE
AND usr_active = TRUE;

-- External users who hit daily limit today
SELECT COUNT(*) as hit_daily_limit
FROM cht_users u
WHERE u.usr_rol = 'ROLE_EXTERNAL'
AND (
    SELECT COUNT(*)
    FROM cht_conversation_messages m
    JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
    WHERE c.cnv_fk_user = u.usr_id
    AND m.cvm_from_me = FALSE
    AND m.cvm_created_at >= CURRENT_DATE
) >= (
    SELECT (prm_data->>'limit')::INT
    FROM cht_parameters
    WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT'
);
```

---

## Institute Users vs External Users

| Feature | Institute Users | External Users |
|---------|----------------|----------------|
| **Registration** | Automatic (via AcademicOK) | Manual (provide email) |
| **Daily Limit** | Unlimited | Configurable (default: 20) |
| **Weekly Limit** | Unlimited | Configurable (default: 100) |
| **Monthly Limit** | Unlimited | Configurable (default: 300) |
| **Rate Limit** | None | 5 messages/minute |
| **Approval** | Not required | Optional (can be enabled) |
| **Auto-Expiry** | Never | After inactivity (default: 30 days) |
| **Access Hours** | 24/7 | Configurable (can restrict) |
| **Priority** | High | Standard |

---

## Migration Guide

### Install the Limits System

```bash
# Run the SQL files in order
psql -h localhost -U postgres -d chatbot_db -f db/08_external_user_limits.sql
psql -h localhost -U postgres -d chatbot_db -f db/09_external_user_limit_procedures.sql
```

This will:
1. ✅ Add all configuration parameters with default values
2. ✅ Add new columns to `cht_users` table (`usr_approved`, `usr_last_activity_at`, `usr_message_count`)
3. ✅ Create stored procedures for limit checking
4. ✅ Create functions for stats and cleanup

### Update Existing External Users

```sql
-- Set default values for existing external users
UPDATE cht_users
SET usr_approved = TRUE,
    usr_last_activity_at = COALESCE(usr_created_at, CURRENT_TIMESTAMP),
    usr_message_count = 0
WHERE usr_rol = 'ROLE_EXTERNAL'
AND usr_approved IS NULL;
```

---

## Recommended Settings by Use Case

### Use Case 1: Public Institute Chatbot

**Goal**: Allow public access but prevent abuse

```sql
-- Daily: 15 messages
-- Weekly: 75 messages
-- Monthly: 200 messages
-- Rate limit: 3 per minute
-- Auto-expire: 30 days
-- No approval required
```

### Use Case 2: Partner/Vendor Access

**Goal**: Provide reliable access to partners

```sql
-- Daily: 50 messages
-- Weekly: 300 messages
-- Monthly: 1000 messages
-- Rate limit: 10 per minute
-- Auto-expire: disabled
-- Approval: required
```

### Use Case 3: Trial/Demo Access

**Goal**: Time-limited trial for prospects

```sql
-- Daily: 20 messages
-- Total: 100 messages (enabled)
-- Auto-expire: 7 days
-- Access hours: business hours only
```

### Use Case 4: Emergency/Temporary Closure

**Goal**: Disable external access temporarily

```sql
-- Registration: disabled
-- All existing external users: active
```

---

## Troubleshooting

### User Claims Limit is Wrong

```sql
-- Check actual message count
SELECT
    COUNT(*) as messages_today,
    (SELECT (prm_data->>'limit')::INT FROM cht_parameters WHERE prm_code = 'EXTERNAL_USER_DAILY_LIMIT') as daily_limit
FROM cht_conversation_messages m
JOIN cht_conversations c ON c.cnv_id = m.cvm_fk_conversation
WHERE c.cnv_fk_user = (SELECT usr_id FROM cht_users WHERE usr_whatsapp = '+593999999999')
AND m.cvm_from_me = FALSE
AND m.cvm_created_at >= CURRENT_DATE;
```

### Limits Not Being Enforced

Check if parameters are active:

```sql
SELECT prm_code, prm_active, prm_data
FROM cht_parameters
WHERE prm_code LIKE 'EXTERNAL_USER%'
ORDER BY prm_code;
```

---

This completes the external user limits system! All parameters are configurable without code changes. 🎉
