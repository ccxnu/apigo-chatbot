--- =====================================================
-- Initial Parameters Data
-- =====================================================

do $$
begin
    -- =====================================================
    -- System ERROR_CODES
    -- =====================================================

    -- OK
    if not exists (select 1 from cht_parameters where prm_code = 'OK') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'OK', '{"message": "Operación exitosa"}'::jsonb, 'Success response code');
    end if;

    -- ERR_INTERNAL_DB
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_INTERNAL_DB') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_INTERNAL_DB', '{"message": "Ha ocurrido un error en la base de datos"}'::jsonb, 'Database error');
    end if;

    -- ERR_HTTP_SERVICE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_SERVICE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_SERVICE', '{"message": "Error en la comunicación HTTP"}'::jsonb, 'HTTP service error');
    end if;

    -- ERR_INTERNAL_SERVER (HTTP 500)
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_INTERNAL_SERVER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_INTERNAL_SERVER', '{"message": "Error interno del servidor"}'::jsonb, 'Internal server error (500)');
    end if;

    -- ERR_PARAM_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_PARAM_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_PARAM_NOT_FOUND', '{"message": "Parámetro no encontrado"}'::jsonb, 'Parameter not found');
    end if;

    -- ERR_PARAM_CODE_EXISTS
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_PARAM_CODE_EXISTS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_PARAM_CODE_EXISTS', '{"message": "El código del parámetro ya existe"}'::jsonb, 'Parameter code already exists');
    end if;

    -- ERR_CREATE_PARAMETER
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_PARAMETER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_PARAMETER', '{"message": "Error al crear el parámetro"}'::jsonb, 'Error creating parameter');
    end if;

    -- ERR_UPDATE_PARAMETER
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_PARAMETER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_PARAMETER', '{"message": "Error al actualizar el parámetro"}'::jsonb, 'Error updating parameter');
    end if;

    -- ERR_DELETE_PARAMETER
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_DELETE_PARAMETER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_DELETE_PARAMETER', '{"message": "Error al eliminar el parámetro"}'::jsonb, 'Error deleting parameter');
    end if;

    -- ERR_USER_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_USER_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_USER_NOT_FOUND', '{"message": "Usuario no encontrado"}'::jsonb, 'User not found');
    end if;

    -- ERR_SESSION_EXPIRED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_SESSION_EXPIRED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_SESSION_EXPIRED', '{"message": "La sesión ha expirado"}'::jsonb, 'Session expired');
    end if;

    -- ERR_INVALID_TOKEN
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_INVALID_TOKEN') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_INVALID_TOKEN', '{"message": "Token de autenticación inválido"}'::jsonb, 'Invalid token');
    end if;

    -- ERR_DOCUMENT_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_DOCUMENT_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_DOCUMENT_NOT_FOUND', '{"message": "Documento no encontrado"}'::jsonb, 'Document not found');
    end if;

    -- ERR_CHUNK_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CHUNK_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CHUNK_NOT_FOUND', '{"message": "Fragmento de documento no encontrado"}'::jsonb, 'Chunk not found');
    end if;

    -- ERR_CHUNK_STATS_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CHUNK_STATS_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CHUNK_STATS_NOT_FOUND', '{"message": "Estadísticas del fragmento no encontradas"}'::jsonb, 'Chunk statistics not found');
    end if;

    -- ERR_CREATE_CHUNK
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_CHUNK') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_CHUNK', '{"message": "Error al crear el fragmento de documento"}'::jsonb, 'Error creating chunk');
    end if;

    -- ERR_UPDATE_CHUNK_EMBEDDING
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_CHUNK_EMBEDDING') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_CHUNK_EMBEDDING', '{"message": "Error al actualizar el embedding del fragmento"}'::jsonb, 'Error updating chunk embedding');
    end if;

    -- ERR_DELETE_CHUNK
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_DELETE_CHUNK') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_DELETE_CHUNK', '{"message": "Error al eliminar el fragmento"}'::jsonb, 'Error deleting chunk');
    end if;

    -- ERR_BULK_CREATE_CHUNKS
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_BULK_CREATE_CHUNKS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_BULK_CREATE_CHUNKS', '{"message": "Error al crear múltiples fragmentos"}'::jsonb, 'Error bulk creating chunks');
    end if;

    -- ERR_INCREMENT_CHUNK_USAGE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_INCREMENT_CHUNK_USAGE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_INCREMENT_CHUNK_USAGE', '{"message": "Error al incrementar el contador de uso del fragmento"}'::jsonb, 'Error incrementing chunk usage');
    end if;

    -- ERR_UPDATE_CHUNK_METRICS
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_CHUNK_METRICS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_CHUNK_METRICS', '{"message": "Error al actualizar las métricas de calidad del fragmento"}'::jsonb, 'Error updating chunk quality metrics');
    end if;

    -- ERR_UPDATE_CHUNK_STALENESS
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_CHUNK_STALENESS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_CHUNK_STALENESS', '{"message": "Error al actualizar la obsolescencia del fragmento"}'::jsonb, 'Error updating chunk staleness');
    end if;

    -- ERR_CREATE_DOCUMENT
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_DOCUMENT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_DOCUMENT', '{"message": "Error al crear el documento"}'::jsonb, 'Error creating document');
    end if;

    -- ERR_UPDATE_DOCUMENT
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_DOCUMENT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_DOCUMENT', '{"message": "Error al actualizar el documento"}'::jsonb, 'Error updating document');
    end if;

    -- ERR_DELETE_DOCUMENT
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_DELETE_DOCUMENT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_DELETE_DOCUMENT', '{"message": "Error al eliminar el documento"}'::jsonb, 'Error deleting document');
    end if;

    -- ERR_WHATSAPP_SESSION_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_WHATSAPP_SESSION_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_WHATSAPP_SESSION_NOT_FOUND', '{"message": "Sesión de WhatsApp no encontrada"}'::jsonb, 'WhatsApp session not found');
    end if;

    -- ERR_IDENTITY_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_IDENTITY_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_IDENTITY_NOT_FOUND', '{"message": "Número de cédula no encontrado en el sistema del instituto"}'::jsonb, 'Identity number not found in institute system');
    end if;

    -- ERR_INVALID_IDENTITY
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_INVALID_IDENTITY') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_INVALID_IDENTITY', '{"message": "Número de cédula inválido o no encontrado"}'::jsonb, 'Invalid or not found identity number');
    end if;

    -- ERR_USER_ALREADY_EXISTS
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_USER_ALREADY_EXISTS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_USER_ALREADY_EXISTS', '{"message": "El usuario ya existe en el sistema"}'::jsonb, 'User already exists');
    end if;

    -- ERR_CREATE_USER
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_USER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_USER', '{"message": "Error al crear el usuario"}'::jsonb, 'Error creating user');
    end if;

    -- ERR_EXTERNAL_USER_INFO_REQUIRED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_EXTERNAL_USER_INFO_REQUIRED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_EXTERNAL_USER_INFO_REQUIRED', '{"message": "Usuario externo - por favor proporciona tu nombre y correo electrónico"}'::jsonb, 'External user information required');
    end if;

    -- ERR_UPDATE_WHATSAPP_SESSION
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_WHATSAPP_SESSION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_WHATSAPP_SESSION', '{"message": "Error al actualizar sesión de WhatsApp"}'::jsonb, 'Error updating WhatsApp session');
    end if;

    -- ERR_UPDATE_QR_CODE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UPDATE_QR_CODE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UPDATE_QR_CODE', '{"message": "Error al actualizar código QR"}'::jsonb, 'Error updating QR code');
    end if;

    -- ERR_CONVERSATION_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CONVERSATION_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CONVERSATION_NOT_FOUND', '{"message": "Conversación no encontrada"}'::jsonb, 'Conversation not found');
    end if;

    -- ERR_CREATE_CONVERSATION
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_CONVERSATION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_CONVERSATION', '{"message": "Error al crear conversación"}'::jsonb, 'Error creating conversation');
    end if;

    -- ERR_LINK_USER_CONVERSATION
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_LINK_USER_CONVERSATION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_LINK_USER_CONVERSATION', '{"message": "Error al vincular usuario a conversación"}'::jsonb, 'Error linking user to conversation');
    end if;

    -- ERR_CREATE_MESSAGE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_CREATE_MESSAGE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_CREATE_MESSAGE', '{"message": "Error al crear mensaje"}'::jsonb, 'Error creating message');
    end if;

    -- ERR_USER_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_USER_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_USER_NOT_FOUND', '{"message": "Usuario no encontrado"}'::jsonb, 'User not found');
    end if;

    -- ERR_PERMISSION_DENIED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_PERMISSION_DENIED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_PERMISSION_DENIED', '{"message": "Permiso denegado"}'::jsonb, 'Permission denied');
    end if;

    -- ERR_BAD_REQUEST
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_BAD_REQUEST') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_BAD_REQUEST', '{"message": "Solicitud incorrecta. Verifique los datos enviados"}'::jsonb, 'Bad request error (400)');
    end if;

    -- ERR_UNAUTHORIZED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UNAUTHORIZED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UNAUTHORIZED', '{"message": "No autorizado. Se requiere autenticación"}'::jsonb, 'Unauthorized error (401)');
    end if;

    -- ERR_FORBIDDEN
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_FORBIDDEN') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_FORBIDDEN', '{"message": "Acceso prohibido. No tiene permisos suficientes"}'::jsonb, 'Forbidden error (403)');
    end if;

    -- ERR_NOT_FOUND
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_NOT_FOUND') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_NOT_FOUND', '{"message": "Recurso no encontrado"}'::jsonb, 'Not found error (404)');
    end if;

    -- ERR_VALIDATION_FAILED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_VALIDATION_FAILED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_VALIDATION_FAILED', '{"message": "Error de validación. Revise los campos enviados"}'::jsonb, 'Validation failed error (422)');
    end if;

    -- ERR_SERVICE_UNAVAILABLE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_SERVICE_UNAVAILABLE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_SERVICE_UNAVAILABLE', '{"message": "Servicio no disponible temporalmente"}'::jsonb, 'Service unavailable error (503)');
    end if;

    -- ERR_UNKNOWN
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_UNKNOWN') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_UNKNOWN', '{"message": "Error desconocido"}'::jsonb, 'Unknown error');
    end if;

    -- ERR_EMBEDDING_GENERATION
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_EMBEDDING_GENERATION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_EMBEDDING_GENERATION', '{"message": "Error al generar embedding del texto"}'::jsonb, 'Embedding generation error');
    end if;

    -- ERR_HTTP_REQUEST_MARSHAL
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_REQUEST_MARSHAL') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_REQUEST_MARSHAL', '{"message": "Error al serializar datos de la petición HTTP"}'::jsonb, 'HTTP request marshal error');
    end if;

    -- ERR_HTTP_REQUEST_CREATE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_REQUEST_CREATE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_REQUEST_CREATE', '{"message": "Error al crear petición HTTP"}'::jsonb, 'HTTP request creation error');
    end if;

    -- ERR_HTTP_REQUEST_FAILED
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_REQUEST_FAILED') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_REQUEST_FAILED', '{"message": "Error al ejecutar petición HTTP"}'::jsonb, 'HTTP request execution failed');
    end if;

    -- ERR_HTTP_RESPONSE_ERROR
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_RESPONSE_ERROR') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_RESPONSE_ERROR', '{"message": "Error en respuesta del servicio HTTP"}'::jsonb, 'HTTP response error status');
    end if;

    -- ERR_HTTP_RESPONSE_DECODE
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_RESPONSE_DECODE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_RESPONSE_DECODE', '{"message": "Error al decodificar respuesta HTTP"}'::jsonb, 'HTTP response decode error');
    end if;

    -- ERR_HTTP_RESPONSE_READ
    if not exists (select 1 from cht_parameters where prm_code = 'ERR_HTTP_RESPONSE_READ') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('ERROR_CODES', 'ERR_HTTP_RESPONSE_READ', '{"message": "Error al leer respuesta HTTP"}'::jsonb, 'HTTP response read error');
    end if;

    -- =====================================================
    -- Application Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'APP_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'APP_CONFIGURATION',
            'APP_CONFIG',
            '{
                "application": "wsChatbot",
                "appEnv": "development",
                "contextTimeout": 2,
                "basicAuth": "X-Auth wiaAchcHks3rBxIhJQem1nLoMDwdoQ=="
            }'::jsonb,
            'General application configuration'
        );
    end if;

    -- =====================================================
    -- Logging Configuration
    -- Note: output is overridden by APP_CONFIG.appEnv:
    --   - development: stdout only (easier for local debugging)
    --   - production: file only (persistent logs, less container noise)
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'LOG_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'LOG_CONFIGURATION',
            'LOG_CONFIG',
            '{
                "level": "info",
                "format": "json",
                "output": "both",
                "filePath": "logs/app.log",
                "maxSizeMB": 100,
                "maxBackups": 5,
                "maxAgeDays": 30
            }'::jsonb,
            'Logging system configuration - slog with file rotation'
        );
    end if;

    -- =====================================================
    -- JWT Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'JWT_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'JWT_CONFIGURATION',
            'JWT_CONFIG',
            '{
                "accessTokenSecret": "Y2xhdmVzdXBlcnNlY3JldGFAMTc2NjQ1NDU2Nw==",
                "accessTokenExpiryHour": 2,
                "refreshTokenSecret": "ttahoeuthaorcuhntuhoatneuh",
                "refreshTokenExpiryHour": 24
            }'::jsonb,
            'JWT token configuration and secrets'
        );
    end if;

    -- =====================================================
    -- Email Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'EMAIL_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'EMAIL_CONFIGURATION',
            'EMAIL_CONFIG',
            '{
                "sender": "noreply@example.com"
            }'::jsonb,
            'Email sender configuration'
        );
    end if;

    -- =====================================================
    -- WhatsApp Connect Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'WPP_CONNECT_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'WPP_CONFIGURATION',
            'WPP_CONNECT_CONFIG',
            '{
                "baseUrl": "http://0.0.0.0:21465"
            }'::jsonb,
            'WhatsApp Connect service base URL'
        );
    end if;

    -- =====================================================
    -- Embedding Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'EMBEDDING_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'EMBEDDING_CONFIGURATION',
            'EMBEDDING_CONFIG',
            '{
                "openaiUrl": "https://api.openai.com/v1/embeddings",
                "openaiApiKey": "sk-proj-f",
                "openaiModel": "text-embedding-3-small",
                "ollamaUrl": "https://localhost:11434/embeddings/:model",
                "ollamaModel": "nomic"
            }'::jsonb,
            'Embedding models configuration (OpenAI and Ollama)'
        );
    end if;

    -- =====================================================
    -- DEPRECATED: Old LLM configurations (replaced by LLM_CONFIG below)
    -- These are kept for reference but should not be used
    -- =====================================================
    -- GROK_API_KEY and CLAUDE_API_KEY have been replaced by unified LLM_CONFIG

    if not exists (select 1 from cht_parameters where prm_code = 'RATE_LIMIT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ( 'LLM_CONFIGURATION', 'RATE_LIMIT', '{"requestsPerMinute": 60, "requestsPerHour": 1000}'::jsonb, 'API rate limiting configuration' );
    end if;

    -- =====================================================
    -- USER_ROLEs
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'ROLE_ADMIN') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('USER_ROLE', 'ROLE_ADMIN', '{}'::jsonb, 'Full system access');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'ROLE_USER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('USER_ROLE', 'ROLE_USER', '{}'::jsonb, 'Standard user access');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'ROLE_STUDENT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('USER_ROLE', 'ROLE_STUDENT', '{}'::jsonb, 'Standard student user access');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'ROLE_GUEST') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('USER_ROLE', 'ROLE_GUEST', '{}'::jsonb, 'Limited guest access');
    end if;

    -- =====================================================
    -- Message Roles
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'MSG_ROLE_USER') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('MESSAGE_ROLE', 'MSG_ROLE_USER', '{}'::jsonb, 'Message from user');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'MSG_ROLE_ASSISTANT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('MESSAGE_ROLE', 'MSG_ROLE_ASSISTANT', '{}'::jsonb, 'Message from AI assistant');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'MSG_ROLE_SYSTEM') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('MESSAGE_ROLE', 'MSG_ROLE_SYSTEM', '{}'::jsonb, 'System message');
    end if;

    -- =====================================================
    -- Session Origins
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'ORIGIN_WEB') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('SESSION_ORIGIN', 'ORIGIN_WEB', '{}'::jsonb, 'Session from web application');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'ORIGIN_MOBILE') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('SESSION_ORIGIN', 'ORIGIN_MOBILE', '{}'::jsonb, 'Session from mobile app');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'ORIGIN_API') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('SESSION_ORIGIN', 'ORIGIN_API', '{}'::jsonb, 'Session from API');
    end if;

    -- =====================================================
    -- Document Categories
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'DOC_CAT_GENERAL') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('DOCUMENT_CATEGORY', 'DOC_CAT_GENERAL', '{}'::jsonb, 'General knowledge documents');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'DOC_CAT_TECHNICAL') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('DOCUMENT_CATEGORY', 'DOC_CAT_TECHNICAL', '{}'::jsonb, 'Technical documentation');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'DOC_CAT_FAQ') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('DOCUMENT_CATEGORY', 'DOC_CAT_FAQ', '{}'::jsonb, 'Frequently asked questions');
    end if;

    -- =====================================================
    -- Functionalities
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'FUNC_CHAT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('FUNCTIONALITY', 'FUNC_CHAT', '{}'::jsonb, 'Access to chat interface');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'FUNC_DOCUMENTS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('FUNCTIONALITY', 'FUNC_DOCUMENTS', '{}'::jsonb, 'Manage documents');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'FUNC_USERS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('FUNCTIONALITY', 'FUNC_USERS', '{}'::jsonb, 'Manage users');
    end if;

    if not exists (select 1 from cht_parameters where prm_code = 'FUNC_PARAMETERS') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ('FUNCTIONALITY', 'FUNC_PARAMETERS', '{}'::jsonb, 'Manage system parameters');
    end if;

    -- =====================================================
    -- System Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'SESSION_TIMEOUT') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values ( 'SESSION_CONFIGURATION', 'SESSION_TIMEOUT', '{"minutes": 30}'::jsonb, 'User session timeout in minutes' );
    end if;

    -- =====================================================
    -- WhatsApp Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'WHATSAPP_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'WHATSAPP_CONFIGURATION',
            'WHATSAPP_CONFIG',
            '{
                "enabled": false,
                "sessionName": "chatbot-session",
                "contextTimeout": 2,
                "autoReconnect": true,
                "maxReconnectAttempts": 5
            }'::jsonb,
            'WhatsApp service configuration - set enabled:true to activate'
        );
    end if;

    -- =====================================================
    -- AcademicOK API Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'ACADEMICOK_CONFIGURATION') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'ACADEMICOK_CONFIGURATION',
            'ACADEMICOK_CONFIGURATION',
            '{
                "personaURL": "https://itsl.academicok.com/api?a=apidatospersona",
                "docenteURL": "https://itsl.academicok.com/api?a=apidatosdocente",
                "personaKey": "j62kDJnltU4wVqp",
                "docenteKey": "123",
                "timeout": 10
            }'::jsonb,
            'AcademicOK API configuration for student and professor validation'
        );
    end if;

    -- =====================================================
    -- RAG Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'RAG_SEARCH_LIMIT') then
           insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
           values ('RAG_CONFIGURATION', 'RAG_SEARCH_LIMIT', '{"value": 5}'::jsonb, 'Maximum number of chunks to retrieve in similarity search');
	end if;

      if not exists (select 1 from cht_parameters where prm_code = 'RAG_MIN_SIMILARITY') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_MIN_SIMILARITY', '{"value": 0.3}'::jsonb, 'Minimum similarity threshold (0.0-1.0) for RAG search results');
      end if;

      -- =====================================================
      -- RAG Configuration - Error Messages
      -- =====================================================
      if not exists (select 1 from cht_parameters where prm_code = 'RAG_ERROR_MESSAGE') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_ERROR_MESSAGE', '{"message": "Lo siento, ocurrió un error al procesar tu mensaje."}'::jsonb, 'Error message when RAG processing fails');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'RAG_NO_RESULTS_MESSAGE') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_NO_RESULTS_MESSAGE', '{"message": "Lo siento, no encontré información relevante sobre tu consulta. ¿Podrías reformular tu pregunta?"}'::jsonb, 'Message when no relevant chunks found');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'RAG_NO_RELEVANT_INFO') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_NO_RELEVANT_INFO', '{"message": "No encontré información relevante."}'::jsonb, 'Simple no results message');
      end if;

      -- =====================================================
      -- RAG Configuration - Response Templates
      -- =====================================================
      if not exists (select 1 from cht_parameters where prm_code = 'RAG_SOURCE_PREFIX') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_SOURCE_PREFIX', '{"template": "## Fuente %d:%s\n"}'::jsonb, 'Template for source prefix in context building (args: index, title)');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'RAG_ANSWER_FORMAT')
 then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_ANSWER_FORMAT', '{"template": "📚 Basado en: *%s*\n\n%s\n\n_Similitud: %.0f%%_"}'::jsonb, 'Template for answer formatting (args: docTitle, content, similarity%)');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'RAG_MULTIPLE_SOURCES_SUFFIX') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'RAG_MULTIPLE_SOURCES_SUFFIX', '{"template":"\n\n_(También encontré información en %d documentos más)_"}'::jsonb, 'Suffix when multiple sources found (args: count)');
      end if;

      -- =====================================================
      -- WhatsApp Bot - Command Messages
      -- =====================================================
      if not exists (select 1 from cht_parameters where prm_code = 'MESSAGE_HELP') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_HELP', '{"message": "👋 *Bienvenido al Asistente del Instituto*\n\nSoy tu asistente virtual y puedo ayudarte con:\n\n🎓 *Información Académica*\n   • Programas y carreras\n   • Requisitos de admisión\n   • Proceso de matrícula\n   • Calendario académico\n\n📚 *Consultas Generales*\n   Solo escribe tu pregunta y te ayudaré a encontrar la información que necesitas.\n\n⚡ *Comandos Disponibles*\n   /help - Muestra esta ayuda\n   /horarios - Consulta horarios de clases\n   /comandos - Lista todos los comandos\n\n💬 También puedes hacer preguntas directamente, por ejemplo:\n   \"¿Cuál es el proceso de matrícula?\"\n \"¿Qué carreras ofrecen?\"\n\n¿En qué puedo ayudarte hoy?"}'::jsonb, 'Help command response text');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'MESSAGE_SCHEDULES') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_SCHEDULES', '{"message": "📅 *Consulta de Horarios*\n\nPara consultar horarios, por favor proporciona:\n   • Nombre de la carrera o programa\n   • Semestre o nivel\n   • (Opcional) Materia específica\n\nEjemplo: \"Horario de Ingeniería en Sistemas, tercer semestre\"\n\nTambién puedo ayudarte con horarios de:\n   🏫 Horarios de atención administrativa\n   📖 Horarios de biblioteca\n   🏃 Horarios de actividades extracurriculares\n\n¿Qué horario necesitas consultar?"}'::jsonb, 'Schedules command response text');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'MESSAGE_COMMANDS') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_COMMANDS', '{"message": "⚡ *Comandos Disponibles*\n\n/help - Muestra ayuda general del bot\n/horarios - Consulta horarios de clases\n/comandos - Muestra esta lista de comandos\n/start - Reinicia la conversación\n\n💡 *Tip*: No necesitas usar comandos para hacer preguntas. ¡Solo escribe tu consulta!"}'::jsonb, 'Commands list response text');
      end if;

      if not exists (select 1 from cht_parameters where prm_code = 'MESSAGE_START') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_START', '{"message": "👋 ¡Hola! Soy el asistente virtual del Instituto.\n\nEstoy aquí para ayudarte con información sobre:\n • Programas académicos\n   • Admisiones y matrículas\n   • Horarios y calendarios\n • Y mucho más...\n\nEscribe /help para ver todo lo que puedo hacer, o simplemente hazme una pregunta.\n\n¿En qué puedo ayudarte?"}'::jsonb, 'Start/welcome command response text');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_UNKNOWN_COMMAND') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_UNKNOWN_COMMAND', '{"message": "❓ Comando no reconocido.\n\nEscribe /help para ver los comandos disponibles, o simplemente hazme tu pregunta directamente."}'::jsonb, 'Unknown command response text');
      end if;

      -- =====================================================
      -- WhatsApp Bot - User Registration Messages
      -- =====================================================
      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_REQUEST_CEDULA') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_REQUEST_CEDULA', '{"message": "👋 ¡Hola! Bienvenido al asistente virtual del Instituto.\n\nPara poder ayudarte, necesito que te registres primero.\n\nPor favor, envíame tu número de cédula (10 dígitos).\n\nEjemplo: 1234567890"}'::jsonb, 'Message requesting user to provide cedula for registration');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_CEDULA_VALIDATION_ERROR') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_CEDULA_VALIDATION_ERROR', '{"message": "❌ No pude validar tu cédula. Por favor verifica que sea correcta e intenta nuevamente."}'::jsonb, 'Error message when cedula validation fails');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_REGISTRATION_ERROR') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_REGISTRATION_ERROR', '{"message": "❌ Ocurrió un error al registrarte. Por favor intenta nuevamente."}'::jsonb, 'Error message when user registration fails');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_EXTERNAL_USER') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_EXTERNAL_USER', '{"message": "👤 No encontré tu cédula en nuestra base de datos.\n\nSi eres un visitante externo, por favor proporciona:\n1. Tu nombre completo\n2. Tu correo electrónico\n\nEjemplo:\nJuan Pérez\njuan.perez@email.com\n\nO si eres estudiante/docente, verifica que tu cédula sea correcta."}'::jsonb, 'Message for external users not found in institute database');
      end if;

      if not exists (select 1 from cht_parameters where prm_code =
 'MESSAGE_WELCOME_REGISTERED') then
          insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
          values ('RAG_CONFIGURATION', 'MESSAGE_WELCOME_REGISTERED', '{"template": "%s ¡Bienvenido, %s!\n\nHas sido registrado exitosamente como %s.\n\nAhora puedes:\n• Hacer preguntas sobre el instituto\n• Consultar horarios con /horarios\n• Ver ayuda con /help\n\n¿En qué puedo ayudarte hoy?"}'::jsonb, 'Welcome message after successful registration (args: roleEmoji, userName, roleText)');
      end if;

      -- =====================================================
      -- LLM Configuration

    -- =====================================================
    -- LLM Configuration
    -- =====================================================
    if not exists (select 1 from cht_parameters where prm_code = 'LLM_CONFIG') then
        insert into cht_parameters (prm_name, prm_code, prm_data, prm_description)
        values (
            'LLM_CONFIGURATION',
            'LLM_CONFIG',
            '{
                "provider": "groq",
                "apiKey": "YOUR_GROQ_API_KEY_HERE",
                "baseURL": "https://api.groq.com/openai/v1",
                "model": "llama-3.3-70b-versatile",
                "temperature": 0.7,
                "maxTokens": 1000,
                "timeout": 30,
                "systemPrompt": "Eres un asistente virtual del instituto educativo. Tu objetivo es ayudar a estudiantes y profesores con información académica de manera clara, precisa y amigable. Siempre basa tus respuestas en el contexto proporcionado."
            }'::jsonb,
            'LLM provider configuration (groq, openai, anthropic). Use baseURL to connect to OpenAI-compatible APIs.'
        );
    end if;

    -- Alternative OpenAI configuration example (commented out)
    -- To use OpenAI instead of Groq, update LLM_CONFIG with:
    -- {
    --   "provider": "openai",
    --   "apiKey": "YOUR_OPENAI_API_KEY_HERE",
    --   "baseURL": "https://api.openai.com/v1",
    --   "model": "gpt-4o-mini",
    --   "temperature": 0.7,
    --   "maxTokens": 1000,
    --   "timeout": 30,
    --   "systemPrompt": "..."
    -- }

