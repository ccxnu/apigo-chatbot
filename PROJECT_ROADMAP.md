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

### **PHASE 1: WhatsApp Core Integration** 🟢 IN PROGRESS

#### 1.1 WhatsApp Client Setup ✅ COMPLETED
- [x] Install whatsmeow dependency (`go get go.mau.fi/whatsmeow`)
- [x] Create WhatsApp client wrapper (`internal/whatsapp/client.go`)
- [x] Implement session management (PostgreSQL storage)
  - [x] Database schema (`cht_whatsapp_sessions` table)
  - [x] Stored procedures (create, update, get session)
  - [x] Repository layer (`repository/whatsapp_session_repository.go`)
  - [x] Use case layer (`usecase/whatsapp_session_usecase.go`)
- [x] Create QR code generation endpoint for admin
  - [x] Service method `GetQRChannel()`
  - [x] API router (`api/route/whatsapp_admin_router.go`)
- [x] Implement connection status monitoring
  - [x] Handle connection events (`Connected`, `Disconnected`)
  - [x] Update database on status changes
- [x] Handle reconnection logic
  - [x] Automatic reconnection handled by whatsmeow
- [x] Add session persistence (survive restarts)
  - [x] Device store integration
  - [x] Session restoration from database

#### 1.2 Message Handler Architecture ✅ COMPLETED
- [x] Create base `MessageHandler` interface (`internal/whatsapp/handler.go`)
- [x] Implement handler registry/dispatcher (`MessageDispatcher`)
- [x] Create filter utilities (commands, user types, etc.)
  - [x] `CanHandle()` method for each handler
- [x] Implement message routing logic
  - [x] Sequential handler checking
  - [x] First-match routing pattern
- [x] Add middleware support (logging, auth, rate limiting)
  - [x] Basic logging in dispatcher
  - [x] Error handling and recovery

#### 1.3 Specific Message Handlers ⚠️ PARTIALLY COMPLETED
- [x] **RAGHandler** - Main Q&A handler using similarity search (`handlers/rag_handler.go`)
  - [x] Similarity search integration
  - [x] OpenAI LLM integration
  - [x] Context building from chunks
  - [ ] Conversation history integration
  - [ ] Citation support
- [x] **CommandHandler** - Handle `/help`, `/horarios`, `/comandos` (`handlers/command_handler.go`)
  - [x] Basic command routing
  - [ ] Implement all specific commands
- [ ] **ScheduleHandler** - Query class schedules
- [ ] **EnrollmentHandler** - Answer enrollment questions
- [ ] **FallbackHandler** - Handle unknown messages
- [ ] **WelcomeHandler** - Onboarding new users

#### 1.4 User Management & Permissions 🔴 NOT STARTED
- [ ] Create user registration on first message
    - For this we will use an institute API. We will require the person ID/document 10 length.
    - And validate against those APIs.
- [ ] Implement role detection (student/professor)
    - This is also part of the validation APIs of the institute, we have two: professor and student.
- [ ] Add permission system for handlers
- [ ] Create user session tracking
- [ ] Implement rate limiting per user
- [ ] Add blacklist/whitelist functionality

---

### **PHASE 2: Enhanced RAG System** ⚠️ PARTIALLY COMPLETED

#### 2.1 Conversation Context 🔴 NOT STARTED
- [ ] Create conversation/session table in DB
- [ ] Store message history per user
- [ ] Implement context window (last N messages)
- [ ] Add conversation embeddings for better context
- [ ] Implement conversation summarization

#### 2.2 LLM Integration ⚠️ PARTIALLY COMPLETED
- [x] Integrate OpenAI for response generation (`internal/llm/openai.go`)
  - [x] Chat completion support
  - [x] Error handling
  - [x] Configurable parameters
- [x] Create prompt templates for different query types
  - [x] RAG-based Q&A template
  - [ ] Command-specific templates
- [x] Implement RAG pipeline (retrieve → generate → respond)
  - [x] Similarity search (`usecase/chunk_usecase.go`)
  - [x] Context building from chunks
  - [x] LLM generation with context
  - [x] Chunk statistics tracking
