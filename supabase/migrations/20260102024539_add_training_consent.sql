-- Add training consent for ethical AI training
-- Users opt-in to allow their public strudels to be used as RAG examples

-- ============================================================================
-- USERS TABLE: Add training_consent column
-- ============================================================================

ALTER TABLE users
ADD COLUMN IF NOT EXISTS training_consent BOOLEAN DEFAULT false;

-- Index for efficient filtering of consenting users
CREATE INDEX IF NOT EXISTS idx_users_training_consent
ON users(training_consent)
WHERE training_consent = true;

COMMENT ON COLUMN users.training_consent IS 'Whether user consents to their public strudels being used for AI training';

-- ============================================================================
-- USER_STRUDELS TABLE: Add training and search columns
-- ============================================================================

-- Add allow_training flag
ALTER TABLE user_strudels
ADD COLUMN IF NOT EXISTS allow_training BOOLEAN DEFAULT false;

-- Add embedding column for vector similarity search
ALTER TABLE user_strudels
ADD COLUMN IF NOT EXISTS embedding extensions.vector(1536);

-- Add tsvector column for BM25 full-text search
ALTER TABLE user_strudels
ADD COLUMN IF NOT EXISTS searchable_tsvector tsvector;

-- Index for trainable strudels
CREATE INDEX IF NOT EXISTS idx_user_strudels_training
ON user_strudels(allow_training)
WHERE allow_training = true;

-- Vector similarity search index
CREATE INDEX IF NOT EXISTS idx_user_strudels_embedding
ON user_strudels
USING ivfflat (embedding extensions.vector_cosine_ops)
WITH (lists = 100);

-- GIN index for BM25 full-text search
CREATE INDEX IF NOT EXISTS idx_user_strudels_searchable_tsvector
ON user_strudels USING GIN(searchable_tsvector);

COMMENT ON COLUMN user_strudels.allow_training IS 'Whether this strudel can be used for AI training (requires user consent + is_public)';
COMMENT ON COLUMN user_strudels.embedding IS 'Embedding of title + description + code for semantic search';

-- ============================================================================
-- POPULATE TSVECTOR FOR EXISTING STRUDELS
-- ============================================================================

UPDATE user_strudels
SET searchable_tsvector = to_tsvector('english',
  COALESCE(title, '') || ' ' ||
  COALESCE(description, '') || ' ' ||
  COALESCE(code, '') || ' ' ||
  COALESCE(array_to_string(tags, ' '), '')
)
WHERE searchable_tsvector IS NULL;

-- ============================================================================
-- TRIGGER TO AUTO-UPDATE TSVECTOR ON INSERT/UPDATE
-- ============================================================================

CREATE OR REPLACE FUNCTION user_strudels_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.searchable_tsvector := to_tsvector('english',
    COALESCE(NEW.title, '') || ' ' ||
    COALESCE(NEW.description, '') || ' ' ||
    COALESCE(NEW.code, '') || ' ' ||
    COALESCE(array_to_string(NEW.tags, ' '), '')
  );
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tsvector_update_user_strudels ON user_strudels;
CREATE TRIGGER tsvector_update_user_strudels
BEFORE INSERT OR UPDATE ON user_strudels
FOR EACH ROW EXECUTE FUNCTION user_strudels_tsvector_trigger();

-- ============================================================================
-- TRIGGER TO ENFORCE allow_training REQUIRES is_public
-- ============================================================================

CREATE OR REPLACE FUNCTION user_strudels_training_check_trigger() RETURNS trigger AS $$
BEGIN
  -- If strudel is being made private, disable training
  IF NEW.is_public = false AND NEW.allow_training = true THEN
    NEW.allow_training := false;
  END IF;
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS training_check_user_strudels ON user_strudels;
CREATE TRIGGER training_check_user_strudels
BEFORE INSERT OR UPDATE ON user_strudels
FOR EACH ROW EXECUTE FUNCTION user_strudels_training_check_trigger();

-- ============================================================================
-- VECTOR SEARCH FUNCTION FOR TRAINABLE USER STRUDELS
-- ============================================================================

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
      AND us.is_public = true
      AND us.embedding IS NOT NULL
      AND u.training_consent = true
    ORDER BY us.embedding <=> query_embedding
    LIMIT match_count;
END;
$$;

COMMENT ON FUNCTION search_user_strudels IS 'Search trainable user strudels by vector similarity (requires user consent + strudel allow_training + is_public)';

-- ============================================================================
-- CLEANUP: DROP EXAMPLE_STRUDELS TABLE AND RELATED OBJECTS
-- ============================================================================

-- Drop the search function first (depends on table)
DROP FUNCTION IF EXISTS search_examples(extensions.vector(1536), int);

-- Drop the trigger and function
DROP TRIGGER IF EXISTS tsvector_update_example_strudels ON example_strudels;
DROP FUNCTION IF EXISTS example_strudels_tsvector_trigger();

-- Drop the table
DROP TABLE IF EXISTS example_strudels;