end $$;

-- =====================================================
-- Default Permissions (Admin has all)
-- =====================================================
do $$
begin
    -- Admin permissions
    insert into cht_permissions (prm_rol, prm_funcionality, prm_active)
    values
        ('ROLE_ADMIN', 'FUNC_CHAT', true),
        ('ROLE_ADMIN', 'FUNC_DOCUMENTS', true),
        ('ROLE_ADMIN', 'FUNC_USERS', true),
        ('ROLE_ADMIN', 'FUNC_PARAMETERS', true)
    on conflict (prm_rol, prm_funcionality) do nothing;

    -- User permissions
    insert into cht_permissions (prm_rol, prm_funcionality, prm_active)
    values
        ('ROLE_USER', 'FUNC_CHAT', true)
    on conflict (prm_rol, prm_funcionality) do nothing;

    -- Student permissions
    insert into cht_permissions (prm_rol, prm_funcionality, prm_active)
    values
        ('ROLE_STUDENT', 'FUNC_CHAT', true)
    on conflict (prm_rol, prm_funcionality) do nothing;

    -- Guest permissions
    insert into cht_permissions (prm_rol, prm_funcionality, prm_active)
    values
        ('ROLE_GUEST', 'FUNC_CHAT', true)
    on conflict (prm_rol, prm_funcionality) do nothing;
end $$;