- [ ] Add citation support (source documents)
- [ ] Implement answer quality checks

#### 2.3 Specialized Query Handlers 🔴 NOT STARTED
- [ ] **Schedule Queries** - Parse date/time from messages
  - [ ] "What class do I have now?" - Get current class based on day/time
  - [ ] "What's my schedule today/this week?" - Show daily/weekly schedule
  - [ ] "Where is my next class?" - Next class location and time
- [ ] **Professor Location Queries** - Real-time professor location
  - [ ] "Where is professor X?" - Current location based on schedule
  - [ ] "When is professor X available?" - Free time slots
  - [ ] "What classroom is professor X in?" - Current classroom
- [ ] **Student Schedule Lookup** - For professors
  - [ ] "Show me students in my current class" - List based on time
  - [ ] "Who are my students in [subject]?" - Students by subject
- [ ] **Enrollment Info** - Structured responses for procedures
- [ ] **Academic Calendar** - Event-based queries
- [ ] **Course Information** - Syllabus, requirements, etc.
- [ ] **Contact Information** - Department contacts, hours

---

### **PHASE 3: Admin Management System** ⚠️ PARTIALLY COMPLETED

#### 3.1 Authentication & Sessions 🔴 NOT STARTED
- [ ] Create admin user table
    - See db/01_create_tables.sql and update/add the necessary fields.
    - Also check the initial_data.sql
- [ ] Implement JWT authentication for admin
- [ ] Create login endpoint
- [ ] Add session management
- [ ] Implement role-based access control (RBAC)

#### 3.2 WhatsApp Session Management ⚠️ PARTIALLY COMPLETED
- [x] **GET** `/admin/whatsapp/qr-code` - Get QR for scanning
  - [x] Router endpoint created
  - [ ] Full implementation and testing
- [x] **GET** `/admin/whatsapp/status` - Check connection status
  - [x] Router endpoint created
  - [ ] Full implementation and testing
- [x] **POST** `/admin/whatsapp/start` - Start WhatsApp session
  - [x] Router endpoint created
  - [ ] Full implementation and testing
- [ ] **POST** `/admin/whatsapp/stop` - Stop session
- [ ] **POST** `/admin/whatsapp/restart` - Restart connection
- [ ] **GET** `/admin/whatsapp/device-info` - Get device info

#### 3.3 Knowledge Base Management ✅ COMPLETED
- [x] **GET** `/api/v1/documents/get-all` - List all documents
- [x] **POST** `/api/v1/documents/add` - Create document
- [x] **PUT** `/api/v1/documents/update` - Update document
- [x] **DELETE** `/api/v1/documents/delete` - Delete document (soft delete)
- [x] **POST** `/api/v1/chunks/add` - Add chunks to documents
- [x] Database procedures for document management
- [ ] **POST** `/admin/documents/import` - Bulk import (PDF, DOCX)

#### 3.4 Parameter Management ✅ COMPLETED
- [x] **POST** `/api/v1/parameters/get-all` - List all parameters
- [x] **POST** `/api/v1/parameters/update` - Update parameter
- [x] **POST** `/api/v1/parameters/add` - Create parameter
- [x] Database procedures for parameter management
- [x] Parameter caching system (`internal/cache/parameter_cache.go`)
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

### **PHASE 5: Database Enhancements** ⚠️ PARTIALLY COMPLETED

#### 5.1 New Tables ⚠️ PARTIALLY COMPLETED
- [x] `cht_whatsapp_sessions` - WhatsApp session data
- [x] `cht_documents` - Document storage
- [x] `cht_chunks` - Document chunks with embeddings
- [x] `cht_chunk_statistics` - RAG quality metrics
- [x] `cht_parameters` - System configuration
- [x] `cht_errors` - Error tracking
- [ ] `cht_conversations` - Conversation tracking
- [ ] `cht_users` - WhatsApp users (students/professors)
- [ ] `cht_api_keys` - External API keys
- [ ] `cht_api_usage` - API usage tracking
- [ ] `cht_scheduled_tasks` - Background jobs
- [ ] `cht_notifications` - Admin notifications

