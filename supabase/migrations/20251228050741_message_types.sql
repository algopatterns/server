-- Migration: Add message_type to session_messages
-- Description: Distinguishes between AI conversation messages and chat messages
-- Created: 2025-01-28

-- Add message_type column with default for existing rows
ALTER TABLE session_messages
ADD COLUMN message_type TEXT NOT NULL DEFAULT 'user_prompt'
CHECK (message_type IN ('user_prompt', 'ai_response', 'chat'));

-- Update existing rows: role 'user' = 'user_prompt', role 'assistant' = 'ai_response'
UPDATE session_messages
SET message_type = CASE
  WHEN role = 'user' THEN 'user_prompt'
  WHEN role = 'assistant' THEN 'ai_response'
END;

-- Remove default after migration
ALTER TABLE session_messages
ALTER COLUMN message_type DROP DEFAULT;

-- Add index for filtering by message type
CREATE INDEX idx_messages_type ON session_messages(session_id, message_type, created_at);

-- Update comment
COMMENT ON COLUMN session_messages.message_type IS 'Type of message: user_prompt (user prompt to AI), ai_response (AI generated code), chat (user chat message)';
