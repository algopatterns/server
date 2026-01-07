-- add AI contribution score to track how much of strudel code came from AI
ALTER TABLE user_strudels
ADD COLUMN IF NOT EXISTS ai_contribution_score FLOAT DEFAULT 0.0;

COMMENT ON COLUMN user_strudels.ai_contribution_score IS 'Score from 0.0-1.0 indicating AI contribution to the code';

CREATE INDEX IF NOT EXISTS idx_user_strudels_ai_score
ON user_strudels(ai_contribution_score)
WHERE ai_contribution_score > 0.3;
