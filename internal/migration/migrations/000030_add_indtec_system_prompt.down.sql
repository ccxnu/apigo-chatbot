-- Remove category-specific system prompt for DOC_INDTEC
DELETE FROM cht_parameters WHERE prm_code = 'RAG_SYSTEM_PROMPT_DOC_INDTEC';
