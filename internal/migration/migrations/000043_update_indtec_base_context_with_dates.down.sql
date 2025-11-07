-- Revert to previous base context without specific dates

UPDATE cht_parameters
SET prm_data = '{
    "context": "INDTEC es el Congreso Internacional de Tecnología, Investigación, Desarrollo e Innovación organizado por el Instituto Tecnológico Sudamericano. Este evento reúne a profesionales, investigadores, estudiantes y empresas del sector tecnológico para compartir conocimiento, presentar investigaciones y fomentar la innovación. El congreso incluye conferencias magistrales, presentación de ponencias, talleres prácticos, y espacios de networking."
}'::jsonb
WHERE prm_code = 'BASE_CONTEXT_DOC_INDTEC';
