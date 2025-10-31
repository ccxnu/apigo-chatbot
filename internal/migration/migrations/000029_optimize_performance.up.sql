-- Optional performance optimization: reduce temperature for faster, more deterministic responses
-- Lower temperature = faster inference, more focused responses
-- You can skip this migration if you prefer more creative responses

-- Reduce temperature from 0.7 to 0.5 for faster responses
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{value}', '0.5'::jsonb)
WHERE prm_code = 'RAG_LLM_TEMPERATURE';

-- Reduce max tokens from 1000 to 800 for faster generation
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{value}', '800'::jsonb)
WHERE prm_code = 'RAG_LLM_MAX_TOKENS';

-- Comment for reference
COMMENT ON TABLE cht_parameters IS 'System configuration parameters. Lower temperature (0.3-0.5) = faster, more focused. Higher (0.7-0.9) = slower, more creative.';
