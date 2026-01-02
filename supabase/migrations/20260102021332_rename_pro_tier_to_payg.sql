-- Rename 'pro' tier to 'payg' (pay as you go)
-- This better reflects the usage-based pricing model

-- Update existing users with 'pro' tier to 'payg'
UPDATE users SET tier = 'payg' WHERE tier = 'pro';

-- Update tier_changes history
UPDATE tier_changes SET old_tier = 'payg' WHERE old_tier = 'pro';
UPDATE tier_changes SET new_tier = 'payg' WHERE new_tier = 'pro';

-- Drop and recreate the constraint with the new valid values
ALTER TABLE users DROP CONSTRAINT IF EXISTS tier_check;
ALTER TABLE users ADD CONSTRAINT tier_check
    CHECK (tier IN ('free', 'payg', 'byok'));

-- Update column comment
COMMENT ON COLUMN users.tier IS 'Subscription tier: free (100/day), payg (pay as you go - usage based), byok (unlimited with own key)';