#### 5.2 Stored Procedures ⚠️ PARTIALLY COMPLETED
- [x] `sp_whatsapp_session_create` - Create WhatsApp session
- [x] `sp_whatsapp_session_update_status` - Update connection status
- [x] `sp_whatsapp_session_get` - Get session by name
- [x] `sp_document_*` - Full CRUD for documents
- [x] `sp_chunk_*` - Full CRUD for chunks
- [x] `sp_chunk_statistics_*` - Track chunk usage
- [x] `sp_parameter_*` - Parameter management
- [ ] `sp_create_conversation` - Start new conversation
- [ ] `sp_add_message_to_conversation` - Store message
- [ ] `sp_get_conversation_context` - Get recent messages
- [ ] `sp_track_api_usage` - Log API calls
- [ ] `sp_get_user_statistics` - User activity stats
- [ ] `sp_get_daily_metrics` - Daily analytics

---

### **PHASE 6: Schedule & Location Intelligence** 🆕 HIGH PRIORITY

#### 6.1 Professor Schedule Data Management
- [ ] Create `cht_professors` table
  - [ ] Store professor basic info (nombre, cedula, email, coordinacion, etc.)
  - [ ] Link to WhatsApp users
- [ ] Create `cht_professor_subjects` table (materias)
  - [ ] Store subject info (nombre, inicio, fin, carrera)
  - [ ] Link to professor
- [ ] Create `cht_schedule_slots` table (horarios)
  - [ ] Store day, classroom, start time, end time
  - [ ] Link to subject
- [ ] Create `cht_subject_enrollments` table
  - [ ] Link students to subjects
  - [ ] Store enrollment status

#### 6.2 Schedule Sync Service
- [ ] Create scheduled job to sync professor data from AcademicOK API
  - [ ] Fetch all active professors
  - [ ] Update materias and horarios
  - [ ] Update student enrollments
  - [ ] Run daily (e.g., 2 AM)
- [ ] Cache professor schedules for quick lookup
- [ ] Invalidate cache when data updates

#### 6.3 Schedule Query Handlers
- [ ] **StudentScheduleHandler** - "What class do I have now?"
  - [ ] Parse current day/time
  - [ ] Query student's enrolled subjects
  - [ ] Find matching schedule slot
  - [ ] Return: subject, professor, classroom, time
- [ ] **ProfessorLocationHandler** - "Where is professor X?"
  - [ ] Parse professor name from message
  - [ ] Get current day/time
  - [ ] Find professor's current class
  - [ ] Return: classroom, subject, time remaining
- [ ] **WeeklyScheduleHandler** - "My schedule this week"
  - [ ] Get all student's subjects
  - [ ] Group by day
  - [ ] Format as readable schedule
- [ ] **NextClassHandler** - "Where is my next class?"
  - [ ] Get current time
  - [ ] Find next upcoming class
  - [ ] Return: time, subject, classroom, professor

#### 6.4 Professor Query Features (for professors)
- [ ] **CurrentStudentsHandler** - "Who are my students now?"
  - [ ] Get professor's current subject (based on time)
  - [ ] List enrolled students
- [ ] **SubjectStudentsHandler** - "Students in [subject name]"
  - [ ] Parse subject name
  - [ ] List all enrolled students
- [ ] **AvailabilityHandler** - "When am I available?"
  - [ ] Show free time slots based on schedule

#### 6.5 Database Functions
- [ ] `fn_get_professor_current_location(cedula, timestamp)` - Get professor location
- [ ] `fn_get_student_current_class(cedula, timestamp)` - Get student's current class
- [ ] `fn_get_student_schedule(cedula, day)` - Get daily schedule
- [ ] `fn_get_weekly_schedule(cedula)` - Get weekly schedule
- [ ] `fn_get_subject_students(subject_id)` - Get students in subject
- [ ] `fn_find_professor_by_name(name_pattern)` - Search professors by name

---

