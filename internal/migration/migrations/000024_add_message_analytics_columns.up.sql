-- Add analytics/stats columns to conversation messages table
ALTER TABLE cht_conversation_messages
ADD COLUMN IF NOT EXISTS cvm_queue_time_ms INT,
ADD COLUMN IF NOT EXISTS cvm_prompt_tokens INT,
ADD COLUMN IF NOT EXISTS cvm_prompt_time_ms INT,
ADD COLUMN IF NOT EXISTS cvm_completion_tokens INT,
ADD COLUMN IF NOT EXISTS cvm_completion_time_ms INT,
ADD COLUMN IF NOT EXISTS cvm_total_tokens INT,
ADD COLUMN IF NOT EXISTS cvm_total_time_ms INT;

-- Comments
COMMENT ON COLUMN cht_conversation_messages.cvm_queue_time_ms IS 'Time spent in queue before processing (milliseconds)';
COMMENT ON COLUMN cht_conversation_messages.cvm_prompt_tokens IS 'Number of prompt tokens used';
COMMENT ON COLUMN cht_conversation_messages.cvm_prompt_time_ms IS 'Time to generate prompt (milliseconds)';
COMMENT ON COLUMN cht_conversation_messages.cvm_completion_tokens IS 'Number of completion tokens generated';
COMMENT ON COLUMN cht_conversation_messages.cvm_completion_time_ms IS 'Time to generate completion (milliseconds)';
COMMENT ON COLUMN cht_conversation_messages.cvm_total_tokens IS 'Total tokens (prompt + completion)';
COMMENT ON COLUMN cht_conversation_messages.cvm_total_time_ms IS 'Total processing time (milliseconds)';

-- Create index for analytics queries
CREATE INDEX IF NOT EXISTS idx_messages_total_tokens ON cht_conversation_messages(cvm_total_tokens) WHERE cvm_total_tokens IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON cht_conversation_messages(cvm_created_at DESC);
