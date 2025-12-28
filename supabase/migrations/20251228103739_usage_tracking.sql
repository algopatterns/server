-- Add usage tracking and subscription tiers

-- ============================================================================
-- USAGE TRACKING TABLE
-- ============================================================================
-- Tracks all code generation requests for analytics and rate limiting
-- Supports both authenticated users and anonymous sessions

CREATE TABLE usage_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id TEXT,              -- For anonymous users (from anonsessions manager)
    provider TEXT NOT NULL,       -- "anthropic", "openai"
    model TEXT NOT NULL,          -- "claude-sonnet-4-20250514", "gpt-4o", etc.
    input_tokens INT,             -- Estimated input tokens
    output_tokens INT,            -- Estimated output tokens
    is_byok BOOLEAN DEFAULT false, -- true if user provided their own API key
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for fast queries
CREATE INDEX idx_usage_user_date ON usage_logs(user_id, created_at DESC);
CREATE INDEX idx_usage_session_date ON usage_logs(session_id, created_at DESC);

-- Index for filtering platform usage (is_byok = false means using our keys)
CREATE INDEX idx_usage_platform ON usage_logs(user_id, is_byok, created_at DESC);

-- ============================================================================
-- SUBSCRIPTION TIERS
-- ============================================================================
-- Add tier column to users table
-- Values: 'free', 'pro', 'byok'
--   - free: 100 generations/day with platform keys
--   - pro: Unlimited with platform keys ($15/mo - future)
--   - byok: Unlimited with user's own API key (free)

ALTER TABLE users ADD COLUMN tier TEXT DEFAULT 'free';

-- Ensure only valid tier values
ALTER TABLE users ADD CONSTRAINT tier_check
    CHECK (tier IN ('free', 'pro', 'byok'));

-- Index for tier-based queries
CREATE INDEX idx_users_tier ON users(tier);

-- ============================================================================
-- TIER CHANGE HISTORY (Optional - for analytics)
-- ============================================================================
-- Track when users change tiers (useful for understanding conversion)

CREATE TABLE tier_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    old_tier TEXT,
    new_tier TEXT NOT NULL,
    reason TEXT,                  -- "upgrade", "downgrade", "byok_added", etc.
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_tier_changes_user ON tier_changes(user_id, changed_at DESC);

-- ============================================================================
-- HELPER FUNCTION: Get user's usage count for today
-- ============================================================================
-- This makes rate limiting queries simpler and faster

CREATE OR REPLACE FUNCTION get_user_usage_today(p_user_id UUID)
RETURNS INT AS $$
    SELECT COUNT(*)::INT
    FROM usage_logs
    WHERE user_id = p_user_id
    AND is_byok = false  -- Only count platform key usage
    AND created_at >= CURRENT_DATE
$$ LANGUAGE SQL STABLE;

-- ============================================================================
-- HELPER FUNCTION: Get session's usage count for today
-- ============================================================================
-- For anonymous users rate limiting

CREATE OR REPLACE FUNCTION get_session_usage_today(p_session_id TEXT)
RETURNS INT AS $$
    SELECT COUNT(*)::INT
    FROM usage_logs
    WHERE session_id = p_session_id
    AND is_byok = false
    AND created_at >= CURRENT_DATE
$$ LANGUAGE SQL STABLE;

-- ============================================================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================================================

COMMENT ON TABLE usage_logs IS 'Tracks all code generation requests for rate limiting and analytics';
COMMENT ON COLUMN usage_logs.user_id IS 'Authenticated user (null for anonymous)';
COMMENT ON COLUMN usage_logs.session_id IS 'Anonymous session ID from anonsessions manager';
COMMENT ON COLUMN usage_logs.is_byok IS 'true if user provided own API key (does not count toward quota)';

COMMENT ON TABLE tier_changes IS 'Audit log of subscription tier changes';

COMMENT ON COLUMN users.tier IS 'Subscription tier: free (100/day), pro (unlimited), byok (unlimited with own key)';