### **PHASE 7: Background Services**

#### 7.1 Scheduled Tasks
- [ ] Daily analytics aggregation
- [ ] Chunk statistics recalculation
- [ ] Session cleanup (old/inactive)
- [ ] Report generation
- [ ] Data export/backup
- [x] **Professor schedule sync** (moved to PHASE 6)

#### 7.2 Message Queue (Optional)
- [ ] Implement async message processing
- [ ] Handle high-volume message bursts
- [ ] Retry failed operations
- [ ] Priority queue for different user types

---

### **PHASE 8: Advanced Features**

#### 8.1 Multi-language Support
- [ ] Detect message language
- [ ] Translate queries if needed
- [ ] Respond in user's language

#### 8.2 Media Handling
- [ ] Accept image messages (OCR for text extraction)
- [ ] Accept document messages (PDF parsing)
- [ ] Send images in responses (charts, schedules)

#### 8.3 Proactive Notifications
- [ ] Broadcast messages to groups
- [ ] Send reminders (enrollment deadlines, etc.)
- [ ] Event notifications

#### 8.4 Feedback System
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

### **PHASE 9: Deployment & Operations** ⚠️ PARTIALLY COMPLETED

#### 9.1 Docker Setup ✅ COMPLETED
- [x] Create Dockerfile
- [x] Docker deployment script (deploy.sh)
- [x] Production Docker setup
- [x] PostgreSQL container integration

#### 9.2 Monitoring ⚠️ PARTIALLY COMPLETED
- [x] Health check endpoints
- [ ] Prometheus metrics
- [x] Logging system (`internal/logger/logger.go`)
  - [x] Lumberjack log rotation
  - [x] Structured logging with slog
  - [x] File and console output
- [x] Error tracking and logging
  - [x] Database error table
  - [x] Comprehensive error logging in use cases
- [ ] External error tracking (Sentry/similar)

#### 9.3 Documentation ⚠️ PARTIALLY COMPLETED
- [x] API documentation (OpenAPI3.1)
  - [x] Handled by Huma.rocks
  - [x] Auto-generated docs at `/docs`
- [x] Architecture documentation
  - [x] `docs/CLAUDE.md` - Development guidelines
  - [x] `docs/API_RESPONSE_FORMAT.md` - Response format
  - [x] `docs/API_ENDPOINTS.md` - Endpoint documentation
  - [x] `docs/HUMA_GUIDE.md` - Huma framework guide
  - [x] `docs/LOGGING.md` - Logging system
  - [x] `docs/PARAMETERS.md` - Parameter system
  - [x] `docs/USING_PARAMETERS.md` - Parameter usage
  - [x] `docs/WHATSAPP_COMPARISON.md` - WhatsApp libraries
  - [x] `db/WRITESQL.md` - SQL writing guide
- [ ] Admin guide
- [ ] Deployment guide

---

## 🎯 Current Status Summary

### **What's Been Completed** ✅

**Core Infrastructure:**
- ✅ PostgreSQL database with advanced features (procedures, triggers, error tracking)
- ✅ Complete RAG system with chunk statistics and quality metrics
- ✅ OpenAI integration (embeddings + LLM)
- ✅ Parameter management system with caching
- ✅ Comprehensive logging system with rotation
- ✅ Error tracking and monitoring
- ✅ API structure with Huma framework (OpenAPI 3.1)
- ✅ Docker deployment setup

**WhatsApp Integration (PHASE 1):**
- ✅ whatsmeow client wrapper
- ✅ Session management (database + persistence)
- ✅ QR code generation
- ✅ Connection status monitoring
- ✅ Message handler architecture (dispatcher + routing)
- ✅ RAG handler (Q&A with similarity search)
- ✅ Command handler (basic structure)
- ✅ Admin endpoints (start, stop, status, QR)
- ✅ AcademicOK API integration (student/professor validation)
- ✅ User auto-registration flow

**Knowledge Base (PHASE 3.3):**
- ✅ Full CRUD operations for documents
- ✅ Chunk management with embeddings
- ✅ Automatic statistics tracking
- ✅ Similarity search

