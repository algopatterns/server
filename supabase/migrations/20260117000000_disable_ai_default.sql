-- Disable AI features by default (free tier disabled, BYOK required)
-- To re-enable free tier, revert this migration

-- Change default for new users
ALTER TABLE users
ALTER COLUMN ai_features_enabled SET DEFAULT false;

-- Update existing users to have AI disabled by default
-- They can re-enable it after configuring BYOK
UPDATE users SET ai_features_enabled = false WHERE ai_features_enabled = true;

COMMENT ON COLUMN users.ai_features_enabled IS 'Whether AI features are enabled for this user. Requires BYOK API key to use.';
