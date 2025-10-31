-- Remove analytics columns from conversation messages table
ALTER TABLE cht_conversation_messages
DROP COLUMN IF EXISTS cvm_queue_time_ms,
DROP COLUMN IF EXISTS cvm_prompt_tokens,
DROP COLUMN IF EXISTS cvm_prompt_time_ms,
DROP COLUMN IF EXISTS cvm_completion_tokens,
DROP COLUMN IF EXISTS cvm_completion_time_ms,
DROP COLUMN IF EXISTS cvm_total_tokens,
DROP COLUMN IF EXISTS cvm_total_time_ms;

-- Drop indexes
DROP INDEX IF EXISTS idx_messages_total_tokens;
DROP INDEX IF EXISTS idx_messages_created_at;
