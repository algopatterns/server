-- Add BM25 (Full-Text Search) support for hybrid search
-- This migration adds tsvector columns and indexes for keyword-based search

-- ============================================================================
-- ADD TSVECTOR COLUMNS
-- ============================================================================

-- Add tsvector column for doc_embeddings
ALTER TABLE doc_embeddings
ADD COLUMN IF NOT EXISTS content_tsvector tsvector;

-- Add tsvector column for example_strudels
ALTER TABLE example_strudels
ADD COLUMN IF NOT EXISTS searchable_tsvector tsvector;

-- ============================================================================
-- CREATE GIN INDEXES FOR FAST FULL-TEXT SEARCH
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_doc_content_tsvector
ON doc_embeddings USING GIN(content_tsvector);

CREATE INDEX IF NOT EXISTS idx_example_searchable_tsvector
ON example_strudels USING GIN(searchable_tsvector);

-- ============================================================================
-- POPULATE TSVECTOR COLUMNS WITH EXISTING DATA
-- ============================================================================

-- Populate doc_embeddings tsvector
UPDATE doc_embeddings
SET content_tsvector = to_tsvector('english',
  COALESCE(page_name, '') || ' ' ||
  COALESCE(section_title, '') || ' ' ||
  COALESCE(content, '')
);

-- Populate example_strudels tsvector
UPDATE example_strudels
SET searchable_tsvector = to_tsvector('english',
  COALESCE(title, '') || ' ' ||
  COALESCE(description, '') || ' ' ||
  COALESCE(code, '') || ' ' ||
  COALESCE(array_to_string(tags, ' '), '')
);

-- ============================================================================
-- CREATE TRIGGERS TO AUTO-UPDATE TSVECTOR ON INSERT/UPDATE
-- ============================================================================

-- Trigger function for doc_embeddings
CREATE OR REPLACE FUNCTION doc_embeddings_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  new.content_tsvector := to_tsvector('english',
    COALESCE(new.page_name, '') || ' ' ||
    COALESCE(new.section_title, '') || ' ' ||
    COALESCE(new.content, '')
  );
  RETURN new;
END
$$ LANGUAGE plpgsql;

-- Create trigger for doc_embeddings
DROP TRIGGER IF EXISTS tsvector_update_doc_embeddings ON doc_embeddings;
CREATE TRIGGER tsvector_update_doc_embeddings
BEFORE INSERT OR UPDATE ON doc_embeddings
FOR EACH ROW EXECUTE FUNCTION doc_embeddings_tsvector_trigger();

-- Trigger function for example_strudels
CREATE OR REPLACE FUNCTION example_strudels_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  new.searchable_tsvector := to_tsvector('english',
    COALESCE(new.title, '') || ' ' ||
    COALESCE(new.description, '') || ' ' ||
    COALESCE(new.code, '') || ' ' ||
    COALESCE(array_to_string(new.tags, ' '), '')
  );
  RETURN new;
END
$$ LANGUAGE plpgsql;

-- Create trigger for example_strudels
DROP TRIGGER IF EXISTS tsvector_update_example_strudels ON example_strudels;
CREATE TRIGGER tsvector_update_example_strudels
BEFORE INSERT OR UPDATE ON example_strudels
FOR EACH ROW EXECUTE FUNCTION example_strudels_tsvector_trigger();
