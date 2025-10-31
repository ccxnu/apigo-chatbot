-- =====================================================
-- Rollback LLM Model Fix
-- =====================================================

-- Revert to llama-3.3-70b-versatile
UPDATE cht_parameters
SET prm_data = jsonb_set(
    prm_data,
    '{model}',
    '"llama-3.3-70b-versatile"'
)
WHERE prm_code = 'LLM_CONFIG';