**User Management:**
- ✅ Database schema for users, conversations, messages
- ✅ Repository and use case layers
- ✅ AcademicOK API validation (student/professor detection)
- ✅ Auto-registration with role assignment

**Overall Progress:** ~45% of MVP complete

---

## 🎯 Next Priority Tasks

**IMMEDIATE (This Week):**
1. **Test WhatsApp Integration End-to-End**
   - Start service and connect via QR
   - Send test messages
   - Verify RAG handler responses
   - Check database updates

2. **Implement Conversation History (PHASE 2.1)**
   - Create `cht_conversations` and `cht_messages` tables
   - Store incoming/outgoing messages
   - Implement context window for RAG

3. **Complete User Management (PHASE 1.4)**
   - Create `cht_users` table
   - Implement registration on first message
   - Integrate institute API validation
   - Add role detection (student/professor)

**SHORT TERM (Next 2 Weeks):**
4. **Enhance Message Handlers (PHASE 1.3)**
   - Complete command implementations
   - Add FallbackHandler
   - Add WelcomeHandler
   - Implement rate limiting

5. **Admin Authentication (PHASE 3.1)**
   - Create admin user table
   - Implement JWT authentication
   - Add RBAC for admin endpoints

**MEDIUM TERM (Next Month):**
6. **External API (PHASE 4)**
   - Design API key system
   - Implement `/v1/chat/completions` endpoint
   - Add streaming support

7. **Analytics & Reports (PHASE 3.6)**
   - Dashboard metrics
   - Message statistics
   - User activity tracking

---

## 📊 Phase Completion Status

| Phase | Status | Completion | Priority |
|-------|--------|-----------|----------|
| **PHASE 1: WhatsApp Integration** | 🟢 In Progress | 75% | HIGH - MVP |
| **PHASE 2: Enhanced RAG System** | ⚠️ Partial | 40% | HIGH - MVP |
| **PHASE 3: Admin Management** | ⚠️ Partial | 35% | HIGH - MVP |
| **PHASE 4: External API** | 🔴 Not Started | 0% | MEDIUM |
| **PHASE 5: Database Enhancements** | ⚠️ Partial | 55% | MEDIUM |
| **PHASE 6: Schedule & Location Intelligence** 🆕 | 🔴 Not Started | 0% | **HIGH** |
| **PHASE 7: Background Services** | 🔴 Not Started | 0% | MEDIUM |
| **PHASE 8: Advanced Features** | 🔴 Not Started | 0% | LOW |
| **PHASE 9: Testing & Quality** | 🔴 Not Started | 0% | LOW |
| **PHASE 10: Deployment & Operations** | ⚠️ Partial | 60% | HIGH - MVP |

**Legend:**
- ✅ Completed
- 🟢 In Progress
- ⚠️ Partially Completed
- 🔴 Not Started

---

## 🏗️ Technical Stack (Implemented)

**Backend:**
- ✅ Go 1.21+
- ✅ Huma framework (OpenAPI 3.1)
- ✅ PostgreSQL 14+
- ✅ whatsmeow (WhatsApp client)

**AI/ML:**
- ✅ OpenAI API (text-embedding-3-small, gpt-4o-mini)
- ✅ pgvector for similarity search
- ✅ RAG pipeline with context building

**Infrastructure:**
- ✅ Docker + Docker Compose
- ✅ Structured logging (slog + lumberjack)
- ✅ Parameter caching system
- ✅ Database procedures & triggers

---

## Priority Order

**HIGH (Must have for MVP):**
- Phase 1: WhatsApp Integration (70% complete)
- Phase 2: Enhanced RAG with LLM (40% complete)
- Phase 3: Basic Admin (35% complete)
- Phase 9: Deployment & Operations (60% complete)

**MEDIUM (Important for v1.0):**
- Phase 4: External API
- Phase 5: Database Enhancements
- Phase 6: Background Services

**LOW (Nice to have for future):**
- Phase 7: Advanced Features
- Phase 8: Comprehensive Testing
