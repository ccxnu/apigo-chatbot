# Admin Chat Panel - Implementation Guide

## Overview

WhatsApp-like admin panel for managing conversations without needing WhatsApp Web.

## Features Implemented

### ✅ Database Layer (Complete)
- [x] User blocking fields (usr_blocked, usr_blocked_at, usr_blocked_by, usr_block_reason)
- [x] Conversation management (cnv_blocked, cnv_admin_intervened, cnv_unread_count, cnv_temporary)
- [x] Message tracking (cvm_sender_type, cvm_admin_id, cvm_read, cvm_read_at)
- [x] Stored procedures for all operations
- [x] Indexes for performance

### 🔄 Backend Layer (In Progress)
- [ ] Domain types (AdminConversationList, AdminMessage, etc.)
- [ ] Repository layer
- [ ] Use case layer
- [ ] API endpoints

### ⏳ Real-time Layer (Planned)
- [ ] Server-Sent Events (SSE) for new messages
- [ ] WebSocket alternative

## API Endpoints Design

### 1. **GET /admin/conversations**
List all conversations (WhatsApp-like)

**Query Parameters:**
- `filter`: `all` | `unread` | `blocked` | `active` (default: `all`)
- `limit`: int (default: 50)
- `offset`: int (default: 0)

**Response:**
```json
{
  "success": true,
  "data": {
    "conversations": [
      {
        "id": 1,
        "chatId": "593959006203@s.whatsapp.net",
        "phoneNumber": "+593959006203",
        "contactName": "Juan Pérez",
        "isGroup": false,
        "lastMessageAt": "2025-10-22T18:30:00Z",
        "messageCount": 45,
        "unreadCount": 3,
        "blocked": false,
        "adminIntervened": false,
        "temporary": false,
        "user": {
          "id": 10,
          "name": "Juan Pérez",
          "identityNumber": "1105608424",
          "role": "ROLE_STUDENT",
          "blocked": false
        },
        "lastMessagePreview": "¿Cuándo son las matriculas?",
        "lastMessageFromMe": false
      }
    ],
    "total": 123,
    "hasMore": true
  }
}
```

### 2. **GET /admin/conversations/:id/messages**
Get conversation history

**Query Parameters:**
- `limit`: int (default: 100)

**Response:**
```json
{
  "success": true,
  "data": {
    "conversationId": 1,
    "messages": [
      {
        "id": 1,
        "messageId": "3EB0C767D4B7F84C097C",
        "fromMe": false,
        "senderName": "Juan Pérez",
        "senderType": "user",
        "messageType": "text",
        "body": "Hola, ¿cuándo son las matrículas?",
        "timestamp": 1729625400,
        "read": true,
        "createdAt": "2025-10-22T18:30:00Z"
      },
      {
        "id": 2,
        "messageId": "bot-12345",
        "fromMe": true,
        "senderName": "Alfibot",
        "senderType": "bot",
        "messageType": "text",
        "body": "Las matrículas son del 1-15 de noviembre...",
        "timestamp": 1729625405,
        "read": true,
        "createdAt": "2025-10-22T18:30:05Z"
      },
      {
        "id": 3,
        "messageId": "admin-67890",
        "fromMe": true,
        "senderName": "Admin (admin@ists.edu.ec)",
        "senderType": "admin",
        "messageType": "text",
        "body": "Hola Juan, si tienes más preguntas contáctanos al 099-123-4567",
        "timestamp": 1729625410,
        "adminName": "admin",
        "read": true,
        "createdAt": "2025-10-22T18:30:10Z"
      }
    ]
  }
}
```

### 3. **POST /admin/conversations/:id/send**
Admin sends a message

**Request:**
```json
{
  "message": "Hola, ¿en qué puedo ayudarte?"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "messageId": 123,
    "whatsappMessageId": "admin-msg-uuid",
    "sentAt": "2025-10-22T18:35:00Z"
  }
}
```

### 4. **POST /admin/conversations/:id/mark-read**
Mark conversation as read

**Response:**
```json
{
  "success": true,
  "data": {
    "markedCount": 5
  }
}
```

### 5. **POST /admin/users/:id/block**
Block a user

**Request:**
```json
{
  "blocked": true,
  "reason": "Spam messages"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "userId": 10,
    "blocked": true,
    "blockedAt": "2025-10-22T18:40:00Z"
  }
}
```

### 6. **DELETE /admin/conversations/:id**
Delete (soft) a conversation

**Response:**
```json
{
  "success": true,
  "data": {
    "deleted": true
  }
}
```

### 7. **POST /admin/conversations/:id/temporary**
Enable temporary chat

**Request:**
```json
{
  "temporary": true,
  "hoursUntilExpiry": 24
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "temporary": true,
    "expiresAt": "2025-10-23T18:40:00Z"
  }
}
```

