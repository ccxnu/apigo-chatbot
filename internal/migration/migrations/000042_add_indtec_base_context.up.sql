-- Add base context parameter for DOC_INDTEC category
-- This provides essential information that is always injected when event_filter=DOC_INDTEC
-- This solves the problem of generic queries not retrieving relevant information

DO $$
BEGIN
    -- Add BASE_CONTEXT_DOC_INDTEC parameter with essential event information
    IF NOT EXISTS (SELECT 1 FROM cht_parameters WHERE prm_code = 'BASE_CONTEXT_DOC_INDTEC') THEN
        INSERT INTO cht_parameters (prm_name, prm_code, prm_data, prm_description)
        VALUES (
            'EVENT_BASE_CONTEXT',
            'BASE_CONTEXT_DOC_INDTEC',
            '{
                "context": "INDTEC es el Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación organizado por el Instituto Tecnológico Sudamericano. Este evento reúne a profesionales, investigadores, estudiantes y empresas del sector tecnológico para compartir conocimiento, presentar investigaciones y fomentar la innovación. El congreso incluye conferencias magistrales, presentación de ponencias, talleres prácticos, y espacios de networking."
            }'::jsonb,
            'Base context automatically injected for INDTEC event queries (DOC_INDTEC)'
        );
    ELSE
        -- Update if already exists
        UPDATE cht_parameters
        SET prm_data = '{
                "context": "INDTEC es el Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación organizado por el Instituto Tecnológico Sudamericano. Este evento reúne a profesionales, investigadores, estudiantes y empresas del sector tecnológico para compartir conocimiento, presentar investigaciones y fomentar la innovación. El congreso incluye conferencias magistrales, presentación de ponencias, talleres prácticos, y espacios de networking."
            }'::jsonb,
            prm_description = 'Base context automatically injected for INDTEC event queries (DOC_INDTEC)'
        WHERE prm_code = 'BASE_CONTEXT_DOC_INDTEC';
    END IF;
END $$;
