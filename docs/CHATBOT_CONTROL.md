# Chatbot Control - Hot Reload Guide

This guide explains how to enable/disable the WhatsApp chatbot without restarting the application.

## Overview

The chatbot has two levels of control:

1. **`WHATSAPP_CONFIG.enabled`** - Determines if WhatsApp service starts when app launches (requires restart)
2. **`CHATBOT_ACTIVE`** - Enables/disables message processing at runtime (hot-reloadable, **use this!**)

## Hot-Reloadable Control (Recommended)

### Deactivate Chatbot (Maintenance Mode)

```sql
-- Step 1: Update database
UPDATE cht_parameters
SET prm_data = '{"active": false}'::jsonb
WHERE prm_code = 'CHATBOT_ACTIVE';
```

```bash
-- Step 2: Reload cache (takes effect immediately)
curl -X POST http://localhost:8080/api/v1/parameters/reload-cache \
  -H "Content-Type: application/json" \
  -d '{
    "idSession": "admin",
    "idRequest": "550e8400-e29b-41d4-a716-446655440000",
    "process": "chatbot-control",
    "idDevice": "admin-device",
    "publicIp": "127.0.0.1",
    "dateProcess": "2025-10-22T00:00:00Z"
  }'
```

**What happens:**
- ✅ WhatsApp connection stays active
- ✅ Users receive: "🔧 El chatbot está temporalmente desactivado por mantenimiento. Por favor, intenta más tarde."
- ✅ All incoming messages are acknowledged but not processed
- ✅ Takes effect in ~1 second

### Reactivate Chatbot

```sql
-- Step 1: Update database
UPDATE cht_parameters
SET prm_data = '{"active": true}'::jsonb
WHERE prm_code = 'CHATBOT_ACTIVE';
```

```bash
-- Step 2: Reload cache
curl -X POST http://localhost:8080/api/v1/parameters/reload-cache \
  -H "Content-Type: application/json" \
  -d '{
    "idSession": "admin",
    "idRequest": "550e8400-e29b-41d4-a716-446655440001",
    "process": "chatbot-control",
    "idDevice": "admin-device",
    "publicIp": "127.0.0.1",
    "dateProcess": "2025-10-22T00:00:00Z"
  }'
```

### Customize Deactivation Message

```sql
UPDATE cht_parameters
SET prm_data = '{"message": "🚧 Sistema en mantenimiento hasta las 18:00. Disculpa las molestias."}'::jsonb
WHERE prm_code = 'CHATBOT_DEACTIVATED_MESSAGE';
```

Then reload the cache as shown above.

## Helper Script

Create a bash script for easy control:

```bash
#!/bin/bash
# chatbot-control.sh

DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"
DB_PASS="lo0G4Rfaw7gtHw0wvpm4aqi4"
DB_NAME="chatbot_db"
API_URL="http://localhost:8080"

function deactivate() {
    echo "Deactivating chatbot..."
    docker exec -i cnt_postgres psql -U $DB_USER -d $DB_NAME <<EOF
UPDATE cht_parameters
SET prm_data = '{"active": false}'::jsonb
WHERE prm_code = 'CHATBOT_ACTIVE';
EOF
    reload_cache
    echo "✅ Chatbot deactivated"
}

function activate() {
    echo "Activating chatbot..."
    docker exec -i cnt_postgres psql -U $DB_USER -d $DB_NAME <<EOF
UPDATE cht_parameters
SET prm_data = '{"active": true}'::jsonb
WHERE prm_code = 'CHATBOT_ACTIVE';
EOF
    reload_cache
    echo "✅ Chatbot activated"
}

function reload_cache() {
    echo "Reloading parameter cache..."
    curl -s -X POST $API_URL/api/v1/parameters/reload-cache \
      -H "Content-Type: application/json" \
      -d '{
        "idSession": "admin",
        "idRequest": "'$(uuidgen)'",
        "process": "chatbot-control",
        "idDevice": "admin-device",
        "publicIp": "127.0.0.1",
        "dateProcess": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
      }' | jq .
}

function status() {
    echo "Current chatbot status:"
    docker exec -i cnt_postgres psql -U $DB_USER -d $DB_NAME <<EOF
SELECT prm_code, prm_data->>'active' as active
FROM cht_parameters
WHERE prm_code = 'CHATBOT_ACTIVE';
EOF
}

case "$1" in
    on|activate)
        activate
        ;;
    off|deactivate)
        deactivate
        ;;
    status)
        status
        ;;
    reload)
        reload_cache
        ;;
    *)
        echo "Usage: $0 {on|off|status|reload}"
        exit 1
        ;;
esac
```

**Usage:**
```bash
chmod +x chatbot-control.sh

./chatbot-control.sh off      # Deactivate chatbot
./chatbot-control.sh on       # Activate chatbot
./chatbot-control.sh status   # Check current status
./chatbot-control.sh reload   # Just reload cache
```

## When to Use What

| Scenario | Solution | Requires Restart |
|----------|----------|-----------------|
| **Maintenance/Updates** | `CHATBOT_ACTIVE = false` | ❌ No |
| **Emergency disable** | `CHATBOT_ACTIVE = false` | ❌ No |
| **Change LLM model** | Update `LLM_CONFIG` + reload | ❌ No |
| **Change any RAG setting** | Update parameter + reload | ❌ No |
| **Change WhatsApp phone** | `WHATSAPP_CONFIG.enabled` | ✅ Yes |
| **Initial setup** | `WHATSAPP_CONFIG.enabled` | ✅ Yes |

## Architecture

```
User sends WhatsApp message
    ↓
WhatsApp Client receives
    ↓
MessageDispatcher.Dispatch()
    ↓
    ├─→ Check CHATBOT_ACTIVE (from cache)
    │   ├─→ If false: Send deactivation message, stop
    │   └─→ If true: Continue processing
    ↓
Handler routing (UserValidation → Command → RAG)
```

The parameter cache is checked on **every message**, so changes take effect immediately after reload.

## Best Practices

1. **Before deploying updates:**
   ```bash
   ./chatbot-control.sh off
   # Deploy your changes
   # Restart app if needed
   ./chatbot-control.sh on
   ```

2. **Testing new configurations:**
   ```bash
   # Disable bot
   ./chatbot-control.sh off

   # Update parameters in database
   # Test with a specific user

   # Re-enable when ready
   ./chatbot-control.sh on
   ```

3. **Emergency shutdown:**
   ```bash
   ./chatbot-control.sh off
   ```
   No restart needed - instant effect!

## Troubleshooting

**Q: I disabled the bot but users still get responses**
- A: Make sure you ran the reload-cache endpoint after updating the database

**Q: How do I know if the cache was reloaded?**
- A: Check the API response - it should return `{"success": true, "data": {"count": N}}`

**Q: Can I schedule automatic enable/disable?**
- A: Yes! Use cron:
  ```cron
  # Disable at 2 AM for maintenance
  0 2 * * * /path/to/chatbot-control.sh off

  # Re-enable at 3 AM
  0 3 * * * /path/to/chatbot-control.sh on
  ```
