-- Update base context for DOC_INDTEC with specific event dates and details

UPDATE cht_parameters
SET prm_data = '{
    "context": "INDTEC es el III Congreso Internacional de Investigación Científica 2025, organizado por el Instituto Tecnológico Sudamericano. El evento se realizará los días 27 y 28 de noviembre de 2025 en modalidad híbrida (presencial y virtual). Este congreso reúne a profesionales, investigadores, estudiantes y empresas del sector tecnológico para compartir conocimiento, presentar investigaciones y fomentar la innovación. Sus principales objetivos son difundir resultados de investigación, fomentar publicaciones en revistas indexadas y promover el desarrollo en investigación en la Zona 7 del país."
}'::jsonb
WHERE prm_code = 'BASE_CONTEXT_DOC_INDTEC';
