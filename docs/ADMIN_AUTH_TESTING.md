# Admin Authentication Testing Guide

## Prerequisites

### 1. Apply Database Migration

Run the admin authentication migration:

```bash
PGPASSWORD='lo0G4Rfaw7gtHw0wvpm4aqi4' psql -h localhost -p 5432 -U postgres -d chatbot_db -f db/07_admin_authentication.sql
```

### 2. Add JWT Configuration Parameter

Insert the JWT configuration parameter in the database:

```sql
INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
VALUES (
    'SECURITY',
    'JWT_CONFIG',
    '{
        "accessSecret": "your-super-secret-access-key-change-in-production",
        "refreshSecret": "your-super-secret-refresh-key-change-in-production",
        "accessExpiryHours": 1,
        "refreshExpiryHours": 168
    }'::jsonb,
    'JWT token configuration for admin authentication'
);
```

**Note:** Change the secrets in production! Use strong, randomly generated keys.

### 3. Create Test Admin User

Create a test admin user directly in the database:

```sql
-- Password: Admin123!
-- This is the bcrypt hash for "Admin123!"
INSERT INTO cht_admin_users (
    adm_username,
    adm_email,
    adm_password_hash,
    adm_name,
    adm_role,
    adm_permissions,
    adm_claims,
    adm_is_active
) VALUES (
    'admin',
    'admin@ists.edu.ec',
    '$2a$10$rHqDXE.8qV8WOXqGkM2WXuZ0V8y8vF0OzNqWYX.KqJxY0qWJxMQYC',
    'System Administrator',
    'super_admin',
    '["*"]'::jsonb,
    '{}'::jsonb,
    true
);
```

Or use the API endpoint (requires existing admin or direct DB access first):

```bash
curl -X POST http://localhost:8080/admin/users/create \
  -H "Content-Type: application/json" \
  -d '{
    "idSession": "test-session",
    "idRequest": "550e8400-e29b-41d4-a716-446655440000",
    "process": "admin-creation",
    "idDevice": "test-device",
    "deviceAddress": "127.0.0.1",
    "dateProcess": "2025-10-22T10:00:00Z",
    "username": "testadmin",
    "email": "testadmin@ists.edu.ec",
    "password": "SecurePass123!",
    "name": "Test Administrator",
    "role": "admin",
    "permissions": ["users.read", "parameters.read", "documents.manage"],
    "claims": {"department": "IT"}
  }'
```

## Testing the Authentication Flow

### 1. Login

```bash
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Operación exitosa",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
    "tokenType": "Bearer",
    "expiresIn": 3600,
    "expiresAt": "2025-10-22T11:00:00Z",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@ists.edu.ec",
      "name": "System Administrator",
      "role": "super_admin",
      "permissions": ["*"],
      "claims": {}
    }
  }
}
```

Save the `accessToken` and `refreshToken` for subsequent requests.

### 2. Test Protected Route (Example)

Once you have the access token, you can access protected routes:

```bash
curl -X POST http://localhost:8080/api/v1/parameters/get-all \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "idSession": "test-session",
    "idRequest": "550e8400-e29b-41d4-a716-446655440001",
    "process": "test-process",
    "idDevice": "test-device",
    "deviceAddress": "127.0.0.1",
    "dateProcess": "2025-10-22T10:00:00Z"
  }'
```

### 3. Refresh Token

After the access token expires (1 hour by default), use the refresh token:

```bash
curl -X POST http://localhost:8080/admin/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

**Expected Response:**
New access token and refresh token pair (token rotation).

### 4. Logout

Revoke the refresh token:

```bash
curl -X POST http://localhost:8080/admin/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refreshToken": "YOUR_REFRESH_TOKEN"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "code": "OK",
  "info": "Operación exitosa",
  "data": {
    "message": "Logout exitoso"
  }
}
```

## Security Features to Test

### 1. Account Locking

Try logging in with wrong password 5 times:

```bash
for i in {1..5}; do
  curl -X POST http://localhost:8080/admin/auth/login \
    -H "Content-Type: application/json" \
    -d '{
      "username": "admin",
      "password": "WrongPassword"
    }'
  echo "Attempt $i"
done
```

The account should be locked after 5 failed attempts.

### 2. Token Rotation

Use the same refresh token twice. The second use should fail and revoke the entire token family (security breach detection):

```bash
# First use - should succeed
curl -X POST http://localhost:8080/admin/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "YOUR_REFRESH_TOKEN"}'

# Second use - should fail and revoke family
curl -X POST http://localhost:8080/admin/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "YOUR_REFRESH_TOKEN"}'
```

### 3. Token Expiration

Wait for the access token to expire (1 hour) or modify the `accessExpiryHours` to 0.0003 (about 1 second) for testing:

```sql
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{accessExpiryHours}', '0.0003')
WHERE prm_code = 'JWT_CONFIG';
```

Then try using an expired token - it should be rejected.

## Verification Queries

### Check Auth Logs

```sql
SELECT
    log_username,
    log_action,
    log_status,
    log_ip_address,
    log_user_agent,
    log_details,
    log_created_at
FROM cht_auth_logs
ORDER BY log_created_at DESC
LIMIT 20;
```

### Check Refresh Tokens

```sql
SELECT
    rft_id,
    rft_admin_id,
    rft_token_family,
    rft_is_revoked,
    rft_revoked_reason,
    rft_expires_at,
    rft_created_at
FROM cht_refresh_tokens
ORDER BY rft_created_at DESC
LIMIT 10;
```

### Check Failed Attempts

```sql
SELECT
    adm_username,
    adm_failed_attempts,
    adm_is_locked,
    adm_last_login,
    adm_last_login_ip
FROM cht_admin_users
WHERE adm_username = 'admin';
```

## Troubleshooting

### "ERR_USER_NOT_FOUND"
- Check if admin user exists in `cht_admin_users` table
- Verify username is correct

### "ERR_INVALID_CREDENTIALS"
- Verify password is correct
- Check if account is locked (`adm_is_locked = true`)

### "ERR_ACCOUNT_LOCKED"
- Reset failed attempts: `UPDATE cht_admin_users SET adm_failed_attempts = 0, adm_is_locked = false WHERE adm_username = 'admin';`

### "ERR_INVALID_TOKEN"
- Token may be expired
- Token may be malformed
- JWT secrets may not match

### "ERR_TOKEN_REVOKED"
- Refresh token was already used (token rotation)
- Token family was revoked due to security breach
- Check `cht_refresh_tokens` table

## API Documentation

Once the server is running, view the complete API documentation at:

```
http://localhost:8080/docs
```

This includes all admin authentication endpoints with request/response schemas.
