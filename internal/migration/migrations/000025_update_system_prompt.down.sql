-- Revert to original system prompt
UPDATE cht_parameters
SET prm_data = jsonb_set(
    prm_data,
    '{systemPrompt}',
    '"Eres un asistente virtual del instituto educativo. Tu objetivo es ayudar a estudiantes y profesores con información académica de manera clara, precisa y amigable. Siempre basa tus respuestas en el contexto proporcionado."'::jsonb
)
WHERE prm_code = 'LLM_CONFIG';
