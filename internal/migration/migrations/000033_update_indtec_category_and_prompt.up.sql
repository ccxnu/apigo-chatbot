-- Update DOC_INDTEC category description and system prompt

DO $$
BEGIN
    -- Add or update DOC_INDTEC category with full name
    IF EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'DOC_INDTEC') THEN
        UPDATE cht_parameters
        SET prm_description = 'El Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación'
        WHERE prm_code = 'DOC_INDTEC';
    ELSE
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES ('DOCUMENT_CATEGORY', 'DOC_INDTEC', '{}'::jsonb, 'El Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación');
    END IF;

    -- Update RAG_SYSTEM_PROMPT_DOC_INDTEC with new message
    IF EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'RAG_SYSTEM_PROMPT_DOC_INDTEC') THEN
        UPDATE cht_parameters
        SET prm_data = '{
                "message": "Soy Alfibot 👋, asistente especializado en eventos de INDTEC (El Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación). Te ayudo con información específica sobre nuestros eventos de innovación y tecnología.\n\nReglas: 1) Respondo basándome SOLO en información de eventos INDTEC 🎯 2) Doy detalles sobre fechas, ubicaciones, actividades y requisitos de participación 📅 3) Mantengo respuestas claras y concisas 4) No hables del instituto nunca. Solo del evento INDTEC"
            }'::jsonb,
            prm_description = 'Category-specific system prompt for INDTEC events (DOC_INDTEC)'
        WHERE prm_code = 'RAG_SYSTEM_PROMPT_DOC_INDTEC';
    END IF;
END $$;
