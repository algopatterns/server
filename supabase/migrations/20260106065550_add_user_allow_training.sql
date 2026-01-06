-- add user-controlled allow_training field for per-strudel training opt-in
-- this is separate from admin-controlled use_in_training

ALTER TABLE user_strudels
ADD COLUMN IF NOT EXISTS allow_training BOOLEAN DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_user_strudels_allow_training
ON user_strudels(allow_training)
WHERE allow_training = true;

COMMENT ON COLUMN user_strudels.allow_training IS 'User-controlled: whether user permits this strudel for training';

-- update search function to require all 4 conditions:
-- 1. user.training_consent (global user consent)
-- 2. strudel.allow_training (user per-strudel consent)
-- 3. strudel.use_in_training (admin curation)
-- 4. strudel.is_public (must be public)
CREATE OR REPLACE FUNCTION search_user_strudels(
    query_embedding extensions.vector(1536),
    match_count int DEFAULT 3
)
RETURNS TABLE (
    id UUID,
    title TEXT,
    description TEXT,
    code TEXT,
    tags TEXT[],
    user_id UUID,
    similarity FLOAT
)
LANGUAGE plpgsql STABLE
AS $$
BEGIN
    PERFORM set_config('search_path', 'extensions, public', true);

    RETURN QUERY
    SELECT
        us.id,
        us.title,
        us.description,
        us.code,
        us.tags,
        us.user_id,
        1 - (us.embedding <=> query_embedding) AS similarity
    FROM user_strudels us
    INNER JOIN users u ON us.user_id = u.id
    WHERE us.allow_training = true
      AND us.use_in_training = true
      AND us.is_public = true
      AND us.embedding IS NOT NULL
      AND u.training_consent = true
    ORDER BY us.embedding <=> query_embedding
    LIMIT match_count;
END;
$$;

-- update trigger to auto-disable allow_training when strudel becomes private
CREATE OR REPLACE FUNCTION user_strudels_training_check_trigger() RETURNS trigger AS $$
BEGIN
  IF NEW.is_public = false THEN
    NEW.allow_training := false;
    NEW.use_in_training := false;
  END IF;
  RETURN NEW;
END
$$ LANGUAGE plpgsql;
