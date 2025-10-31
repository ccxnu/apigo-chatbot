-- Revert performance optimizations to original values

-- Restore temperature to 0.7
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{value}', '0.7'::jsonb)
WHERE prm_code = 'RAG_LLM_TEMPERATURE';

-- Restore max tokens to 1000
UPDATE cht_parameters
SET prm_data = jsonb_set(prm_data, '{value}', '1000'::jsonb)
WHERE prm_code = 'RAG_LLM_MAX_TOKENS';
