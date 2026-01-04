-- Add is_code_response column to session_messages
-- This distinguishes code responses (should update editor) from question responses (explanations)

ALTER TABLE session_messages
ADD COLUMN IF NOT EXISTS is_code_response BOOLEAN DEFAULT true;

COMMENT ON COLUMN session_messages.is_code_response IS 'Whether this response should update the editor (true for code, false for explanations/questions)';
