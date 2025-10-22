-- =====================================================
-- Cleanup: Remove old message tables (replaced by conversation system)
-- =====================================================

DROP TABLE IF EXISTS public.cht_message_statistics CASCADE;
DROP TABLE IF EXISTS public.cht_messages CASCADE;

DROP INDEX IF EXISTS idx_cht_messages_session;
DROP INDEX IF EXISTS idx_cht_messages_processed;
DROP INDEX IF EXISTS idx_cht_messages_embedding;
DROP INDEX IF EXISTS idx_cht_msg_stats_message;

-- =====================================================
-- Add LLM usage statistics to conversation messages
-- =====================================================

ALTER TABLE public.cht_conversation_messages
    ADD COLUMN IF NOT EXISTS cvm_queue_time_ms      int,
    ADD COLUMN IF NOT EXISTS cvm_prompt_tokens      int,
    ADD COLUMN IF NOT EXISTS cvm_prompt_time_ms     int,
    ADD COLUMN IF NOT EXISTS cvm_completion_tokens  int,
    ADD COLUMN IF NOT EXISTS cvm_completion_time_ms int,
    ADD COLUMN IF NOT EXISTS cvm_total_tokens       int,
    ADD COLUMN IF NOT EXISTS cvm_total_time_ms      int;

CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_total_tokens ON cht_conversation_messages(cvm_total_tokens DESC);
CREATE INDEX IF NOT EXISTS idx_cht_conv_msgs_total_time ON cht_conversation_messages(cvm_total_time_ms DESC);

COMMENT ON COLUMN cht_conversation_messages.cvm_queue_time_ms IS 'LLM queue time in milliseconds';
COMMENT ON COLUMN cht_conversation_messages.cvm_prompt_tokens IS 'Number of tokens in prompt';
COMMENT ON COLUMN cht_conversation_messages.cvm_prompt_time_ms IS 'Time to process prompt in milliseconds';
COMMENT ON COLUMN cht_conversation_messages.cvm_completion_tokens IS 'Number of tokens in completion';
COMMENT ON COLUMN cht_conversation_messages.cvm_completion_time_ms IS 'Time to generate completion in milliseconds';
COMMENT ON COLUMN cht_conversation_messages.cvm_total_tokens IS 'Total tokens used (prompt + completion)';
COMMENT ON COLUMN cht_conversation_messages.cvm_total_time_ms IS 'Total processing time in milliseconds';
