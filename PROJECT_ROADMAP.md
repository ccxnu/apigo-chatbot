# Instituto Chatbot - Complete Backend Roadmap

## Project Overview
**Purpose:** WhatsApp chatbot for educational institute to handle student/professor questions, schedules, and RAG-based information retrieval.

**Users:**
- **Students** - Limited access via WhatsApp only
- **Professors** - Open access via WhatsApp + API
- **Admin** - Full system management, QR code scanning, analytics

**Modules:**
1. WhatsApp Integration (Primary Interface)
2. RAG System (Already implemented ✅)
3. API for External Integration (Claude-like API)
4. Admin Backend (Management & Analytics)

---

## 📋 COMPLETE TODO LIST

### **PHASE 1: WhatsApp Core Integration** 🟢 CURRENT PRIORITY

#### 1.1 WhatsApp Client Setup
- [ ] Install whatsmeow dependency (`go get go.mau.fi/whatsmeow`)
- [ ] Create WhatsApp client wrapper (`internal/whatsapp/client.go`)
- [ ] Implement session management (SQLite/PostgreSQL storage)
- [ ] Create QR code generation endpoint for admin
- [ ] Implement connection status monitoring
- [ ] Handle reconnection logic
- [ ] Add session persistence (survive restarts)

#### 1.2 Message Handler Architecture
- [ ] Create base `MessageHandler` interface
- [ ] Implement handler registry/dispatcher
- [ ] Create filter utilities (commands, user types, etc.)
- [ ] Implement message routing logic
- [ ] Add middleware support (logging, auth, rate limiting)

#### 1.3 Specific Message Handlers
- [ ] **RAGHandler** - Main Q&A handler using similarity search
- [ ] **CommandHandler** - Handle `/help`, `/horarios`, `/comandos`
- [ ] **ScheduleHandler** - Query class schedules
- [ ] **EnrollmentHandler** - Answer enrollment questions
- [ ] **FallbackHandler** - Handle unknown messages
- [ ] **WelcomeHandler** - Onboarding new users

#### 1.4 User Management & Permissions
- [ ] Create user registration on first message
    - For this we will use an institute API. We will require the person ID/document 10 leght.
    - And validate against those apis.
- [ ] Implement role detection (student/professor)
    - This is also parte of the validation apis of the institute, we have two professor and student.
- [ ] Add permission system for handlers
- [ ] Create user session tracking
- [ ] Implement rate limiting per user
- [ ] Add blacklist/whitelist functionality

---

### **PHASE 2: Enhanced RAG System**

#### 2.1 Conversation Context
- [ ] Create conversation/session table in DB
- [ ] Store message history per user
- [ ] Implement context window (last N messages)
- [ ] Add conversation embeddings for better context
- [ ] Implement conversation summarization

#### 2.2 LLM Integration
- [ ] Integrate Grok/OpenAI for response generation
- [ ] Create prompt templates for different query types
- [ ] Implement RAG pipeline (retrieve → generate → respond)
- [ ] Add citation support (source documents)
- [ ] Implement answer quality checks

#### 2.3 Specialized Query Handlers
- [ ] **Schedule Queries** - Parse date/time from messages
- [ ] **Enrollment Info** - Structured responses for procedures
- [ ] **Academic Calendar** - Event-based queries
- [ ] **Course Information** - Syllabus, requirements, etc.
- [ ] **Contact Information** - Department contacts, hours

---

### **PHASE 3: Admin Management System**

#### 3.1 Authentication & Sessions
- [ ] Create admin user table
    - See db/01_create_tables.sql and update/add the necesary fields.
    - Also check the initial_data.sql
- [ ] Implement JWT authentication for admin
- [ ] Create login endpoint
- [ ] Add session management
- [ ] Implement role-based access control (RBAC)

#### 3.2 WhatsApp Session Management
- [ ] **GET** `/admin/whatsapp/qr-code` - Get QR for scanning
- [ ] **GET** `/admin/whatsapp/status` - Check connection status
- [ ] **POST** `/admin/whatsapp/start` - Start WhatsApp session
- [ ] **POST** `/admin/whatsapp/stop` - Stop session
- [ ] **POST** `/admin/whatsapp/restart` - Restart connection
- [ ] **GET** `/admin/whatsapp/device-info` - Get device info

#### 3.3 Knowledge Base Management
- [ ] **GET** `/admin/documents` - List all documents
- [ ] **POST** `/admin/documents` - Create document
- [ ] **PUT** `/admin/documents/:id` - Update document
- [ ] **DELETE** `/admin/documents/:id` - Delete document
- [ ] **POST** `/admin/documents/:id/chunks` - Regenerate chunks
- [ ] **POST** `/admin/documents/import` - Bulk import (PDF, DOCX)

#### 3.4 Parameter Management
- [ ] **GET** `/admin/parameters` - List all parameters
- [ ] **PUT** `/admin/parameters/:code` - Update parameter
- [ ] **POST** `/admin/parameters` - Create parameter
- [ ] Hot-reload parameters (refresh cache without restart)

#### 3.5 User Management
- [ ] **GET** `/admin/users` - List all users (students/professors)
- [ ] **GET** `/admin/users/:id` - Get user details
- [ ] **PUT** `/admin/users/:id/role` - Change user role
- [ ] **POST** `/admin/users/:id/block` - Block/unblock user
- [ ] **GET** `/admin/users/:id/conversations` - View user chat history