### 8. **GET /admin/conversations/stream** (SSE)
Real-time conversation updates

**Event Types:**
- `new_message`: New message received
- `message_read`: Message marked as read
- `conversation_blocked`: Conversation blocked
- `new_conversation`: New conversation started

**Example Event:**
```
event: new_message
data: {"conversationId": 1, "message": {...}}
```

## Database Schema

### Tables Modified

#### `cht_users`
```sql
usr_blocked boolean DEFAULT false
usr_blocked_at timestamp
usr_blocked_by int (FK to cht_admin_users)
usr_block_reason text
```

#### `cht_conversations`
```sql
cnv_blocked boolean DEFAULT false
cnv_admin_intervened boolean DEFAULT false
cnv_last_admin_message_at timestamp
cnv_temporary boolean DEFAULT false
cnv_expires_at timestamp
cnv_unread_count int DEFAULT 0
```

#### `cht_conversation_messages`
```sql
cvm_sender_type varchar(20) DEFAULT 'user' CHECK IN ('user', 'admin', 'bot')
cvm_admin_id int (FK to cht_admin_users)
cvm_read boolean DEFAULT false
cvm_read_at timestamp
```

## Stored Procedures

1. **fn_get_all_conversations_for_admin** - Get paginated conversations
2. **fn_get_conversation_messages** - Get messages for conversation
3. **sp_block_user** - Block/unblock user
4. **sp_delete_conversation** - Soft delete conversation
5. **sp_send_admin_message** - Admin sends message
6. **sp_mark_messages_as_read** - Mark messages as read
7. **sp_set_conversation_temporary** - Enable/disable temporary chat
8. **fn_cleanup_expired_conversations** - Cleanup expired chats (cron job)

## Implementation Steps

### Phase 1: Basic CRUD (Current)
- [x] Database schema
- [x] Stored procedures
- [ ] Domain types
- [ ] Repository layer
- [ ] Use case layer
- [ ] API endpoints

### Phase 2: Admin Intervention
- [ ] WhatsApp message sending from admin
- [ ] Message tracking
- [ ] Conversation state management

### Phase 3: Real-time Updates
- [ ] SSE implementation
- [ ] Message notifications
- [ ] Unread count updates

### Phase 4: Advanced Features
- [ ] Temporary chat auto-cleanup (cron)
- [ ] Bulk operations (block multiple users)
- [ ] Export conversation history
- [ ] Search conversations

## Frontend Integration

### React/Vue Components Needed

1. **ConversationList** - WhatsApp-like sidebar
2. **ConversationView** - Message thread
3. **MessageInput** - Admin message composer
4. **UserInfo** - User details sidebar
5. **AdminPanel** - Main container

### Example Frontend Flow

```javascript
// 1. Load conversations
GET /admin/conversations?filter=unread&limit=50

// 2. Select conversation
GET /admin/conversations/123/messages

// 3. Mark as read
POST /admin/conversations/123/mark-read

// 4. Send message
POST /admin/conversations/123/send
{ "message": "Hello!" }

// 5. Real-time updates (SSE)
const eventSource = new EventSource('/admin/conversations/stream');
eventSource.addEventListener('new_message', (e) => {
  const data = JSON.parse(e.data);
  // Update UI
});
```

## Security Considerations

1. **JWT Authentication** - All admin endpoints require JWT
2. **Admin Roles** - Check permissions before allowing actions
3. **Rate Limiting** - Prevent admin message spam
4. **Audit Logging** - Log all admin actions
5. **CORS** - Configure for admin frontend domain

## Performance Optimizations

1. **Indexes** - Added on critical columns
2. **Pagination** - Limit conversation lists
3. **Message Limits** - Default 100 messages per load
4. **Caching** - Redis for conversation list (optional)
5. **SSE Backpressure** - Limit concurrent SSE connections

## Testing

### Manual Testing Script
```bash
# 1. Get conversations
curl -X GET http://localhost:8080/admin/conversations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 2. Get messages
curl -X GET http://localhost:8080/admin/conversations/1/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 3. Send message
curl -X POST http://localhost:8080/admin/conversations/1/send \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello from admin!"}'

# 4. Block user
curl -X POST http://localhost:8080/admin/users/10/block \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"blocked": true, "reason": "Spam"}'
```

## Next Steps

1. Create domain types in `domain/admin_conversation.go`
2. Create repository in `repository/admin_conversation_repository.go`
3. Create use case in `usecase/admin_conversation_usecase.go`
4. Create API routes in `api/route/admin_conversation_router.go`
5. Implement SSE for real-time updates
6. Add to main router setup
7. Test end-to-end

## Cron Job Setup

Add to cron for auto-cleanup of temporary conversations:

```bash
# Every hour, cleanup expired temporary conversations
0 * * * * curl -X POST http://localhost:8080/admin/conversations/cleanup
```
