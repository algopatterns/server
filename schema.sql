-- enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- create document_chunks table
CREATE TABLE IF NOT EXISTS document_chunks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_name TEXT NOT NULL,
    page_url TEXT NOT NULL,
    section_title TEXT,
    content TEXT NOT NULL,
    embedding vector(1536),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- create vector similarity search index
-- ivfflat is good for datasets < 1M vectors
CREATE INDEX IF NOT EXISTS document_chunks_embedding_idx 
ON document_chunks 
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- create index on page_name for faster lookups
CREATE INDEX IF NOT EXISTS document_chunks_page_name_idx 
ON document_chunks(page_name);

-- create index on created_at for maintenance queries
CREATE INDEX IF NOT EXISTS document_chunks_created_at_idx 
ON document_chunks(created_at DESC);

-- verify the table was created
SELECT 
    tablename, 
    schemaname 
FROM pg_tables 
WHERE tablename = 'document_chunks';

-- check if vector extension is enabled
SELECT * FROM pg_extension WHERE extname = 'vector';