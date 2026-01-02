-- Add ai_features_enabled flag to users table
-- Allows users to disable AI features and use the editor manually

ALTER TABLE users
ADD COLUMN IF NOT EXISTS ai_features_enabled BOOLEAN DEFAULT true;

COMMENT ON COLUMN users.ai_features_enabled IS 'Whether AI features are enabled for this user (prompt bar, code generation)';
