-- =====================================================
-- Fix LLM Model - Change to llama-3.1-8b-instant
-- =====================================================

-- Update LLM_CONFIG to use llama 3.1 8B instead of 70B
UPDATE cht_parameters
SET prm_data = jsonb_set(
    prm_data,
    '{model}',
    '"llama-3.1-8b-instant"'
)
WHERE prm_code = 'LLM_CONFIG';

-- Verify the update
DO $$
DECLARE
    v_model TEXT;
BEGIN
    SELECT prm_data->>'model' INTO v_model FROM cht_parameters WHERE prm_code = 'LLM_CONFIG';
    RAISE NOTICE 'LLM_CONFIG updated. Current model: %', v_model;
END $$;