#### 3.6 Analytics & Reports
- [ ] **GET** `/admin/analytics/overview` - Dashboard metrics
- [ ] **GET** `/admin/analytics/messages` - Message statistics
- [ ] **GET** `/admin/analytics/top-queries` - Most asked questions
- [ ] **GET** `/admin/analytics/chunk-usage` - Most used knowledge chunks
- [ ] **GET** `/admin/analytics/response-quality` - Quality metrics
- [ ] **GET** `/admin/analytics/user-activity` - User engagement
- [ ] **GET** `/admin/reports/daily` - Daily report
- [ ] **GET** `/admin/reports/weekly` - Weekly summary
- [ ] **POST** `/admin/reports/export` - Export data (CSV/JSON)

#### 3.7 Conversation Monitoring
- [ ] **GET** `/admin/conversations` - List all conversations
- [ ] **GET** `/admin/conversations/:id` - View full conversation
- [ ] **GET** `/admin/conversations/active` - Active chats now
- [ ] **POST** `/admin/conversations/:id/intervene` - Admin can send message
- [ ] **GET** `/admin/conversations/flagged` - Flagged/problematic chats

---

### **PHASE 4: External API (Claude-like Interface)**

#### 4.1 API Design
- [ ] Design API key system for external clients
- [ ] Create rate limiting per API key
- [ ] Implement usage quotas
- [ ] Add request/response logging

#### 4.2 Chat Completions API
- [ ] **POST** `/v1/chat/completions` - Claude-style chat API
- [ ] Support streaming responses (SSE)
- [ ] Include RAG context in responses
- [ ] Add model selection (different RAG strategies)
- [ ] Implement token counting/billing

#### 4.3 Embeddings API
- [ ] **POST** `/v1/embeddings` - Generate embeddings
- [ ] Batch embedding support
- [ ] Multiple model support

#### 4.4 Search API
- [ ] **POST** `/v1/search` - Direct knowledge base search
- [ ] Support filters (category, date range)
- [ ] Return sources with citations

---

### **PHASE 5: Database Enhancements**

#### 5.1 New Tables
- [ ] `cht_whatsapp_sessions` - WhatsApp session data
- [ ] `cht_conversations` - Conversation tracking
- [ ] `cht_api_keys` - External API keys
- [ ] `cht_api_usage` - API usage tracking
- [ ] `cht_scheduled_tasks` - Background jobs
- [ ] `cht_notifications` - Admin notifications

#### 5.2 Stored Procedures
- [ ] `sp_create_conversation` - Start new conversation
- [ ] `sp_add_message_to_conversation` - Store message
- [ ] `sp_get_conversation_context` - Get recent messages
- [ ] `sp_track_api_usage` - Log API calls
- [ ] `sp_get_user_statistics` - User activity stats
- [ ] `sp_get_daily_metrics` - Daily analytics

---

### **PHASE 6: Background Services**

#### 6.1 Scheduled Tasks
- [ ] Daily analytics aggregation
- [ ] Chunk statistics recalculation
- [ ] Session cleanup (old/inactive)
- [ ] Report generation
- [ ] Data export/backup

#### 6.2 Message Queue (Optional)
- [ ] Implement async message processing
- [ ] Handle high-volume message bursts
- [ ] Retry failed operations
- [ ] Priority queue for different user types

---

### **PHASE 7: Advanced Features**

#### 7.1 Multi-language Support
- [ ] Detect message language
- [ ] Translate queries if needed
- [ ] Respond in user's language

#### 7.2 Media Handling
- [ ] Accept image messages (OCR for text extraction)
- [ ] Accept document messages (PDF parsing)
- [ ] Send images in responses (charts, schedules)

#### 7.3 Proactive Notifications
- [ ] Broadcast messages to groups
- [ ] Send reminders (enrollment deadlines, etc.)
- [ ] Event notifications

#### 7.4 Feedback System
- [ ] Thumbs up/down for answers
- [ ] Collect user feedback
- [ ] Use feedback to improve quality metrics

---

### **PHASE 8: Testing & Quality**

#### 8.1 Unit Tests
- [ ] WhatsApp client tests
- [ ] Handler tests
- [ ] RAG pipeline tests
- [ ] API endpoint tests

#### 8.2 Integration Tests
- [ ] End-to-end message flow
- [ ] Database operations
- [ ] External API calls

#### 8.3 Performance Testing
- [ ] Load testing (concurrent messages)
- [ ] RAG query performance
- [ ] Database query optimization

---

### **PHASE 9: Deployment & Operations**

#### 9.1 Docker Setup
- [ ] Create Dockerfile
    -- I already have a Dockerfile
- [ ] Docker Compose for local dev
    - I use deploy.sh for this. Never docker-compose.
- [ ] Production Docker setup

#### 9.2 Monitoring
- [ ] Health check endpoints
- [ ] Prometheus metrics
- [ ] Logging aggregation
- [ ] Error tracking (Sentry/similar)

#### 9.3 Documentation
- [ ] API documentation (OpenAPI3.1)
    - Handle by Huma.rocks
- [ ] Admin guide
- [ ] Deployment guide
- [ ] Architecture documentation

---

## 🎯 Current Focus: PHASE 1 - WhatsApp Integration

**Next Steps:**
1. Install whatsmeow
2. Create WhatsApp client wrapper
3. Implement basic message handler architecture
4. Create RAG handler for Q&A
5. Add QR code endpoint for admin

**Dependencies Already Complete:** ✅
- PostgreSQL database
- RAG system with chunk statistics
- Embedding service (OpenAI)
- Parameter management
- Basic API structure

---

## Priority Order

**HIGH (Must have for MVP):**
- Phase 1: WhatsApp Integration
- Phase 2: Enhanced RAG with LLM
- Phase 3: Basic Admin (QR, Knowledge, Analytics)

**MEDIUM (Important):**
- Phase 4: External API
- Phase 5: Database Enhancements
- Phase 6: Background Services

**LOW (Nice to have):**
- Phase 7: Advanced Features
- Phase 8: Comprehensive Testing
- Phase 9: Production Hardening
