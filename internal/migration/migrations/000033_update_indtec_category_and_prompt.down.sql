-- Rollback: Restore original INDTEC prompt and category description

DO $$
BEGIN
    -- Restore original DOC_INDTEC description
    IF EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'DOC_INDTEC') THEN
        UPDATE cht_parameters
        SET prm_description = 'INDTEC event category'
        WHERE prm_code = 'DOC_INDTEC';
    END IF;

    -- Restore original RAG_SYSTEM_PROMPT_DOC_INDTEC message
    IF EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'RAG_SYSTEM_PROMPT_DOC_INDTEC') THEN
        UPDATE cht_parameters
        SET prm_data = '{
                "message": "Soy Alfibot 👋, asistente especializado en eventos de INDTEC (Instituto Tecnológico Sudamericano). Te ayudo con información específica sobre nuestros eventos de innovación y tecnología.\n\nReglas: 1) Respondo basándome SOLO en información de eventos INDTEC 🎯 2) Doy detalles sobre fechas, ubicaciones, actividades y requisitos de participación 📅 3) Mantengo respuestas claras y concisas 4) Si necesitas info general del instituto, te sugiero hacer una consulta sin especificar categoría 🏫"
            }'::jsonb,
            prm_description = 'Category-specific system prompt for INDTEC events (DOC_INDTEC)'
        WHERE prm_code = 'RAG_SYSTEM_PROMPT_DOC_INDTEC';
    END IF;
END $$;
